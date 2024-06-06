import json
import requests
import base64
import cv2
import numpy as np

def function(params):
    params = json.loads(params)
    
    url = params['url']
    response = requests.get(url)
    image_data = np.frombuffer(response.content, np.uint8)
    img = cv2.imdecode(image_data, cv2.IMREAD_COLOR)
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    
    face_cascade = cv2.CascadeClassifier(cv2.data.haarcascades + 'haarcascade_frontalface_default.xml')
    faces = face_cascade.detectMultiScale(gray, 1.1, 4)
    
    info_str = f'OpenCV Haar detector detected {len(faces)} faces!'
    
    for (x, y, w, h) in faces:
        cv2.rectangle(img, (x, y), (x+w, y+h), (255, 0, 0), 2)
        
    is_success, im_buf_arr = cv2.imencode(".jpg", img)
    byte_array = im_buf_arr.tobytes()
    image_str = base64.b64encode(byte_array).decode()
    
    resp = {
        "info": info_str,
        "image": image_str
    }
    
    return json.dumps(resp)