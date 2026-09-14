#!/usr/bin/env python3
"""Check that the deployed API health endpoint is reachable and healthy."""

import sys
import json
import requests


BASE_URL = "https://clownfish-app-fet99.ondigitalocean.app"
TIMEOUT_SECONDS = 10

def add_item(payload) -> requests.Response | None:
    url = f"{BASE_URL}/api/items"
    try:
        response = requests.post(url, data=json.dumps(payload))
    except requests.RequestException as error:
        print(f"Request failed: {error}", file=sys.stderr)
        return None
    return response

def main() -> int:
    url = f"{BASE_URL}/health"

    try:
        response = requests.get(url, timeout=TIMEOUT_SECONDS)
    except requests.RequestException as error:
        print(f"Request failed: {error}", file=sys.stderr)
        return 1

    if response.status_code != requests.codes.ok:
        print(
            f"Health check failed: HTTP {response.status_code} from {url}",
            file=sys.stderr,
        )
        return 1

    try:
        payload = response.json()
    except ValueError:
        print("Health check failed: response was not valid JSON", file=sys.stderr)
        return 1

    if payload.get("status") != "ok":
        print(f"Health check failed: unexpected response {payload!r}", file=sys.stderr)
        return 1

    response = add_item({'title':'asdf'})
    if response:
        print(f'add_item | {response.text}')
    else:
        print('add_item no response')

    print(f"Health check passed: {url}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
