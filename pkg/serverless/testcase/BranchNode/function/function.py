import json

def function(params):
    params = json.loads(params)
    handler_type = params["handler"]
    resp = 0
    
    if handler_type == 'opencv':
        resp = 1
    else:
        resp = 0
        
    return str(resp)