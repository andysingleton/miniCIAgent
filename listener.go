package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"sync"
)

var (
	mutex sync.Mutex
	state AgentState
)

type NetworkManagerInterface interface {
	Get() (string, error)
	AddHandler()
	Listen()
}

type NetworkManager struct {
	backend string
	webPort int
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

func (net NetworkManager) AddHandler() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		output, err := json.Marshal(state)
		check(err)
		fmt.Fprintf(w, string(output))
	})
}

func (net NetworkManager) Listen() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", net.webPort), nil))
}

type LocalStateUpdater interface {
	InitState(ipGetter NetworkManager, backend string, executionId string)
	SetState(newState string)
	SetBuilding(workflowName string)
	AddDone(workflowName string)
	AddArtefact(artefact Artefact)
}

type AgentState struct {
	Ip          string
	ExecutionId uuid.UUID
	State       string
	Building    string
	Done        []string
	// { "name": type }
	Artefacts []Artefact
}

type Artefact struct {
	Name string
	Type string
}

func (st *AgentState) initState(ipGetter NetworkManagerInterface, executionId uuid.UUID) {
	var err error

	mutex.Lock()
	st.Ip, err = ipGetter.Get()
	check(err)
	st.ExecutionId = executionId
	st.State = "Starting"
	mutex.Unlock()
}

func (st *AgentState) setState(newState string) {
	mutex.Lock()
	st.State = newState
	mutex.Unlock()
}

func (st *AgentState) setBuilding(workflowName string) {
	mutex.Lock()
	st.Building = workflowName
	mutex.Unlock()
}

func (st *AgentState) addDone(workflowName string) {
	mutex.Lock()
	st.Done = append(st.Done, workflowName)
	mutex.Unlock()
}

func (st *AgentState) addArtefact(artefact Artefact) {
	mutex.Lock()
	st.Artefacts = append(st.Artefacts, artefact)
	mutex.Unlock()
}

func listener(netManager NetworkManager, executionId uuid.UUID) {
	fmt.Println(executionId, ": Starting Listener")

	state.initState(netManager, executionId)

	netManager.AddHandler()
	netManager.Listen()
}
