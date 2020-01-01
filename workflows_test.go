package main

import (
	"testing"
)

// WorkflowManager test fixture
var getWorkflowsCalled = false

type StubWorkflowManager struct{}

func (StubWorkflowManager) getWorkflows() { getWorkflowsCalled = true }
func (StubWorkflowManager) getAvailableWorkflow(AgentHandlerInterface, NetworkManagerInterface) (Workflow, error) {
	return Workflow{
		Name: "foobar",
	}, nil
}

func TestGetWorkflows(t *testing.T) {
	manager := WorkflowManager{manifest: "test-manifest.json"}
	manager.readWorkflows()

	workflowSuccess := false
	for workflow := range manager.Workflows {
		if manager.Workflows[workflow].Name == "Building" {
			workflowSuccess = true
		}
	}
	if workflowSuccess == false {
		t.Errorf("Workflow structure did not contain a test workflow")
	}
}

func TestGetAvailableWorkflow_success(t *testing.T) {
	agentHandler := StubAgentManager{}
	networkManager := StubNetworkManager{}
	agentHandler.UpdateAgentStates(networkManager)
	manager := WorkflowManager{manifest: "test-manifest.json"}
	manager.readWorkflows()
	manager.updateAvailableWorkflows(&agentHandler, networkManager)
	workflow, _ := manager.getAvailableWorkflow()
	if workflow.Name != "Building" {
		t.Errorf("Available workflow was not returned")
	}
}

func TestGetAvailableWorkflow_fail_availability(t *testing.T) {
	agentHandler := StubAgentManagerBusy{}
	networkManager := StubNetworkManager{}
	manager := WorkflowManager{manifest: "test-manifest.json"}
	manager.readWorkflows()
	manager.updateAvailableWorkflows(&agentHandler, networkManager)
	workflow, err := manager.getAvailableWorkflow()

	if err == nil {
		t.Errorf("No workflows should have been returned %s", workflow)
	}
}

func TestGetAvailableWorkflow_fail_no_artefacts(t *testing.T) {
	agentHandler := StubAgentManagerNoArtefacts{}
	networkManager := StubNetworkManager{}
	manager := WorkflowManager{manifest: "test-manifest.json"}
	manager.readWorkflows()
	manager.updateAvailableWorkflows(&agentHandler, networkManager)
	workflow, err := manager.getAvailableWorkflow()

	if err == nil {
		t.Errorf("No workflows should have been returned %s", workflow)
	}
}

func TestGetAvailableWorkflowsAmount(t *testing.T) {
	expectedResult := 1
	agentHandler := StubAgentManager{}
	networkManager := StubNetworkManager{}
	manager := WorkflowManager{manifest: "test-manifest.json"}
	manager.readWorkflows()
	manager.updateAvailableWorkflows(&agentHandler, networkManager)
	result := manager.getAvailableWorkflowNumber()
	if result != expectedResult {
		t.Errorf("The wrong number of workflows was returned: %d", result)
	}
}
