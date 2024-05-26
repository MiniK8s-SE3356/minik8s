import json

def function(params):
    params = json.loads(params)
    x = params["x"]
    y = params["y"]
    x = x + y
    resp = {
        "sum": x
    }
    return json.dumps(resp)