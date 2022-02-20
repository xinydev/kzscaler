import time

from flask import Flask

app = Flask(__name__)

service_cache = {
    "service1": 0,
    "service2": 1
}


@app.route("/service", methods=['GET'])
def get_all_service():
    # service1%10&service2%10
    service_list = []
    for k in service_cache:
        service_list.append(f"{k}%{service_cache[k]}")
    return "&".join(service_list)


@app.route("/scale_up/<service_name>", methods=['GET'])
def scale_up(service_name):
    time.sleep(20)
    print("scale up request:", service_name)
    return "OK"


if __name__ == "__main__":
    app.run("0.0.0.0", 10021)
