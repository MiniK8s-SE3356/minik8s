package server

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	gpu_types "github.com/MiniK8s-SE3356/minik8s/pkg/gpu/types"
	"github.com/MiniK8s-SE3356/minik8s/pkg/gpu/utils/ssh"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/poller"
	minik8s_zip "github.com/MiniK8s-SE3356/minik8s/pkg/utils/zip"
)

var JobManagerUrl string

const (
	prefix              = "/api/v1"
	RequireGPUJobSuffix = prefix + "/RequireGPUJob"
	UpdateGPUJobSuffix  = prefix + "/UpdateGPUJob"

	ZipFileName = "job.zip"
	// ZipDirName  = "job"

	A100NodeHostname = "sylogin.hpc.sjtu.edu.cn"
	DGX2NodeHostname = "pilogin.hpc.sjtu.edu.cn"
)

type JobServer struct {
	Job gpu_types.SlurmJob
	// TODO: Does sshclient need to be protected by mutex?
	SSHClient *ssh.SSHClient
}

func NewJobServer(jobName string) (*JobServer, error) {
	job, zipContent, err := GetJobDetailFromJobManager(jobName)
	if err != nil {
		fmt.Println("failed to create new job server because GetJobDetailFromJobManager failed")
		return nil, err
	}

	// write zipContent to file
	err = os.WriteFile(ZipFileName, *zipContent, os.ModePerm)
	if err != nil {
		fmt.Println("failed to create new job server because write zipContent to file failed")
		return nil, err
	}

	// decompress zip file to /app/job directory
	fileDir := job.Metadata.UUID
	err = minik8s_zip.DecompressZipFile(ZipFileName, fileDir)
	if err != nil {
		fmt.Println("failed to create new job server because decompress zip file failed")
		return nil, err
	}

	workDir := filepath.Join(fileDir, job.WorkDir)

	// create a .slurm script file in /app/job directory
	slurmScript, err := os.Create(filepath.Join(workDir, "job.slurm"))
	if err != nil {
		fmt.Println("failed to create new job server because create slurm script file failed")
		return nil, err
	}
	defer slurmScript.Close()

	var modules string
	for _, module := range job.Modules {
		modules += module + " "
	}

	var executionCmds string
	for _, cmd := range job.RunCmds {
		executionCmds += cmd + "\n"
	}

	slurmContent := gpu_types.SBATCH_HEADER +
		fmt.Sprintf(gpu_types.SBATCH_JOB_NAME, job.Metadata.Name) +
		fmt.Sprintf(gpu_types.SBATCH_PARTITION, job.Partition) +
		fmt.Sprintf(gpu_types.SBATCH_NODES, job.Nodes) +
		fmt.Sprintf(gpu_types.SBATCH_NTASKS_PER_NODE, job.NTasksPerNode) +
		fmt.Sprintf(gpu_types.SBATCH_CPU_PER_TASK, job.CPUPerTask) +
		fmt.Sprintf(gpu_types.SBATCH_GPUS, job.GPUNum) +
		// fmt.Sprintf(gpu_types.SBATCH_OUTPUT_FILE, job.OutputFile) +
		// fmt.Sprintf(gpu_types.SBATCH_ERR_FILE, job.ErrFile)
		//! Temporarily use the default output and error file: 'jobid'.out/.err
		fmt.Sprintf(gpu_types.SBATCH_OUTPUT_FILE, "%j.out") +
		fmt.Sprintf(gpu_types.SBATCH_ERR_FILE, "%j.err") + "\n" +
		fmt.Sprintf(gpu_types.SBATCH_MODULE_LOAD, modules) + "\n" +
		// Execution commands
		executionCmds

	_, err = slurmScript.WriteString(slurmContent)
	if err != nil {
		fmt.Println("failed to create new job server because write slurm script file failed")
		return nil, err
	}

	sshConfig := ssh.SSHConfig{
		Username: job.Username,
		Password: job.Password,
	}
	if job.Partition == "a100" {
		sshConfig.Hostname = A100NodeHostname
	} else if job.Partition == "dgx2" {
		sshConfig.Hostname = DGX2NodeHostname
	} else {
		fmt.Println("failed to create new job server because partition is invalid")
		return nil, fmt.Errorf("partition is invalid")
	}

	//!debug//
	fmt.Println("sshConfig: ")
	fmt.Println("Username: ", sshConfig.Username)
	fmt.Println("Password:", sshConfig.Password)
	fmt.Println("Hostname: ", sshConfig.Hostname)
	//!debug//

	sshClient, err := ssh.NewSSHClient(sshConfig)
	if err != nil {
		fmt.Println("failed to create new job server because NewSSHClient failed")
		return nil, err
	}

	result := &JobServer{
		Job:       *job,
		SSHClient: sshClient,
	}

	return result, nil
}

