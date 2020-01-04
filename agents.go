package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"miniCIAgent/memberlist"
	"os"
	"sync"
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

func (st AgentState) GetAgentState() AgentState {
	return st
}

func (st *AgentState) InitState(ipGetter NetworkManagerInterface) {
	var err error

	mutex.Lock()
	st.Ip, err = ipGetter.Get()
	check(err)
	st.ExecutionId = executionId
	st.State = "Starting"
	mutex.Unlock()
}

func (st *AgentState) SetStatus(newStatus string) {
	mutex.Lock()
	st.State = newStatus
	mutex.Unlock()
}

func (st *AgentState) SetBuilding(workflowName string) {
	mutex.Lock()
	st.Building = workflowName
	mutex.Unlock()
}

func (st *AgentState) AddDone(workflowName string) {
	mutex.Lock()
	st.Done = append(st.Done, workflowName)
	mutex.Unlock()
}

func (st *AgentState) AddArtefact(artefact string) {
	mutex.Lock()
	st.Artefacts = append(st.Artefacts, artefact)
	mutex.Unlock()
}

func listener(networkManager NetworkManagerInterface, stateManager AgentStateInterface) {
	fmt.Println(executionId, ": Starting Listener")

	stateManager.InitState(networkManager)

	networkManager.AddHandler(stateManager)
	networkManager.Listen()
}

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

func AgentLoop(agentHandler AgentHandlerInterface, networkManager NetworkManagerInterface,
	workflowManager WorkflowManagerInterface, localStateManager AgentStateInterface, waitgGroup *sync.WaitGroup) {
	for true {
		agentHandler.UpdateAgentStates(networkManager)
		agents := agentHandler.GetStates()
		workflow, err := workflowManager.GetAvailableWorkflow(agents)
		if err != nil {
			// terminate here
			// todo: how do we handle artefacts hosted by the agent?
			fmt.Println(executionId, ": No workflows available. Terminating agent")
			defer waitgGroup.Done()
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
