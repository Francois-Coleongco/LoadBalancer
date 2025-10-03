import argparse


from flask import Flask, request

app = Flask(__name__)

@app.route("/", methods=["GET", "POST"])
def unauthed_get_post_endpoint():
    if request.method == "GET":
        return "<p>YOOOOOOO</p>"
    else:
        data = request.form['data']
        return f"received data: {data}"

if __name__ == "__main__":

    parser = argparse.ArgumentParser(description="port configuration")


    parser.add_argument("-p", "--port", type=int, help="enter a valid port 0-65535")

    args = parser.parse_args()

    app.run(debug=True, port=args.port)
