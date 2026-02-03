#!/usr/bin/env bash
set -euo pipefail

SECRET="${JWT_SECRET:-dev-secret}"
SUBJECT="${1:-user-1}"
IAT="${IAT:-$(date +%s)}"
EXP="${EXP:-$((IAT + 3600))}"

python3 - "$SECRET" "$SUBJECT" "$IAT" "$EXP" <<'PY'
import base64, hashlib, hmac, json, sys

secret, sub, iat, exp = sys.argv[1], sys.argv[2], int(sys.argv[3]), int(sys.argv[4])

header = {"alg": "HS256", "typ": "JWT"}
payload = {"sub": sub, "iat": iat, "exp": exp}

def b64url(data: bytes) -> str:
    return base64.urlsafe_b64encode(data).rstrip(b"=").decode()

h = b64url(json.dumps(header, separators=(",", ":")).encode())
p = b64url(json.dumps(payload, separators=(",", ":")).encode())
msg = f"{h}.{p}".encode()
sig = hmac.new(secret.encode(), msg, hashlib.sha256).digest()
print(f"{h}.{p}.{b64url(sig)}")
PY