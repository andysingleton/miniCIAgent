package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"miniCIAgent/memberlist"
	"time"
)

type AgentHandlerInterface interface {
	GetRemoteState(string) (AgentState, error)
	GetMembers() []*memberlist.Node
	GetStates() []AgentState
	UpdateAgentStates(NetworkManagerInterface)
}

type AgentManager struct {
	gossip      memberlist.Memberlist
	agentStates []AgentState
}

type AgentState struct {
	Ip          string
	ExecutionId uuid.UUID
	State       string
	Building    string
	Done        []string
	// { "name": type }
	Artefacts []string
}

func (handler AgentManager) GetStates() []AgentState {
	return handler.agentStates
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

func (handler *AgentManager) UpdateAgentStates(managerInterface NetworkManagerInterface) {
	memberList := handler.GetMembers()
	var agents []AgentState

	for member := range memberList {
		var agentState AgentState
		connectionString := fmt.Sprintf("http://%s:%d", memberList[member], managerInterface.Webport())
		agentState, err := handler.GetRemoteState(connectionString)
		check(err)
		agents = append(agents, agentState)
	}
	handler.agentStates = agents
}

func AgentLoop(agentHandler AgentHandlerInterface, networkManager NetworkManagerInterface, workflowManager WorkflowManagerInterface) {
	for true {
		workflowManager.updateAvailableWorkflows(agentHandler, networkManager)
		workflow, err := workflowManager.getAvailableWorkflow()
		if err != nil {
			fmt.Println(executionId, ": No work available. Terminating agent")
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
