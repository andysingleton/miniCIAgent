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
var executionId uuid.UUID

type Pipeline struct {
	Name             string `json:"name"`
	ExecutorBackend  string `json:"executors"`
	Dockerfile       string `json:"dockerfile"`
	MiniciBinaryPath string `json:"agentBinaryPath"`
	WebPort          int    `json:"webPort"`
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

func setExecutionId(id string) error {
	var err error
	if id == "new" {
		executionId = uuid.New()
	} else {
		executionId, err = uuid.Parse(id)
	}
	return err
}

func main() {
	executionIdString := flag.String("pipeline-id", "new", "Name of the Pipeline manifest to build")
	manifestFile := flag.String("manifest", "default.json", "Name of the Pipeline manifest to build")
	//parentIp := flag.String("parent-ip", "", "IP address of the process that started us")
	flag.Parse()

	err := setExecutionId(*executionIdString)
	check(err)
	pipelineConfig, err := getPipelineConfig(*manifestFile)
	check(err)

	workflowManager := WorkflowManager{*manifestFile, Workflows{}}
	workflowManager.GetWorkflows()

	// Start a gossip cluster
	agentStates, err := memberlist.Create(memberlist.DefaultLocalConfig(), executionId)
	check(err)
	agentManager := AgentManager{
		*agentStates,
	}

	// spawn the status webservice
	networkManager := NetworkManager{pipelineConfig.ExecutorBackend, pipelineConfig.WebPort}
	stateManager := AgentState{}
	go listener(networkManager, &stateManager)
	time.Sleep(1 * time.Second)

	// Local behaviour is to run all workflows in one process
	if pipelineConfig.ExecutorBackend == "local" {
		AgentLoop(agentManager, networkManager, workflowManager)
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
