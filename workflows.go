package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Workflows struct {
	Workflows []Workflow `json:"workflows"`
}

type Workflow struct {
	Name         string                   `json:"name"`
	Tags         []string                 `json:"tags"`
	Source       []map[string]string      `json:"source"`
	Wants        []Artefact               `json:"wants"`
	Steps        []map[string]interface{} `json:"steps"`
	StepsIterate []map[string]string      `json:"stepsIterate"`
	Provides     []map[string]string      `json:"provides"`
	State        string
}

type WorkflowManagerInterface interface {
	GetWorkflows()
	GetAvailableWorkflow(AgentHandlerInterface, NetworkManagerInterface) (Workflow, error)
}

type WorkflowManager struct {
	manifest  string
	workflows Workflows
}

func (manager WorkflowManager) GetWorkflows() {
	jsonFile, err := os.Open(manager.manifest)
	check(err)
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	check(err)

	err = json.Unmarshal(byteValue, &manager.workflows)
	check(err)
}

func (manager WorkflowManager) GetAvailableWorkflow(agentHandler AgentHandlerInterface, managerInterface NetworkManagerInterface) (Workflow, error) {
	var agents []AgentState
	agents = getAllAgentStates(agentHandler, managerInterface)
	fmt.Println(executionId, ": Got agent states", agents)
	var unavailableList []string
	var artefacts []Artefact

	// todo: implement something better here

	for agent := range agents {
		unavailableList = append(unavailableList, agents[agent].Done...)
		unavailableList = append(unavailableList, agents[agent].Building)
		artefacts = append(artefacts, agents[agent].Artefacts...)
	}

	for workflow := range manager.workflows.Workflows {
		selectWorkflow := true

		// is it available
		for unavailable := range unavailableList {
			if manager.workflows.Workflows[workflow].Name == unavailableList[unavailable] {
				selectWorkflow = false
			}
			return manager.workflows.Workflows[workflow], nil
		}

		// does it have all its pre-requisites
		for workflowArtefact := range manager.workflows.Workflows[workflow].Wants {
			// todo: is there set theory we could use here?
			artefactOk := false
			for artefact := range artefacts {
				if manager.workflows.Workflows[workflow].Wants[workflowArtefact].Name == artefacts[artefact].Name {
					artefactOk = true
				}
			}
			if artefactOk == false {
				selectWorkflow = false
			}
		}
		if selectWorkflow == true {
			return manager.workflows.Workflows[workflow], nil
		}
	}
	return Workflow{}, errors.New("There are no workflows")
}

func getAvailableWorkflowsAmount(workflows Workflows, memberWorkStates map[string]map[string]string) int {
	fmt.Printf("workflows are %s", workflows)
	//for workflow := range manifest["workflows"] {
	//
	//}
	//return "test workflow", nil
	return 2
}

func getCompletedWorkflowsAmount(workflows Workflows, memberWorkStates map[string]map[string]string) int {
	fmt.Printf("workflows are %s", workflows)
	//for workflow := range manifest["workflows"] {
	//
	//}
	//return "test workflow", nil
	return 2
}
