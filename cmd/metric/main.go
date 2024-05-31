package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/url"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/poller"
)

type Target struct {
	Targets []string `json:"targets"`
	Labels  struct {
		// Env string `json:"env"`
		Job string `json:"job"`
	} `json:"labels"`
}

var filePath string

func routing() {
	result, err := httpRequest.GetRequest(url.RootURL + url.GetMetricPoint)
	if err != nil {
		fmt.Println(err)
		return
	}
	var name_url map[string]string
	err = json.Unmarshal([]byte(result), &name_url)
	if err != nil {
		fmt.Println(err)
		return
	}

	targets := make([]Target, 0)
	for k, v := range name_url {
		var t Target
		t.Labels.Job = k

		t.Targets = make([]string, 0)
		t.Targets = append(t.Targets, v)

		targets = append(targets, t)
	}
	fileContent, err := json.Marshal(targets)
	if err != nil {
		fmt.Println(err)
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = file.Write(fileContent)
	if err != nil {
		fmt.Println(err)
	}
	file.Close()

}

func main() {
	metricFilePath := flag.String("metric", "metric.json", "metricFilePath")
	filePath = *metricFilePath
	poller.PollerStaticPeriod(30*time.Second, routing, true)
}
