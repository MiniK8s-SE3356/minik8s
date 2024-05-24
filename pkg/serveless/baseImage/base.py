import function
import json
from flask import Flask, request

app = Flask(__name__)

@app.route('/', methods=['GET'])
def hello_minik8s_serveless():
    return 'Hello, Welcome to Minik8s Serveless!'

@app.route('/callfunc', methods=['POST'])
def callfunc():
    try:
        arguments = json.loads(request.get_data())
    except json.JSONDecodeError:
        arguments = ""
    finally:
        res = function.function(arguments)
    
    return json.dumps(res), 200

if __name__ == '__main__':
    app.run(
        debug=True,
        threaded=True,
    )