package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type WorkflowManagerInterface interface {
	readWorkflows()
	updateAvailableWorkflows(AgentHandlerInterface, NetworkManagerInterface)
	getAvailableWorkflowNumber() int
	getAvailableWorkflow() (Workflow, error)
}

type WorkflowManager struct {
	manifest           string
	Workflows          []Workflow `json:"workflows"`
	AvailableWorkflows []Workflow
}

type Workflow struct {
	Name         string                   `json:"name"`
	Tags         []string                 `json:"tags"`
	Source       []string                 `json:"source"`
	Wants        []string                 `json:"wants"`
	Events       []map[string]interface{} `json:"events"`
	Steps        []map[string]interface{} `json:"steps"`
	StepsIterate []map[string]string      `json:"stepsIterate"`
	Provides     []string                 `json:"provides"`
	State        string
}

func (manager *WorkflowManager) readWorkflows() {
	jsonFile, err := os.Open(manager.manifest)
	check(err)
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	check(err)

	err = json.Unmarshal(byteValue, &manager)
	check(err)
}

func (manager *WorkflowManager) updateAvailableWorkflows(agentHandler AgentHandlerInterface, networkManager NetworkManagerInterface) {
	agentHandler.UpdateAgentStates(networkManager)
	agents := agentHandler.GetStates()
	var unavailableList []string
	var artefacts []string

	// todo: implement something better here
	var availableWorkflows []Workflow

	for agent := range agents {
		unavailableList = append(unavailableList, agents[agent].Done...)
		unavailableList = append(unavailableList, agents[agent].Building)
		artefacts = append(artefacts, agents[agent].Artefacts...)
	}

	for workflow := range manager.Workflows {
		selectWorkflow := true

		// is it available
		for unavailable := range unavailableList {
			if manager.Workflows[workflow].Name == unavailableList[unavailable] {
				selectWorkflow = false
			}
		}

		// does it have all its pre-requisites
		for workflowArtefact := range manager.Workflows[workflow].Wants {
			artefactName := manager.Workflows[workflow].Wants[workflowArtefact]
			artefactOk := false
			for artefact := range artefacts {
				if artefacts[artefact] == artefactName {
					artefactOk = true
				}
			}
			if artefactOk == false {
				selectWorkflow = false
			}
		}
		if selectWorkflow == true {
			availableWorkflows = append(availableWorkflows, manager.Workflows[workflow])
		}
	}
	manager.AvailableWorkflows = availableWorkflows
}

func (manager *WorkflowManager) getAvailableWorkflowNumber() int {
	return len(manager.AvailableWorkflows)
}

func (manager *WorkflowManager) getAvailableWorkflow() (Workflow, error) {
	if len(manager.AvailableWorkflows) >= 1 {
		return manager.AvailableWorkflows[0], nil
	}
	return Workflow{}, errors.New("There are no available workflows")
}
