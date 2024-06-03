package formatprint

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func PrintNodes(str string) {
	var nodes []struct {
		Metadata struct {
			ID     string            `json:"id"`
			Name   string            `json:"name"`
			Labels map[string]string `json:"labels"`
		} `json:"Metadata"`
		Status struct {
			Hostname   string   `json:"hostname"`
			IP         string   `json:"ip"`
			Condition  []string `json:"condition"`
			CPUPercent float64  `json:"cpuPercent"`
			MemPercent float64  `json:"memPercent"`
			NumPods    int      `json:"numPods"`
			UpdateTime string   `json:"updateTime"`
		} `json:"Status"`
	}

	err := json.Unmarshal([]byte(str), &nodes)
	if err != nil {
		fmt.Println("Error parsing JSON: ", err)
		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)

	fmt.Fprintln(writer, "ID\tName\tHostname\tIP\tCondition\tCPU Usage (%)\tMemory Usage (%)\tPod Count\tLast Update")

	for _, node := range nodes {
		// Join the condition array into a single string
		conditionStr := strings.Join(node.Status.Condition, ", ")
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%.2f\t%.2f\t%d\t%s\n",
			node.Metadata.ID, node.Metadata.Name, node.Status.Hostname,
			node.Status.IP, conditionStr, node.Status.CPUPercent*100, node.Status.MemPercent*100,
			node.Status.NumPods, node.Status.UpdateTime)
	}

	writer.Flush()
}
