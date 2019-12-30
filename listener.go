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

type NetworkIpGetter interface {
	Get() (string, error)
}

type ipNetworkGetter string

func (backend ipNetworkGetter) Get() (string, error) {
	switch backend {
	case "local":
		return "172.0.0.1", nil
	case "docker":
		return "172.17.0.1", nil
	}
	return "", errors.New(fmt.Sprintf("Could not handle backend of %s", backend))
}

type LocalStateUpdater interface {
	InitState(ipGetter NetworkIpGetter, backend string, executionId string)
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

func (st AgentState) initState(ipGetter NetworkIpGetter, backend string, executionId uuid.UUID) {
	var err error

	mutex.Lock()
	st.Ip, err = ipGetter.Get()
	check(err)
	st.ExecutionId = executionId
	st.State = "Starting"
	mutex.Unlock()
}

func (st AgentState) setState(newState string) {
	mutex.Lock()
	st.State = newState
	mutex.Unlock()
}

func (st AgentState) setBuilding(workflowName string) {
	mutex.Lock()
	st.Building = workflowName
	mutex.Unlock()
}

func (st AgentState) addDone(workflowName string) {
	mutex.Lock()
	st.Done = append(st.Done, workflowName)
	mutex.Unlock()
}

func (st AgentState) addArtefact(artefact Artefact) {
	mutex.Lock()
	st.Artefacts = append(st.Artefacts, artefact)
	mutex.Unlock()
}

func listener(ipGetter NetworkIpGetter, backend string, executionId uuid.UUID, webPort int) {
	fmt.Println(executionId, ": Starting Listener")

	state.initState(ipGetter, backend, executionId)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		output, err := json.Marshal(state)
		check(err)
		fmt.Fprintf(w, string(output))
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", webPort), nil))
}
