import function
import json
from flask import Flask, request

app = Flask(__name__)

@app.route('/', methods=['GET'])
def hello_minik8s_serverless():
    return 'Hello, Welcome to Minik8s Serverless!'

@app.route('/api/v1/callfunc', methods=['POST'])
def callfunc():
    print("-----------------[A serverless function is being called!]-----------------")
    try:
        body = json.loads(request.get_data())
    except json.JSONDecodeError:
        body = ""
    finally:
        params = body["params"]
        # make sure is a string
        if not isinstance(params, str):
            return json.dumps({"error": "params must be a string"}), 400
        
        print("params: ", params)
        
        res = function.function(params)
        # make sure res is str
        if not isinstance(res, str):
            return json.dumps({"error": "function must return a string"}), 500
        
        print("res: ", res)
        print("-----------------[A serverless function has been called!]-----------------")
    
    return json.dumps({"result": res}), 200

if __name__ == '__main__':
    app.run(
        host='0.0.0.0',
        # debug=True,
        threaded=True,
    )