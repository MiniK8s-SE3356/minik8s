#include <stdio.h>

#define WIDTH 16 // 定义矩阵的宽度

// CUDA内核函数，用于矩阵乘法
__global__ void matrixMultiply(const float* A, const float* B, float* C, int width)
{
    int col = blockIdx.x * blockDim.x + threadIdx.x;
    int row = blockIdx.y * blockDim.y + threadIdx.y;

    if (row < width && col < width)
    {
        float result = 0.0f;
        for (int k = 0; k < width; ++k)
        {
            result += A[row * width + k] * B[k * width + col];
        }
        C[row * width + col] = result;
    }
}

int main()
{
    int numElements = WIDTH * WIDTH;
    size_t size = numElements * sizeof(float);

    // 分配主机内存
    float* h_A = (float*)malloc(size);
    float* h_B = (float*)malloc(size);
    float* h_C = (float*)malloc(size);

    // 初始化矩阵A和B
    for (int i = 0; i < numElements; ++i)
    {
        h_A[i] = rand()/(float)RAND_MAX;
        h_B[i] = rand()/(float)RAND_MAX;
    }

    // 分配设备内存
    float* d_A = NULL;
    float* d_B = NULL;
    float* d_C = NULL;
    cudaMalloc((void**)&d_A, size);
    cudaMalloc((void**)&d_B, size);
    cudaMalloc((void**)&d_C, size);

    // 将矩阵数据从主机复制到设备
    cudaMemcpy(d_A, h_A, size, cudaMemcpyHostToDevice);
    cudaMemcpy(d_B, h_B, size, cudaMemcpyHostToDevice);

    // 定义线程块和网格大小
    dim3 threadsPerBlock(16, 16);
    dim3 blocksPerGrid((WIDTH + threadsPerBlock.x - 1) / threadsPerBlock.x,
                       (WIDTH + threadsPerBlock.y - 1) / threadsPerBlock.y);

    // 启动内核计算矩阵乘法
    matrixMultiply<<<blocksPerGrid, threadsPerBlock>>>(d_A, d_B, d_C, WIDTH);

    // 将结果从设备复制回主机
    cudaMemcpy(h_C, d_C, size, cudaMemcpyDeviceToHost);

    // 打印一些结果进行检查
    printf("Matrix A:\n");
    for (int i = 0; i < WIDTH; ++i)
    {
        for (int j = 0; j < WIDTH; ++j)
        {
            printf("%f ", h_A[i * WIDTH + j]);
        }
        printf("\n");
    }

    printf("Matrix B:\n");
    for (int i = 0; i < WIDTH; ++i)
    {
        for (int j = 0; j < WIDTH; ++j)
        {
            printf("%f ", h_B[i * WIDTH + j]);
        }
        printf("\n");
    }

    printf("Matrix C:\n");
    for (int i = 0; i < WIDTH; ++i)
    {
        for (int j = 0; j < WIDTH; ++j)
        {
            printf("%f ", h_C[i * WIDTH + j]);
        }
        printf("\n");
    }

    // 释放设备内存
    cudaFree(d_A);
    cudaFree(d_B);
    cudaFree(d_C);

    // 释放主机内奇
    free(h_A);
    free(h_B);
    free(h_C);

    return 0;
}