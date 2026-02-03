import http from 'k6/http';
import ws from 'k6/ws';
import { check, sleep } from 'k6';

const BASE_URL = __ENV.GATEWAY_URL || 'http://localhost:8080';
const ACCESS_TOKEN = __ENV.ACCESS_TOKEN || '';
const SUBJECT = __ENV.SUBJECT || 'storm.events';
const CHANNEL_ID = __ENV.CHANNEL_ID || '';

const PUB_RATE = Number(__ENV.PUB_RATE || 10);
const WS_VUS = Number(__ENV.WS_VUS || 5);
const PUB_VUS = Number(__ENV.PUB_VUS || 5);
const DURATION = __ENV.DURATION || '30s';

export const options = {
  scenarios: {
    publish_http: {
      executor: 'constant-vus',
      vus: PUB_VUS,
      duration: DURATION,
      exec: 'publishHttp',
    },
    ws_stream: {
      executor: 'constant-vus',
      vus: WS_VUS,
      duration: DURATION,
      exec: 'wsStream',
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.02'],
    http_req_duration: ['p(95)<500'],
  },
};

function authHeaders() {
  if (!ACCESS_TOKEN) {
    return { 'Content-Type': 'application/json' };
  }
  return {
    'Content-Type': 'application/json',
    Authorization: `Bearer ${ACCESS_TOKEN}`,
  };
}

function wsURL() {
  let url = BASE_URL.replace('http://', 'ws://').replace('https://', 'wss://');
  let qs = '';
  if (CHANNEL_ID) {
    qs = `channel_id=${encodeURIComponent(CHANNEL_ID)}`;
  } else {
    qs = `subject=${encodeURIComponent(SUBJECT)}`;
  }
  return `${url}/ws?${qs}`;
}

export function publishHttp() {
  const payload = JSON.stringify({
    payload: JSON.stringify({
      user: 'loadtest',
      message: `ping-${__VU}-${__ITER}`,
    }),
  });
  const res = http.post(
    `${BASE_URL}/publish?subject=${encodeURIComponent(SUBJECT)}`,
    payload,
    { headers: authHeaders() },
  );
  check(res, {
    'publish status 202/200': (r) => r.status === 202 || r.status === 200,
  });
  sleep(1 / Math.max(PUB_RATE, 1));
}

export function wsStream() {
  const params = { headers: authHeaders() };
  const res = ws.connect(wsURL(), params, function (socket) {
    socket.on('open', function () {
      socket.send(JSON.stringify({ user: 'loadtest', message: 'hello' }));
    });
    socket.on('message', function () {});
    socket.on('error', function () {});
    socket.setTimeout(function () {
      socket.close();
    }, 1000);
  });
  check(res, { 'ws connected': (r) => r && r.status === 101 });
  sleep(1);
}
