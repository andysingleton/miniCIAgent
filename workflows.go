package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type WorkflowManagerInterface interface {
	ReadWorkflows(string) error
	GetAvailableWorkflow([]AgentState) (Workflow, error)
	IsWorkflowAvailable(Workflow) bool
	UpdateCompletions([]AgentState)
}

type WorkflowManager struct {
	Workflows       []Workflow `json:"workflows"`
	UnavailableList []string
	Artefacts       []string
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

func (manager *WorkflowManager) ReadWorkflows(manifest string) error {
	jsonFile, err := os.Open(manifest)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &manager)
	if err != nil {
		return err
	}
	return nil
}

func (manager *WorkflowManager) IsWorkflowAvailable(workflow Workflow) bool {
	selectWorkflow := true

	// is it available
	for unavailable := range manager.UnavailableList {
		if workflow.Name == manager.UnavailableList[unavailable] {
			selectWorkflow = false
		}
	}

	// does it have all its pre-requisites
	for workflowArtefact := range workflow.Wants {
		artefactName := workflow.Wants[workflowArtefact]
		artefactOk := false
		for artefact := range manager.Artefacts {
			if manager.Artefacts[artefact] == artefactName {
				artefactOk = true
			}
		}
		if artefactOk == false {
			selectWorkflow = false
		}
	}
	return selectWorkflow
}

func (manager *WorkflowManager) UpdateCompletions(agents []AgentState) {
	var unavailableList []string
	var artefacts []string

	for agent := range agents {
		unavailableList = append(unavailableList, agents[agent].Done...)
		unavailableList = append(unavailableList, agents[agent].Building)
		unavailableList = append(unavailableList, agents[agent].Pending)
		artefacts = append(artefacts, agents[agent].Artefacts...)
	}
	manager.UnavailableList = unavailableList
	manager.Artefacts = artefacts
}

func (manager *WorkflowManager) GetAvailableWorkflow(agents []AgentState) (Workflow, error) {
	manager.UpdateCompletions(agents)
	for workflow := range manager.Workflows {
		selectWorkflow := manager.IsWorkflowAvailable(manager.Workflows[workflow])
		if selectWorkflow == true {
			return manager.Workflows[workflow], nil
		}
	}
	return Workflow{}, errors.New("No available workflows")
}

func processWorkflow(workflow Workflow, localStateManager AgentStateInterface) error {
	// todo: processing workflow
	fmt.Println(executionId, ": Processing workflow: ", workflow.Name)
	time.Sleep(10 * time.Second)
	// Success
	for artefact := range workflow.Provides {
		localStateManager.AddArtefact(workflow.Provides[artefact])
	}
	localStateManager.PromoteToDone(workflow.Name)
	return nil
}
