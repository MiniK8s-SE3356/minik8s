import json
import requests
import base64
import dlib
import cv2
import numpy as np

def function(params):
    params = json.loads(params)
    url = params['url']
    response = requests.get(url)
    image_data = np.frombuffer(response.content, np.uint8)
    img = cv2.imdecode(image_data, cv2.IMREAD_COLOR)
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    
    detector = dlib.get_frontal_face_detector()
    predictor = dlib.shape_predictor('shape_predictor_68_face_landmarks.dat')
    faces = detector(gray, 1)
    
    for face in faces:
        x1, y1, x2, y2 = face.left(), face.top(), face.right(), face.bottom()
        cv2.rectangle(img, (x1, y1), (x2, y2), (0, 255, 0), 3)
        landmarks = predictor(gray, face)
        for n in range(0, 68):
            x = landmarks.part(n).x
            y = landmarks.part(n).y
            cv2.circle(img, (x, y), 4, (255, 0, 0), -1)
    
    info_str = f'dlib detector detected {len(faces)} faces!'
    
    is_success, im_buf_arr = cv2.imencode(".jpg", img)
    byte_array = im_buf_arr.tobytes()
    image_str = base64.b64encode(byte_array).decode()
    
    resp = {
        "info": info_str,
        "image": image_str
    }
    
    return json.dumps(resp)