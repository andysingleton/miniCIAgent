package main

import (
	"fmt"
	"github.com/google/uuid"
	"miniCIAgent/memberlist"
	"time"
)

func launcherLoop(IpAddress string, executionId uuid.UUID, pipeline Pipeline, workflows Workflows, agentStates memberlist.Memberlist) {
	fmt.Printf("%s: Starting local agent\n", executionId)
	fmt.Printf("%s: Entering pipeline \"%s\"\n", executionId, pipeline.Name)

	for true {
		totalWorkflows := len(workflows.Workflows)
		completedWorkflows := getCompletedWorkflowsAmount(workflows, agentStates.MemberWorkStates())

		if totalWorkflows == completedWorkflows {
			fmt.Printf("%s: All workflows completed\n", executionId)
			break
		}

		// todo: currently stubbed
		availableWorkflows := getAvailableWorkflowsAmount(workflows, agentStates.MemberWorkStates())
		availableAgents := getAvailableAgents(agentStates.MemberWorkStates())

		// If we have more work than agents, launch a new agent
		// -1 is to discount the local agent, which does no work
		if availableWorkflows > (availableAgents - 1) {
			launchAgent(IpAddress, pipeline)
		}

		time.Sleep(1 * time.Second)
		break
	}
}
