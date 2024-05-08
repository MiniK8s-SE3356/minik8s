package cmdline

import "github.com/spf13/cobra"

var DeleteFuncTable = map[string]func(names []string) error{
	"Pod":        deletePod,
	"Service":    deleteService,
	"ReplicaSet": deleteReplicaSet,
	"Namespace":  deleteNamespace,
}

func DeleteCmdHandler(cmd *cobra.Command, args []string) {

}

func deletePod(names []string) error {
	return nil
}

func deleteService(names []string) error {
	return nil
}

func deleteReplicaSet(names []string) error {
	return nil
}

func deleteNamespace(names []string) error {
	return nil
}
