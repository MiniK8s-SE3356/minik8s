import os
import json
import base64
import uuid
import oss2

def function(params):
    params = json.loads(params)
    
    info_str = params['info']
    image_data = params['image']
    image_bytes = base64.b64decode(image_data)
    image_uid = str(uuid.uuid4())
    
    with open(f'{image_uid}.jpg', 'wb') as f:
        f.write(image_bytes)
    
    auth = oss2.Auth(
        'LTAI5tMwHrYu3hi5bJu7tq34',
        'IrWRqqRkv1GEG5V0xwzoxwmVsjVCM4'
    )
    
    bucket = oss2.Bucket(
        auth,
        'oss-cn-shanghai.aliyuncs.com',
        'xubbbb-chartbed'
    )
    
    bucket.put_object_from_file('minik8s/' + image_uid + '.jpg', f'{image_uid}.jpg')
    url = "https://xubbbb-chartbed.oss-cn-shanghai.aliyuncs.com/minik8s/" + image_uid + ".jpg"
    
    os.remove(f'{image_uid}.jpg')
    
    resp = {
        "result_img_url": url,
        "info": info_str
    }
    
    return json.dumps(resp)
    