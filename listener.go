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

type AgentState struct {
	Ip          string `json:"ip"`
	ExecutionId string
	State       string   `json:"state"`
	Building    string   `json:"building"`
	Done        []string `json:"done"`
	// { "name": type }
	Artefacts []Artefact `json:"artefacts"`
}

type Artefact struct {
	Name string
	Type string
}

type NetworkIpGetter interface {
	Get(backend string) (string, error)
}

type ipNetworkGetter string

func (ipNetworkGetter) Get(backend string) (string, error) {
	switch backend {
	case "local":
		return "172.0.0.1", nil
	case "docker":
		return "172.17.0.1", nil
	}
	return "", errors.New(fmt.Sprintf("Could not handle backend of %s", backend))
}

func initState(ipGetter NetworkIpGetter, backend string, executionId string) {
	var err error

	mutex.Lock()
	state.Ip, err = ipGetter.Get(backend)
	check(err)
	state.ExecutionId = executionId
	state.State = "Starting"
	mutex.Unlock()
}

func setState(newState string) {
	mutex.Lock()
	state.State = newState
	mutex.Unlock()
}

func setBuilding(workflowName string) {
	mutex.Lock()
	state.Building = workflowName
	mutex.Unlock()
}

func addDone(workflowName string) {
	mutex.Lock()
	state.Done = append(state.Done, workflowName)
	mutex.Unlock()
}

func addArtefact(artefact Artefact) {
	mutex.Lock()
	state.Artefacts = append(state.Artefacts, artefact)
	mutex.Unlock()
}

func listener(backend string, executionId uuid.UUID, webPort int) {
	fmt.Println(executionId, ": Starting Listener")

	var ipGetter ipNetworkGetter
	initState(ipGetter, backend, fmt.Sprintf("%s", executionId))
	fmt.Println("State is currently", state.State)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		output, err := json.Marshal(state)
		check(err)
		fmt.Println("State is currently", string(output))
		//json, err := fmt.Sprintf("Returning %s", string(output))
		//check(err)
		fmt.Fprintf(w, string(output))
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", webPort), nil))
}
