package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	mutex sync.Mutex
	//state AgentState
)

type NetworkManagerInterface interface {
	Get() (string, error)
	AddHandler(AgentStateInterface)
	Listen()
	Webport() int
}

type NetworkManager struct {
	backend string
	webPort int
}

func (net NetworkManager) Webport() int {
	return net.webPort
}

func (net NetworkManager) Get() (string, error) {
	switch net.backend {
	case "local":
		return "172.0.0.1", nil
	case "docker":
		return "172.17.0.1", nil
	}
	return "", errors.New(fmt.Sprintf("Could not handle backend of %s", net.backend))
}

func (net NetworkManager) AddHandler(state AgentStateInterface) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		output, err := json.Marshal(state.GetAgentState())
		check(err)
		fmt.Fprintf(w, string(output))
	})
}

func (net NetworkManager) Listen() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", net.webPort), nil))
}

type AgentStateInterface interface {
	InitState(NetworkManagerInterface)
	SetStatus(string)
	SetBuilding(string)
	AddDone(string)
	AddArtefact(string)
	GetAgentState() AgentState
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
