from flask import Flask, request
app = Flask(__name__)

@app.route('/')
def verify():
    print("I ran")
    print(request.headers)
    return 200


if __name__ == "__main__":
    app.run(debug=True)