package main

import (
	"reflect"
	"testing"
)

// NetworkManager test fixture
var addHandlerCalled = false
var listenCalled = false

type StubNetworkManager struct{}

func (StubNetworkManager) Listen()      { listenCalled = true }
func (StubNetworkManager) Webport() int { return 1001 }
func (StubNetworkManager) Get() (string, error) {
	return "10.0.0.1", nil
}
func (StubNetworkManager) AddHandler(stateInterface AgentStateInterface) {
	addHandlerCalled = true
}

// AgentState test fixture
var initStateCalled = false

type StubAgentState struct{}

func (StubAgentState) InitState(NetworkManagerInterface) { initStateCalled = true }
func (StubAgentState) SetStatus(string)                  {}
func (StubAgentState) SetBuilding(string)                {}
func (StubAgentState) AddDone(string)                    {}
func (StubAgentState) AddArtefact(string)                {}
func (StubAgentState) PromoteToBuilding(string)          {}
func (StubAgentState) PromoteToDone(string)              {}
func (StubAgentState) SetPendingWorkflow(string)         {}
func (StubAgentState) GetAgentState() AgentState {
	return AgentState{
		State: "foobar",
	}
}

func TestGetNetworkIp_local(t *testing.T) {
	expectedResult := "172.0.0.1"
	result, err := NetworkManager{"local", 8080}.Get()

	if err != nil {
		t.Errorf("Function raised an error: %s", err)
	}
	if result != expectedResult {
		t.Errorf("Did not return expected result: %s", result)
	}
}

func TestGetNetworkIp_dockerl(t *testing.T) {
	expectedResult := "172.17.0.1"
	result, err := NetworkManager{"docker", 8080}.Get()

	if err != nil {
		t.Errorf("Function raised an error: %s", err)
	}
	if result != expectedResult {
		t.Errorf("Did not return expected result: %s", result)
	}
}

func TestGetNetworkIp_fail(t *testing.T) {
	_, err := NetworkManager{"foobar", 8080}.Get()

	if err == nil {
		t.Errorf("Function did not raise expected error")
	}
}

func TestInitState(t *testing.T) {
	ipGetter := StubNetworkManager{}
	testStatus := AgentState{}
	testStatus.InitState(ipGetter)
	// Since executionId is random, we have to copy it for the DeepEqual
	expectedState := AgentState{
		Ip:          "10.0.0.1",
		ExecutionId: testStatus.ExecutionId,
		State:       "Starting",
	}

	eq := reflect.DeepEqual(expectedState, testStatus)
	if eq == false {
		t.Errorf("Object does not match expected state: %s", testStatus)
	}
}

func TestSetState(t *testing.T) {
	ipGetter := StubNetworkManager{}
	testStatus := AgentState{}
	testStatus.InitState(ipGetter)
	testStatus.SetStatus("foobar")

	expectedState := AgentState{
		Ip:          "10.0.0.1",
		ExecutionId: testStatus.ExecutionId,
		State:       "foobar",
	}

	eq := reflect.DeepEqual(expectedState, testStatus)
	if eq == false {
		t.Errorf("Object does not match expected state: %s", testStatus)
	}
}

func TestSetBuilding(t *testing.T) {
	ipGetter := StubNetworkManager{}
	testStatus := AgentState{}
	testStatus.InitState(ipGetter)
	testStatus.SetBuilding("foobar")

	expectedState := AgentState{
		Ip:          "10.0.0.1",
		ExecutionId: testStatus.ExecutionId,
		State:       "Starting",
		Building:    "foobar",
	}

	eq := reflect.DeepEqual(expectedState, testStatus)
	if eq == false {
		t.Errorf("Object does not match exspected state: %s", testStatus)
	}
}

func TestAddDone(t *testing.T) {
	ipGetter := StubNetworkManager{}
	testStatus := AgentState{}
	testStatus.InitState(ipGetter)
	testStatus.AddDone("foobar")

	expectedState := AgentState{
		Ip:          "10.0.0.1",
		ExecutionId: testStatus.ExecutionId,
		State:       "Starting",
		Done:        []string{"foobar"},
	}

	eq := reflect.DeepEqual(expectedState, testStatus)
	if eq == false {
		t.Errorf("Object does not match exspected state: %s", testStatus)
	}
}

func TestAddArtefact(t *testing.T) {
	ipGetter := StubNetworkManager{}
	testStatus := AgentState{}
	testStatus.InitState(ipGetter)
	testArtefact := "foobar"
	testStatus.AddArtefact(testArtefact)

	expectedState := AgentState{
		Ip:          "10.0.0.1",
		ExecutionId: testStatus.ExecutionId,
		State:       "Starting",
		Artefacts:   []string{testArtefact},
	}

	eq := reflect.DeepEqual(expectedState, testStatus)
	if eq == false {
		t.Errorf("Object does not match expected state: %s", testStatus)
	}
}

func TestGetAgentState(t *testing.T) {
	testAgentState := AgentState{}
	testAgentState.State = "test"
	result := testAgentState.GetAgentState()

	if result.State != testAgentState.State {
		t.Errorf("Resulting object does not match passed object")
	}

}

func TestWebport(t *testing.T) {
	expectedResult := 8080
	testNetManager := NetworkManager{webPort: expectedResult}

	result := testNetManager.Webport()
	if result != expectedResult {
		t.Errorf("Webport is not correct: %d", result)
	}

}

func TestListener(t *testing.T) {
	initStateCalled = false
	addHandlerCalled = false
	listenCalled = false
	ipGetter := StubNetworkManager{}
	stateManager := StubAgentState{}

	listener(ipGetter, stateManager)
	if initStateCalled == false {
		t.Errorf("InitState was not called")
	}
	if addHandlerCalled == false {
		t.Errorf("AddHandler was not called")
	}
	if listenCalled == false {
		t.Errorf("Listen was not called")
	}
}
