package main

import (
	"encoding/json"
	"flag"
	"github.com/google/uuid"
	"io/ioutil"
	"miniCIAgent/memberlist"
	"net/http"
	"os"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

type Pipeline struct {
	Name             string `json:"name"`
	ExecutorBackend  string `json:"executors"`
	Dockerfile       string `json:"dockerfile"`
	MiniciBinaryPath string `json:"agentBinaryPath"`
	WebPort          int    `json:"webPort"`
}

type Workflows struct {
	Workflows []Workflow `json:"workflows"`
}

type Workflow struct {
	Name         string                   `json:"name"`
	Tags         []string                 `json:"tags"`
	Source       []map[string]string      `json:"source"`
	Wants        []Artefact               `json:"wants"`
	Steps        []map[string]interface{} `json:"steps"`
	StepsIterate []map[string]string      `json:"stepsIterate"`
	Provides     []map[string]string      `json:"provides"`
	State        string
}

func getPipelineConfig(manifestFile string) (Pipeline, error) {
	jsonFile, err := os.Open(manifestFile)
	if err != nil {
		return Pipeline{}, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return Pipeline{}, err
	}

	var pipeline Pipeline
	err = json.Unmarshal(byteValue, &pipeline)

	return pipeline, err
}

func getExecutionId(id string) (uuid.UUID, error) {
	var executionId uuid.UUID
	var err error
	if id == "new" {
		executionId = uuid.New()
	} else {
		executionId, err = uuid.Parse(id)
	}
	return executionId, err
}

func main() {
	executionIdString := flag.String("pipeline-id", "new", "Name of the Pipeline manifest to build")
	manifestFile := flag.String("manifest", "default.json", "Name of the Pipeline manifest to build")
	parentIp := flag.String("parent-ip", "", "IP address of the process that started us")
	flag.Parse()

	executionId, err := getExecutionId(*executionIdString)
	check(err)
	pipelineConfig, err := getPipelineConfig(*manifestFile)
	check(err)
	workflowDefinition := getWorkflows(*manifestFile)

	// Start a gossip cluster
	agentStates, err := memberlist.Create(memberlist.DefaultLocalConfig(), executionId)
	check(err)

	// spawn the status webservice
	ipGetter := ipNetworkGetter(pipelineConfig.ExecutorBackend)
	go listener(ipGetter, pipelineConfig.ExecutorBackend, executionId, pipelineConfig.WebPort)
	time.Sleep(1 * time.Second)

	// Local behaviour is to run all workflows in one process
	if pipelineConfig.ExecutorBackend == "local" {
		AgentLoop(*parentIp, pipelineConfig, workflowDefinition, *agentStates)
		os.Exit(0)
	}

	// parentIp determines whether the agent thinks it is a launcher, or a true agent
	//if *parentIp != "" {
	//	_, err := agentStates.Join([]string{*parentIp})
	//	check(err)
	//	AgentLoop(*parentIp, pipelineConfig, workflowDefinition, *agentStates)
	//} else {
	//	myIp, err := getNetworkIp(pipelineConfig.ExecutorBackend)
	//	check(err)
	//	launcherLoop(myIp, executionId, pipelineConfig, workflowDefinition, *agentStates)
	//}
}
