package types

const (
	SBATCH_HEADER          = "#!/bin/bash\n"
	SBATCH_JOB_NAME        = "#SBATCH --job-name=%s\n"
	SBATCH_PARTITION       = "#SBATCH --partition=%s\n"
	SBATCH_NODES           = "#SBATCH --nodes=%s\n"
	SBATCH_NTASKS_PER_NODE = "#SBATCH --ntasks-per-node=%s\n"
	SBATCH_CPU_PER_TASK    = "#SBATCH --cpus-per-task=%s\n"
	SBATCH_GPUS            = "#SBATCH --gres=gpu:%s\n"
	SBATCH_OUTPUT_FILE     = "#SBATCH --output=%s\n"
	SBATCH_ERR_FILE        = "#SBATCH --error=%s\n"

	SBATCH_MODULE_LOAD = "module load %s\n"

	SBATCH_SUBMIT = "sbatch %s\n"
)
