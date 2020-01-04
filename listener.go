package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// todo: Fix this abstraction
// Its part webservice manager, part network config
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
