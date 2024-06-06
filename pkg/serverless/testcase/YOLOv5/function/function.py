import io
import json
import requests
from PIL import Image
from yolov5 import helpers
import base64

def function(params):
    params = json.loads(params)
    url = params['url']
    response = requests.get(url)
    img = Image.open(io.BytesIO(response.content))
    
    model = helpers.load_model(
        model_path='./yolov5s.pt',
        device='cpu'
    )
    
    class_names = model.module.names if hasattr(model, 'module') else model.names
    
    results = model(img)
    predictions = results.pred[0].tolist()
    
    info_str = f'YOLOv5 detected {len(predictions)} objects!\n'
    for x1, y1, x2, y2, conf, cls in predictions:
        class_name = class_names[int(cls)]
        info_str += f'Bounding box: ({x1}, {y1}, {x2}, {y2}), confidence: {conf}, class: {class_name}\n'
        
    results.render()
    rendered_image = Image.fromarray(results.ims[0])
    byte_array = io.BytesIO()
    rendered_image.save(byte_array, format='JPEG')
    img_bytes = byte_array.getvalue()
    img_str = base64.b64encode(img_bytes).decode()
    
    resp = {
        "info": info_str,
        "image": img_str
    }
    
    return json.dumps(resp)