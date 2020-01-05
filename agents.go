package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	memberlist "github.com/hashicorp/memberlist"
	"io/ioutil"
	"os"
	"time"
)

type AgentStateInterface interface {
	InitState(NetworkManagerInterface)
	SetStatus(string)
	SetBuilding(string)
	AddDone(string)
	AddArtefact(string)
	GetAgentState() AgentState
	SetPendingWorkflow(string)
	PromoteToBuilding(string)
	PromoteToDone(string)
}

type AgentState struct {
	Ip          string
	ExecutionId uuid.UUID
	State       string
	Building    string
	Pending     string
	Done        []string
	Artefacts   []string
}

func (agentState *AgentState) SetPendingWorkflow(workflowName string) {
	agentState.Pending = workflowName
}

func (agentState *AgentState) PromoteToDone(workflowName string) {
	agentState.Done = append(agentState.Done, workflowName)
	agentState.Building = ""
}

func (agentState *AgentState) PromoteToBuilding(workflowName string) {
	agentState.Building = workflowName
	agentState.Pending = ""
}

func (agentState AgentState) GetAgentState() AgentState {
	return agentState
}

func (st *AgentState) InitState(ipGetter NetworkManagerInterface) {
	var err error
	st.Ip, err = ipGetter.Get()
	check(err)
	st.ExecutionId = executionId
	st.State = "Starting"
}

func (st *AgentState) SetStatus(newStatus string) {
	st.State = newStatus
}

func (st *AgentState) SetBuilding(workflowName string) {
	st.Building = workflowName
}

func (st *AgentState) AddDone(workflowName string) {
	st.Done = append(st.Done, workflowName)
}

func (st *AgentState) AddArtefact(artefact string) {
	st.Artefacts = append(st.Artefacts, artefact)
}

type AgentManagerInterface interface {
	GetRemoteState(string) (AgentState, error)
	GetMembers() []*memberlist.Node
	GetStates() []AgentState
	UpdateAgentStates(NetworkManagerInterface)
}

type AgentManager struct {
	gossip      memberlist.Memberlist
	agentStates []AgentState
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

func (handler *AgentManager) UpdateAgentStates(networkManager NetworkManagerInterface) {
	memberList := handler.GetMembers()
	var agents []AgentState

	for member := range memberList {
		var agentState AgentState
		connectionString := fmt.Sprintf("http://%s:%d", memberList[member], networkManager.Webport())
		agentState, err := handler.GetRemoteState(connectionString)
		check(err)
		agents = append(agents, agentState)
	}
	handler.agentStates = agents
}

func AgentLoop(agentHandler AgentManagerInterface, networkManager NetworkManagerInterface,
	workflowManager WorkflowManagerInterface, localStateManager AgentStateInterface) {
	for true {
		agentHandler.UpdateAgentStates(networkManager)
		agents := agentHandler.GetStates()
		workflow, err := workflowManager.GetAvailableWorkflow(agents)
		if err != nil {

			// todo: how do we handle artefacts hosted by the agent?
			fmt.Println(executionId, ": No workflows available. Terminating agent")
			os.Exit(0)
		}

		localStateManager.SetPendingWorkflow(workflow.Name)
		agentHandler.UpdateAgentStates(networkManager)
		if workflowManager.IsWorkflowAvailable(workflow) {
			// Workflow is finally available
			localStateManager.PromoteToBuilding(workflow.Name)
			err := processWorkflow(workflow, localStateManager)
			check(err)
			// todo: Process workflow

		} else {
			localStateManager.SetPendingWorkflow("")
			time.Sleep(1 * time.Second)
		}
	}
}
