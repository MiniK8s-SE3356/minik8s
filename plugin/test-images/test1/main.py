import math
import time

def compute_factorial(n):
    start_time = time.time()  # 记录开始时间
    result = math.factorial(n)
    end_time = time.time()  # 记录结束时间
    # print(f"计算 {n} 的阶乘结果为：{result} (耗时 {end_time - start_time:.4f} 秒)")

def main():
    while True:  # 无限循环
        compute_factorial(50000)  # 计算一个较大数的阶乘

if __name__ == "__main__":
    main()