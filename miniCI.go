package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Workflow struct {
	WorkflowName string `json:"name"`
	// todo: enforce no spaces in workflow name
	Builders []map[string]string `json:"build"`
}

type Manifest struct {
	PipelineName string `json:"name"`
	// todo: enforce no spaces in pipeline name
	StartingWorkflow         string `json:"starting workflow"`
	ExecutorBackend          string `json:"executors"`
	ExecutorDockerDockerfile string `json:"dockerfile"`
}

func getFile(filename string) []uint8 {
	data, err := ioutil.ReadFile(filename)
	check(err)

	return data
}

func getManifest(manifestFile string) Manifest {
	manifest := Manifest{}

	data := getFile(manifestFile)

	err := json.Unmarshal([]byte(data), &manifest)
	check(err)

	return manifest
}

func getWorkflow(workflow_name string) Workflow {
	workflow := Workflow{}
	workflow_file := fmt.Sprintf("%s/default.json", workflow_name)

	data := getFile(workflow_file)

	err := json.Unmarshal([]byte(data), &workflow)
	check(err)

	return workflow
}

func getQualifiedFilename(filename string) string {
	if filepath.IsAbs(filename) == true {
		return filename
	}

	if _, err := os.Stat(filename); err == nil {
		cwd, err := os.Getwd()
		check(err)
		return filepath.Join(cwd, filename)
	}

	return filename
}

func main() {
	var manifest Manifest

	manifest_file := flag.String("manifest", "default.json", "Name of the Pipeline manifest to build")
	manifest = getManifest(*manifest_file)
	workflow_name := flag.String("workflow", manifest.StartingWorkflow, "Workflow to begin")
	sub_process := flag.String("subprocess", "", "Flag denoting we were started as a sub-process")
	flag.Parse()

	fmt.Printf("Sub process flag is %s\n", *sub_process)

	workflow := getWorkflow(*workflow_name)

	fmt.Printf("Manifest name is %s\n", manifest.PipelineName)
	fmt.Printf("Workflow name is %s\n", workflow.WorkflowName)
	fmt.Printf("builder command is %s\n", workflow.Builders[0]["make"])

	if *sub_process == "" {
		fmt.Printf("We ARE NOT a subprocess\n")
		if manifest.ExecutorBackend == "docker" {
			fmt.Printf("We are using Docker\n")

			dockerFilePath := getQualifiedFilename(manifest.ExecutorDockerDockerfile)

			BuildContainer(dockerFilePath, manifest.PipelineName)

			LaunchContainer(manifest.PipelineName, manifest.StartingWorkflow, *manifest_file)
		}
	} else {
		fmt.Printf("We ARE a subprocess\n")

	}
}
