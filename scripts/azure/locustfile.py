import json
import os
import ssl
import time
from urllib.parse import urlencode

import websocket
from locust import User, between, events, task


HOST = os.getenv("TARGET_HOST", "http://localhost:8080").rstrip("/")
WS_PATH = os.getenv("WS_PATH", "/ws")
WS_SUBJECT = os.getenv("WS_SUBJECT", "storm.events")
WS_CHANNEL_ID = os.getenv("WS_CHANNEL_ID", "").strip()
ACCESS_TOKEN = os.getenv("ACCESS_TOKEN", "").strip()
PING_INTERVAL = float(os.getenv("WS_PING_INTERVAL_SECONDS", "30"))
MESSAGE_INTERVAL = float(os.getenv("WS_MESSAGE_INTERVAL_SECONDS", "0"))
CONNECTION_HOLD_SECONDS = float(os.getenv("WS_CONNECTION_HOLD_SECONDS", "600"))
VERIFY_TLS = os.getenv("VERIFY_TLS", "true").lower() == "true"


def ws_base_url() -> str:
    if HOST.startswith("https://"):
        return "wss://" + HOST[len("https://") :]
    if HOST.startswith("http://"):
        return "ws://" + HOST[len("http://") :]
    return HOST


class StormWebSocketUser(User):
    wait_time = between(0.1, 1.0)
    abstract = False

    def on_start(self) -> None:
        self.ws = None
        self.connected_at = 0.0
        self.last_ping = 0.0
        self.last_message = 0.0
        self._connect()

    def on_stop(self) -> None:
        self._close()

    def _connect(self) -> None:
        params = {}
        if WS_CHANNEL_ID:
            params["channel_id"] = WS_CHANNEL_ID
        else:
            params["subject"] = WS_SUBJECT
        if ACCESS_TOKEN:
            params["token"] = ACCESS_TOKEN

        url = f"{ws_base_url()}{WS_PATH}?{urlencode(params)}"
        sslopt = {"cert_reqs": ssl.CERT_REQUIRED if VERIFY_TLS else ssl.CERT_NONE}

        start = time.perf_counter()
        try:
            self.ws = websocket.create_connection(url, timeout=30, sslopt=sslopt)
            self.connected_at = time.time()
            self.last_ping = self.connected_at
            self.last_message = self.connected_at
            elapsed_ms = (time.perf_counter() - start) * 1000
            events.request.fire(
                request_type="WS",
                name="connect",
                response_time=elapsed_ms,
                response_length=0,
                exception=None,
                context={},
            )
        except Exception as exc:
            elapsed_ms = (time.perf_counter() - start) * 1000
            events.request.fire(
                request_type="WS",
                name="connect",
                response_time=elapsed_ms,
                response_length=0,
                exception=exc,
                context={},
            )
            raise

    def _close(self) -> None:
        if self.ws is None:
            return
        try:
            self.ws.close()
        except Exception:
            pass
        finally:
            self.ws = None

    def _send_ping(self) -> None:
        if self.ws is None:
            return
        start = time.perf_counter()
        try:
            self.ws.ping("storm")
            elapsed_ms = (time.perf_counter() - start) * 1000
            events.request.fire(
                request_type="WS",
                name="ping",
                response_time=elapsed_ms,
                response_length=0,
                exception=None,
                context={},
            )
        except Exception as exc:
            elapsed_ms = (time.perf_counter() - start) * 1000
            events.request.fire(
                request_type="WS",
                name="ping",
                response_time=elapsed_ms,
                response_length=0,
                exception=exc,
                context={},
            )
            raise

    def _send_message(self) -> None:
        if self.ws is None:
            return
        payload = json.dumps(
            {
                "user": "azure-load",
                "message": f"ping-{self.environment.runner.user_count}",
                "ts": int(time.time()),
            }
        )
        start = time.perf_counter()
        try:
            self.ws.send(payload)
            elapsed_ms = (time.perf_counter() - start) * 1000
            events.request.fire(
                request_type="WS",
                name="send",
                response_time=elapsed_ms,
                response_length=len(payload),
                exception=None,
                context={},
            )
        except Exception as exc:
            elapsed_ms = (time.perf_counter() - start) * 1000
            events.request.fire(
                request_type="WS",
                name="send",
                response_time=elapsed_ms,
                response_length=len(payload),
                exception=exc,
                context={},
            )
            raise

    @task
    def hold_connection(self) -> None:
        now = time.time()
        if self.ws is None:
            self._connect()
            return

        if now - self.connected_at >= CONNECTION_HOLD_SECONDS:
            self._close()
            self._connect()
            return

        try:
            if PING_INTERVAL > 0 and now - self.last_ping >= PING_INTERVAL:
                self._send_ping()
                self.last_ping = now

            if MESSAGE_INTERVAL > 0 and now - self.last_message >= MESSAGE_INTERVAL:
                self._send_message()
                self.last_message = now
        except Exception:
            self._close()
            self._connect()

