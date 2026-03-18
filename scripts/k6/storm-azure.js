import http from 'k6/http';
import ws from 'k6/ws';
import { check, sleep } from 'k6';
import { Trend, Counter } from 'k6/metrics';

// ─── Native helper (no remote jslib) ──────────
function generateRandomString(length) {
    const chars = 'abcdefghijklmnopqrstuvwxyz0123456789';
    let res = '';
    for (let i = 0; i < length; i++) {
        res += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return res;
}

// ─── Custom metrics ──────────────────────────
const authLoginTrend = new Trend('auth_login_duration');
const authRegisterTrend = new Trend('auth_register_duration');
const httpMessageTrend = new Trend('http_message_duration');
const wsConnectErrors = new Counter('ws_connect_errors');

// ─── Configuration via env vars ──────────────
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const WS_URL = __ENV.WS_URL || 'ws://localhost:8080/ws';
const TARGET = parseInt(__ENV.TARGET_VUS || '250');
const RAMP_UP = __ENV.RAMP_UP || '2m';
const HOLD = __ENV.HOLD || '5m';
const RAMP_DN = __ENV.RAMP_DOWN || '30s';

export const options = {
    stages: [
        { duration: RAMP_UP, target: TARGET },
        { duration: HOLD, target: TARGET },
        { duration: RAMP_DN, target: 0 },
    ],
    thresholds: {
        'http_req_failed': ['rate<0.01'],    // < 1% errors
        'auth_login_duration': ['p(95)<500'],    // p95 login < 500ms
    },
    thresholdAbortOnFail: false,
};

export default function () {
    const username = `u_${generateRandomString(6)}_${__VU}_${__ITER}`;
    const password = 'password123';

    const params = {
        headers: { 'Content-Type': 'application/json' },
    };

    // 1. Register
    const registerRes = http.post(`${BASE_URL}/auth/register`, JSON.stringify({
        user_id: username,
        password: password,
        display_name: `Load User ${__VU}`,
    }), params);
    authRegisterTrend.add(registerRes.timings.duration);
    check(registerRes, {
        'register ok': (r) => r.status === 201 || r.status === 409, // 409 is user already exists
    });

    // 2. Login
    const loginRes = http.post(`${BASE_URL}/auth/login`, JSON.stringify({
        user_id: username,
        password: password,
    }), params);
    authLoginTrend.add(loginRes.timings.duration);
    check(loginRes, { 'login ok': (r) => r.status === 200 });

    if (loginRes.status !== 200) return;

    let token = '';
    const cookies = loginRes.cookies['access_token'];
    if (cookies && cookies.length > 0) {
        token = cookies[0].value;
    }

    // 3. HTTP message
    const msgRes = http.post(`${BASE_URL}/channels/1/messages`, JSON.stringify({
        content: `msg from ${username}`,
    }), { headers: { 'Content-Type': 'application/json' } });
    httpMessageTrend.add(msgRes.timings.duration);
    check(msgRes, { 'message ok': (r) => r.status === 201 });

    // 4. WebSocket – hold connection for 60s
    const wsUrl = `${WS_URL}?token=${token}`;
    const res = ws.connect(wsUrl, params, function (socket) {
        socket.on('open', function () {
            socket.setInterval(function () {
                socket.send(JSON.stringify({
                    channel_id: 1,
                    content: `ws from ${username}`,
                }));
            }, 10000);
        });

        socket.on('error', function () {
            wsConnectErrors.add(1);
        });

        socket.setTimeout(function () {
            socket.close();
        }, 60000);
    });

    check(res, { 'ws connected': (r) => r && r.status === 101 });
    if (!res || res.status !== 101) {
        wsConnectErrors.add(1);
    }

    sleep(1);
}

export function handleSummary(data) {
    const vus = data.metrics.vus ? data.metrics.vus.values.max : 0;
    const duration = data.state.testRunDurationMs / 1000;

    const md = `
# 🌩️ STORM – Azure Load Test Report
**Max VUs :** ${vus}
**Duration :** ${duration.toFixed(0)}s
`;

    return {
        'stdout': JSON.stringify(data.metrics), // Minimal output to avoid external libs
        '/scripts/azure-report.md': md,
    };
}
