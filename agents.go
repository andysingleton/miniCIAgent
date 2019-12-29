package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"miniCIAgent/memberlist"
	"time"
)

func getJson(url string) AgentState {
	r, err := myClient.Get(url)
	check(err)
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	check(err)

	var result AgentState
	json.Unmarshal([]byte(body), &result)
	return result
}

func getAllAgentStates(agentStates memberlist.Memberlist, webPort int) []AgentState {
	memberList := agentStates.Members()
	var agents []AgentState
	for member := range memberList {
		var agentState AgentState
		connectionString := fmt.Sprintf("http://%s:%d", memberList[member], webPort)
		agentState = getJson(connectionString)
		agents = append(agents, agentState)
	}
	return agents
}

func getAvailableAgents(memberWorkStates map[string]map[string]string) int {
	available := int(0)
	for _, value := range memberWorkStates {
		if value["workstate"] == "available" {
			available += 1
		}
	}
	return available
}

func launchAgent(IpAddress string, pipeline Pipeline) {
	//if pipeline.ExecutorBackend == "docker" {
	//	dockerFilePath := getQualifiedFilename(pipeline.ExecutorDockerDockerfile)
	//	BuildContainer(dockerFilePath, pipeline.Name)
	//	LaunchContainer(manifest.PipelineName, manifest.StartingWorkflow, *manifestFile, *workflowId, manifest.MiniciBinaryPath)
	//}

}

func AgentLoop(IpAddress string, pipeline Pipeline, workflows Workflows, agentStates memberlist.Memberlist) {
	for true {
		workflow, err := GetAvailableWorkflow(workflows, agentStates, pipeline.WebPort)
		if err != nil {
			fmt.Printf("No work available. Terminating agent\n")
			break
		}
		fmt.Println("Got workflow", workflow)
		break
		//if agentStates.DoTask(workflowId) == true {
		//	// do work
		//	workflow := getWorkflow(*workflowName)
		//}
		time.Sleep(1 * time.Second)
	}
}
