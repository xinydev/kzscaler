import requests


def query_service():
    headers = {
        "HOST": "service1:123"
    }
    req = requests.get("http://127.0.0.1:18000", headers=headers)
    print(req.status_code, req.content, req.headers)


if __name__ == "__main__":
    query_service()