func (js *JobServer) Run() {
	fmt.Println("JobServer is running")

	//? Maybe we should change ZipDirName to job's UUID
	fileDir := filepath.Join(js.Job.Metadata.UUID, js.Job.WorkDir)
	js.SSHClient.PostDirectory(fileDir, fileDir)
	// Execute the compile commands
	cmds := []string{
		fmt.Sprintf("module load %s", "gcc cuda"),
		fmt.Sprintf("cd %s", fileDir),
	}
	cmds = append(cmds, js.Job.CompileCmds...)
	// submit the job
	cmds = append(cmds, fmt.Sprintf(gpu_types.SBATCH_SUBMIT, "job.slurm"))

	out, err := js.SSHClient.BatchCmd(cmds)
	if err != nil {
		fmt.Println("failed to run job because BatchCmd failed")
		return
	}
	fmt.Println("JobServer run result: ", out)

	// parse out to get job id
	re := regexp.MustCompile(`Submitted batch job (\d+)`)
	matches := re.FindStringSubmatch(out)
	if len(matches) > 1 {
		js.Job.JobID = matches[1]
		fmt.Println("JobServer run job id: ", js.Job.JobID)
	} else {
		// TODO: handle this error
		fmt.Println("failed to get job id")
		return
	}

	poller.PollerStaticPeriod(
		10*time.Second,
		js.GetJobState,
		true,
	)
}

func (js *JobServer) GetJobOutput() (string, error) {
	fileDir := filepath.Join(js.Job.Metadata.UUID, js.Job.WorkDir)
	cmds := []string{
		fmt.Sprintf("cd %s", fileDir),
		fmt.Sprintf("cat %s.out", js.Job.JobID),
	}

	out, err := js.SSHClient.BatchCmd(cmds)
	if err != nil {
		fmt.Println("failed to get job output because BatchCmd failed")
		return "", err
	}
	return out, nil
}

func (js *JobServer) GetJobState() {
	cmd := fmt.Sprintf(gpu_types.JOBSTATE_CHECK, js.Job.JobID)
	out, err := js.SSHClient.BatchCmd([]string{cmd})
	if err != nil {
		fmt.Println("failed to get job state because BatchCmd failed")
		return
	}

	outlines := strings.Split(out, "\n")

	// because in this output, it may output all related jobs' state
	// only the first line is the job's state we want

	if len(outlines) > 0 {
		stateInfos := strings.Split(outlines[0], " ")
		if len(stateInfos) < 6 {
			fmt.Println("failed to get job state because stateInfos is invalid")
			return
		}

		newState := stateInfos[5]
		isUpdate := false
		if newState != js.Job.State {
			js.Job.State = newState
			fmt.Println("JobServer job state: ", js.Job.State)
			isUpdate = true
		}

		if newState == "COMPLETED" {
			// TODO: fetch job output files
			output, err := js.GetJobOutput()
			if err != nil {
				fmt.Println("failed to get job output")
				return
			}
			js.Job.Result = output
			isUpdate = true
		}

		if isUpdate {
			var response_str string
			status_code, err := httpRequest.PostRequestByObject(
				JobManagerUrl+UpdateGPUJobSuffix,
				js.Job,
				&response_str,
			)
			if err != nil {
				fmt.Println("failed to update job state")
				return
			}
			if status_code != 200 {
				fmt.Println("failed to update job state, status_code: ", status_code)
				return
			}
			fmt.Println("JobServer update job state result: ", response_str)
		}
	}

	fmt.Println("JobServer get job state result: ", out)
}

func GetJobDetailFromJobManager(jobName string) (*gpu_types.SlurmJob, *[]byte, error) {
	param_list := make(map[string]string)
	param_list["jobName"] = jobName

	var response struct {
		Job        gpu_types.SlurmJob `json:"job"`
		ZipContent []byte             `json:"zipContent"`
	}

	statusCode, err := httpRequest.GetRequestByObject(
		JobManagerUrl+RequireGPUJobSuffix,
		param_list,
		&response,
	)

	if err != nil {
		fmt.Println("failed to get job")
		return nil, nil, err
	}

	if statusCode != 200 {
		fmt.Println("failed to get job, statusCode: ", statusCode)
		return nil, nil, fmt.Errorf("failed to get job, statusCode: %d", statusCode)
	}

	return &response.Job, &response.ZipContent, nil
}
