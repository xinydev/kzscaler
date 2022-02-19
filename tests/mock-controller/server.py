from flask import Flask, request

app = Flask(__name__)

enabled_service = ["service1", "service2"]
zerostate_service = ["service1"]


@app.route("/enabled")
def get_enabled_service():
    print(request.headers)
    return "|".join(enabled_service)


@app.route("/zerostate")
def get_zerostate_service():
    print(request.headers)
    return "|".join(zerostate_service)


if __name__ == "__main__":
    app.run("0.0.0.0", 10021)
