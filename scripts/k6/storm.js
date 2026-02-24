import http from 'k6/http';
import ws from 'k6/ws';
import { check, sleep } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import { textSummary } from 'https://jslib.k6.io/k6-summary/0.0.2/index.js';
import { Trend } from 'k6/metrics';

// Custom metrics for granular reporting
const authLoginTrend = new Trend('auth_login_duration');
const authRegisterTrend = new Trend('auth_register_duration');
const httpMessageTrend = new Trend('http_message_duration');

export const options = {
    stages: [
        { duration: '30s', target: 50 },
        { duration: '1m', target: 50 },
        { duration: '10s', target: 0 },
    ],
};

const BASE_URL = 'http://localhost:8080';
const WS_URL = 'ws://localhost:8080/ws';

export default function () {
    const username = `user_${randomString(8)}_${__VU}`;
    const password = 'password123';

    const registerPayload = JSON.stringify({
        user_id: username,
        password: password,
        display_name: `Test User ${__VU}`,
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    // 1. Register
    const registerRes = http.post(`${BASE_URL}/auth/register`, registerPayload, params);
    authRegisterTrend.add(registerRes.timings.duration);

    check(registerRes, {
        'register success or conflict': (r) => r.status === 201 || r.status === 500,
    });

    // 2. Login
    const loginPayload = JSON.stringify({
        user_id: username,
        password: password,
    });

    const loginRes = http.post(`${BASE_URL}/auth/login`, loginPayload, params);
    authLoginTrend.add(loginRes.timings.duration);

    check(loginRes, {
        'login success': (r) => r.status === 200,
    });

    if (loginRes.status !== 200) return;

    let cookies = loginRes.cookies['access_token'];
    let token = '';
    if (cookies && cookies.length > 0) {
        token = cookies[0].value;
    }

    // 3. Test HTTP Message (once per VU to measure API latency)
    const msgPayload = JSON.stringify({
        content: `HTTP message from ${username}`,
    });
    const msgRes = http.post(`${BASE_URL}/channels/1/messages`, msgPayload, {
        headers: { 'Content-Type': 'application/json' },
    });
    httpMessageTrend.add(msgRes.timings.duration);
    check(msgRes, { 'http message success': (r) => r.status === 201 });

    // 4. WebSocket Flow
    const wsUrlWithAuth = `${WS_URL}?token=${token}`;
    const res = ws.connect(wsUrlWithAuth, params, function (socket) {
        socket.on('open', function () {
            socket.setInterval(function timeout() {
                const msg = JSON.stringify({
                    channel_id: 1,
                    content: `WS message from ${username}`,
                });
                socket.send(msg);
            }, 5000); // 5s to leave air for HTTP tests
        });

        socket.setTimeout(function () {
            socket.close();
        }, 20000);
    });

    check(res, { 'websocket connected successfully': (r) => r && r.status === 101 });
    sleep(1);
}

export function handleSummary(data) {
    const vus = data.metrics.vus ? data.metrics.vus.values.max : 0;
    const wsReceived = data.metrics.ws_msgs_received ? data.metrics.ws_msgs_received.values.count : 0;
    const wsSent = data.metrics.ws_msgs_sent ? data.metrics.ws_msgs_sent.values.count : 0;
    const duration = data.state.testRunDurationMs / 1000;

    const markdownReport = `
# 🌩️ STORM Day - Rapport de Tir k6

**Date :** ${new Date().toLocaleString()}
**Durée du test :** ${duration.toFixed(2)} secondes
**Utilisateurs virtuels max (VUs) :** ${vus}

## 📊 Métriques WebSockets
- **Messages échangés (Total) :** ${wsReceived + wsSent}
- **Temps de connexion WS (médiane) :** ${data.metrics.ws_connecting ? data.metrics.ws_connecting.values.med.toFixed(2) : 0} ms

## ⚡ Performance API Auth
- **Login (moyenne) :** ${data.metrics.auth_login_duration ? data.metrics.auth_login_duration.values.avg.toFixed(2) : 0} ms
- **Register (moyenne) :** ${data.metrics.auth_register_duration ? data.metrics.auth_register_duration.values.avg.toFixed(2) : 0} ms

## 💬 Performance Messaging (HTTP)
- **Envoi message (moyenne) :** ${data.metrics.http_message_duration ? data.metrics.http_message_duration.values.avg.toFixed(2) : 0} ms

## ✅ Fiabilité
- **Taux de succès HTTP :** ${data.metrics.http_req_failed ? (100 - (data.metrics.http_req_failed.values.rate * 100)).toFixed(2) : 100}%
`;

    return {
        'stdout': textSummary(data, { indent: ' ', enableColors: true }),
        '/scripts/report.md': markdownReport,
    };
}

