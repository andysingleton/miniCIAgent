package main

import (
	"flag"
	"github.com/google/uuid"
	"miniCIAgent/memberlist"
	"net/http"
	"os"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}
var executionId = uuid.New()

func main() {
	manifestFile := flag.String("manifest", "default.json", "Name of the Pipeline manifest to build")
	//parentIp := flag.String("parent-ip", "", "IP address of the process that started us")
	flag.Parse()

	pipelineManifest, err := readPipelineManifest(*manifestFile)
	check(err)
	workflowManager := WorkflowManager{[]Workflow{}, []string{}, []string{}}
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

	AgentLoop(&agentManager, networkManager, &workflowManager, &localStateManager)
	os.Exit(0)
}
