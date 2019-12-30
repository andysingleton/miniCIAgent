package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"miniCIAgent/memberlist"
	"os"
)

func getWorkflows(manifestFile string) Workflows {
	jsonFile, err := os.Open(manifestFile)
	check(err)
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	check(err)

	var workflows Workflows
	err = json.Unmarshal(byteValue, &workflows)
	check(err)

	return workflows
}

func GetAvailableWorkflow(workflows Workflows, agentStates memberlist.Memberlist, webPort int) (Workflow, error) {
	var agents []AgentState
	agents = getAllAgentStates(agentStates, webPort)
	fmt.Println("Got agents", agents)
	var unavailableList []string
	var artefacts []Artefact

	// todo: implement something better here

	for agent := range agents {
		unavailableList = append(unavailableList, agents[agent].Done...)
		unavailableList = append(unavailableList, agents[agent].Building)
		artefacts = append(artefacts, agents[agent].Artefacts...)
	}

	for workflow := range workflows.Workflows {
		selectWorkflow := true

		// is it available
		for unavailable := range unavailableList {
			if workflows.Workflows[workflow].Name == unavailableList[unavailable] {
				selectWorkflow = false
			}
			return workflows.Workflows[workflow], nil
		}

		// does it have all its pre-requisites
		for workflowArtefact := range workflows.Workflows[workflow].Wants {
			// todo: is there set theory we could use here?
			artefactOk := false
			for artefact := range artefacts {
				if workflows.Workflows[workflow].Wants[workflowArtefact].Name == artefacts[artefact].Name {
					artefactOk = true
				}
			}
			if artefactOk == false {
				selectWorkflow = false
			}
		}
		if selectWorkflow == true {
			return workflows.Workflows[workflow], nil
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
