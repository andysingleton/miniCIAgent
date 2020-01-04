package main

import (
	"encoding/json"
	"flag"
	"github.com/google/uuid"
	"io/ioutil"
	"miniCIAgent/memberlist"
	"net/http"
	"os"
	"sync"
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

func readPipelineManifest(manifestFile string) (Pipeline, error) {
	jsonFile, err := os.Open(manifestFile)
	if err != nil {
		return Pipeline{}, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return Pipeline{}, err
	}

	var manifest Pipeline
	err = json.Unmarshal(byteValue, &manifest)

	return manifest, err
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
	executionIdString := flag.String("pipeline-id", "new", "Pre-defined ID for this execution")
	manifestFile := flag.String("manifest", "default.json", "Name of the Pipeline manifest to build")
	externalArtefacts := flag.String("artifact", "", "Artefacts provided externally to this execution")
	//parentIp := flag.String("parent-ip", "", "IP address of the process that started us")
	flag.Parse()

	err := setExecutionId(*executionIdString)
	check(err)
	pipelineManifest, err := readPipelineManifest(*manifestFile)
	check(err)
	workflowManager := WorkflowManager{[]Workflow{}, []string{}, []string{*externalArtefacts}}
	err = workflowManager.ReadWorkflows(*manifestFile)
	check(err)

	// Start a gossip cluster
	gossipCluster, err := memberlist.Create(memberlist.DefaultLocalConfig(), executionId)
	check(err)
	agentManager := AgentManager{
		*gossipCluster,
		[]AgentState{},
	}

	// spawn the status webservice
	networkManager := NetworkManager{pipelineManifest.ExecutorBackend, pipelineManifest.WebPort}
	localStateManager := AgentState{}
	go listener(networkManager, &localStateManager)
	time.Sleep(1 * time.Second)

	// Local behaviour is to run all workflows in one process
	if pipelineManifest.ExecutorBackend == "local" {
		var waitGroup sync.WaitGroup
		waitGroup.Add(1)
		go AgentLoop(&agentManager, networkManager, &workflowManager, &localStateManager, &waitGroup)
		waitGroup.Wait()
		os.Exit(0)
	}

	// parentIp determines whether the agent thinks it is a launcher, or a true agent
	//if *parentIp != "" {
	//	_, err := agentStates.Join([]string{*parentIp})
	//	check(err)
	//	AgentLoop(*parentIp, pipelineManifest, workflowDefinition, *agentStates)
	//} else {
	//	myIp, err := getNetworkIp(pipelineManifest.ExecutorBackend)
	//	check(err)
	//	launcherLoop(myIp, executionId, pipelineManifest, workflowDefinition, *agentStates)
	//}
}
