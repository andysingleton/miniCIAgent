package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"miniCIAgent/memberlist"
	"time"
)

type AgentHandlerInterface interface {
	GetRemoteState(string) (AgentState, error)
	GetMembers() []*memberlist.Node
}

type AgentManager struct {
	gossip memberlist.Memberlist
}

func (handler AgentManager) GetRemoteState(url string) (AgentState, error) {
	r, err := myClient.Get(url)
	if err != nil {
		return AgentState{}, err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return AgentState{}, err
	}

	var result AgentState
	err = json.Unmarshal([]byte(body), &result)
	return result, err
}

func (handler AgentManager) GetMembers() []*memberlist.Node {
	return handler.gossip.Members()
}

func getAllAgentStates(handler AgentHandlerInterface, managerInterface NetworkManagerInterface) []AgentState {
	memberList := handler.GetMembers()
	var agents []AgentState
	for member := range memberList {
		var agentState AgentState
		connectionString := fmt.Sprintf("http://%s:%d", memberList[member], managerInterface.Webport())
		agentState, err := handler.GetRemoteState(connectionString)
		check(err)
		agents = append(agents, agentState)
	}
	return agents
}

func getAvailableAgents(memberWorkStates map[string]map[string]string) int {
	available := int(0)
	for _, value := range memberWorkStates {
		if value["workstate"] == "available" {
			available += 1
		}
	}
	return available
}

func launchAgent(IpAddress string, pipeline Pipeline) {
	//if pipeline.ExecutorBackend == "docker" {
	//	dockerFilePath := getQualifiedFilename(pipeline.ExecutorDockerDockerfile)
	//	BuildContainer(dockerFilePath, pipeline.Name)
	//	LaunchContainer(manifest.PipelineName, manifest.StartingWorkflow, *manifestFile, *workflowId, manifest.MiniciBinaryPath)
	//}

}

func AgentLoop(agentHandler AgentHandlerInterface, networkManager NetworkManagerInterface, workflowManager WorkflowManagerInterface) {
	for true {

		workflow, err := workflowManager.GetAvailableWorkflow(agentHandler, networkManager)
		if err != nil {
			fmt.Printf("%s: No work available. Terminating agent\n", executionId)
			break
		}
		fmt.Println("Got workflow", workflow)
		break
		//if agentStates.DoTask(workflowId) == true {
		//	// do work
		//	workflow := getWorkflow(*workflowName)
		//}
		time.Sleep(1 * time.Second)
	}
}
