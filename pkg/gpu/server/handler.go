package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
	minik8s_yaml "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/yaml"
	gpu_types "github.com/MiniK8s-SE3356/minik8s/pkg/gpu/types"
	minik8s_container "github.com/MiniK8s-SE3356/minik8s/pkg/types/container"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
)

const EtcdGpuJobPrefix = "/minik8s/gpujob/"
const GPUJobPodNamePrefix = "gpujob-"
const GPUJobPodImage = "levixubbbb/jobserver-image:latest"
const JobManagerIPEnv = "JOBMANAGERIP"
const JobManagerPortEnv = "JOBMANAGERPORT"
const JobNameEnv = "JOBNAME"

func SubmitGPUJobHandler(c *gin.Context) {
	// Get the job from the request
	var jobRequest struct {
		JobDesc    gpu_types.SlurmJob `json:"jobDesc"`
		ZipContent []byte             `json:"zipContent"`
	}
	if err := c.ShouldBindJSON(&jobRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := SubmitGPUJobProccess(&jobRequest.JobDesc, &jobRequest.ZipContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func SubmitGPUJobProccess(job *gpu_types.SlurmJob, zip_content *[]byte) (string, error) {

	req := make(map[string]interface{})

	// Check the job's name
	tmp, err := EtcdCli.Exist(EtcdGpuJobPrefix + job.Metadata.Name)
	if err != nil {
		fmt.Println("failed to check existence in etcd")
		return "failed to check existence in etcd", err
	}
	if tmp {
		fmt.Println("job has existed")
		return "job has existed", nil
	}

	// Save the zip content to directory
	zipFilePath := path.Join(GPUJobZipDir, job.Metadata.UUID+".zip")
	err = os.WriteFile(zipFilePath, *zip_content, os.ModePerm)
	if err != nil {
		fmt.Println("failed to save the zip content to file")
		return "failed to save the zip content to file", err
	}

	// Save Job to etcd
	value, err := json.Marshal(job)
	if err != nil {
		fmt.Println("failed to translate into json ", err.Error())
		return "failed to translate into json ", err
	}
	err = EtcdCli.Put(EtcdGpuJobPrefix+job.Metadata.Name, string(value))
	if err != nil {
		fmt.Println("failed to write to etcd ", err.Error())
		return "failed to write to etcd", err
	}

	// Request APIServer to add pod
	podDesc := minik8s_yaml.PodDesc{
		ApiVersion: "v1",
		Kind:       "Pod",
		Metadata: struct {
			Name   string            `yaml:"name" json:"name"`
			Labels map[string]string `yaml:"labels" json:"labels"`
		}{
			Name: GPUJobPodNamePrefix + job.Metadata.UUID,
			Labels: map[string]string{
				"job": job.Metadata.Name,
			},
		},
		Spec: pod.PodSpec{
			Containers: []minik8s_container.Container{
				{
					Name:  "gpujob",
					Image: GPUJobPodImage,
					Env: []minik8s_container.EnvVar{
						{
							Name: JobManagerIPEnv,
							// TODO: get the job manager ip
							Value: ControlPanelIP,
						},
						{
							Name:  JobManagerPortEnv,
							Value: JobManagerPort,
						},
						{
							Name:  JobNameEnv,
							Value: job.Metadata.Name,
						},
					},
				},
			},
		},
	}

	req["namespace"] = "default"
	req["podDesc"] = podDesc

	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("failed to translate into json")
		return "failed to translate into json", err
	}

	result, err := httpRequest.PostRequest(
		// TODO: get the root url from the config
		"http://localhost:"+APIServerPort+"/api/v1/AddPod",
		jsonData,
	)
	if err != nil {
		fmt.Println("error when post request")
		return "error when post request", err
	}

	fmt.Println(result)

	return "job submitted", nil
}

func GetGPUJobHandler(c *gin.Context) {
	jobName := c.Query("name")

	if jobName == "" {
		pairs, err := EtcdCli.GetWithPrefix(EtcdGpuJobPrefix)
		if err != nil {
			fmt.Println("failed to get the jobs from etcd")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the jobs from etcd"})
			return
		}

		result := make(map[string]interface{}, 0)

		for _, pair := range pairs {
			var job gpu_types.SlurmJob
			err = json.Unmarshal([]byte(pair.Value), &job)
			if err != nil {
				fmt.Println("failed to unmarshal the job")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unmarshal the job"})
				return
			}
			result[pair.Key] = job
		}

		c.JSON(http.StatusOK, result)
	} else {
		value, err := EtcdCli.Get(EtcdGpuJobPrefix + jobName)
		if err != nil {
			fmt.Println("failed to get the job from etcd")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get the job from etcd"})
			return
		}
		var job gpu_types.SlurmJob
		err = json.Unmarshal(value, &job)
		if err != nil {
			fmt.Println("failed to unmarshal the job")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unmarshal the job"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			EtcdGpuJobPrefix + jobName: job,
		})
	}
}

func RequireGPUJobHandler(c *gin.Context) {
	jobName := c.Query("jobName")
	fmt.Println("jobName: ", jobName)
	if jobName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no field 'jobName'"})
		return
	}

	job, zipContent, err := RequireGPUJobProccess(jobName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"job":        job,
		"zipContent": zipContent,
	})
}

func RequireGPUJobProccess(jobName string) (*gpu_types.SlurmJob, *[]byte, error) {
	// Get the job from etcd
	value, err := EtcdCli.Get(EtcdGpuJobPrefix + jobName)
	if err != nil {
		fmt.Println("failed to get the job from etcd")
		return nil, nil, err
	}
	var job gpu_types.SlurmJob
	err = json.Unmarshal(value, &job)
	if err != nil {
		fmt.Println("failed to unmarshal the job")
		return nil, nil, err
	}

	// Get the zip content from the directory
	zipFilePath := path.Join(GPUJobZipDir, job.Metadata.UUID+".zip")
	zipContent, err := os.ReadFile(zipFilePath)
	if err != nil {
		fmt.Println("failed to read the zip content from file")
		return nil, nil, err
	}

	return &job, &zipContent, nil
}

func UpdateGPUJobHandler(c *gin.Context) {
	var jobRequest gpu_types.SlurmJob

	if err := c.ShouldBindJSON(&jobRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if jobRequest.State == "COMPLETED" {
		// Kill this pod
		req := make(map[string]interface{})
		req["name"] = GPUJobPodNamePrefix + jobRequest.Metadata.UUID
		req_url := "http://localhost:" + APIServerPort + "/api/v1/RemovePod"

		var resp_str string
		statusCode, err := httpRequest.PostRequestByObject(
			req_url,
			req,
			&resp_str,
		)
		if err != nil {
			fmt.Println("failed to remove pod")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove pod"})
			return
		}
		if statusCode != 200 {
			fmt.Println("failed to remove pod, status_code: ", statusCode)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove pod"})
			return
		}

		fmt.Println("remove pod result: ", resp_str)
	}

	value, err := json.Marshal(jobRequest)
	if err != nil {
		fmt.Println("failed to translate into json ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to translate into json"})
	}
	err = EtcdCli.Put(EtcdGpuJobPrefix+jobRequest.Metadata.Name, string(value))
	if err != nil {
		fmt.Println("failed to write to etcd ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write to etcd"})
	}
}
