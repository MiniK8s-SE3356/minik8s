from flask import Flask, request, jsonify
import math

app = Flask(__name__)

def simpson_rule(f, a, b, n):
    if n % 2:
        n += 1  # n 必须是偶数
    h = (b - a) / n
    s = f(a) + f(b)
    for i in range(1, n, 2):
        s += 4 * f(a + i * h)
    for i in range(2, n-1, 2):
        s += 2 * f(a + i * h)
    return s * h / 3

@app.route('/', methods=['GET'])
def integrate():
    # 计算函数 sin(x) 在区间 [0, π] 上的积分
    lower_bound = request.args.get('lower_bound', default=0, type=float)
    upper_bound = request.args.get('upper_bound', default=math.pi, type=float)
    result = simpson_rule(math.sin, lower_bound, upper_bound, 100000)
    return jsonify({"integral": result})

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, threaded=True)