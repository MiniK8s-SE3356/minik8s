#include <stdio.h>

// 定义矩阵的维度
#define WIDTH 16
#define HEIGHT 16

// CUDA内核用于矩阵加法
__global__ void matrixAdd(const float* A, const float* B, float* C, int width, int height)
{
    int col = blockIdx.x * blockDim.x + threadIdx.x;
    int row = blockIdx.y * blockDim.y + threadIdx.y;

    if (col < width && row < height)
    {
        int index = row * width + col;
        C[index] = A[index] + B[index];
    }
}

int main()
{
    int numElements = WIDTH * HEIGHT;
    size_t size = numElements * sizeof(float);

    // 为每个矩阵分配主机内存
    float* h_A = (float*)malloc(size);
    float* h_B = (float*)malloc(size);
    float* h_C = (float*)malloc(size);

    // 初始化矩阵数据
    for (int i = 0; i < numElements; ++i)
    {
        h_A[i] = rand()/(float)RAND_MAX;
        h_B[i] = rand()/(float)RAND_MAX;
    }

    // 为每个矩阵分配设备内存
    float* d_A = NULL;
    cudaMalloc((void**)&d_A, size);
    float* d_B = NULL;
    cudaMalloc((void**)&d_B, size);
    float* d_C = NULL;
    cudaMalloc((void**)&d_C, size);

    // 将矩阵数据从主机复制到设备
    cudaMemcpy(d_A, h_A, size, cudaMemcpyHostToDevice);
    cudaMemcpy(d_B, h_B, size, cudaMemcpyHostToDevice);

    // 定义线程块的大小和网格的大小
    dim3 threadsPerBlock(16, 16);
    dim3 blocksPerGrid((WIDTH + threadsPerBlock.x - 1) / threadsPerBlock.x, (HEIGHT + threadsPerBlock.y - 1) / threadsPerBlock.y);
    
    // 启动内核
    matrixAdd<<<blocksPerGrid, threadsPerBlock>>>(d_A, d_B, d_C, WIDTH, HEIGHT);

    // 从设备复制结果回主机
    cudaMemcpy(h_C, d_C, size, cudaMemcpyDeviceToHost);

    // 打印矩阵A
    printf("Matrix A:\n");
    for (int i = 0; i < HEIGHT; ++i)
    {
        for (int j = 0; j < WIDTH; ++j)
        {
            printf("%f ", h_A[i * WIDTH + j]);
        }
        printf("\n");
    }

    // 打印矩阵B
    printf("Matrix B:\n");
    for (int i = 0; i < HEIGHT; ++i)
    {
        for (int j = 0; j < WIDTH; ++j)
        {
            printf("%f ", h_B[i * WIDTH + j]);
        }
        printf("\n");
    }

    // 打印矩阵C
    printf("Matrix C:\n");
    for (int i = 0; i < HEIGHT; ++i)
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

    // 释放主机内存
    free(h_A);
    free(h_B);
    free(h_C);

    return 0;
}