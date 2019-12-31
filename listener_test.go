package main

import (
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

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

//type stubNetworkManager struct {
//	backend string
//	webPort int
//}
type StubNetworkManager struct{}

func (StubNetworkManager) Get() (string, error) {
	return "10.0.0.1", nil
}

func (StubNetworkManager) AddHandler(stateInterface AgentStateInterface) {
	fmt.Println("thing")
}

func (StubNetworkManager) Listen() {
	fmt.Println("thing")
}

func (StubNetworkManager) Webport() int {
	return 1001
}

func TestInitState(t *testing.T) {
	ipGetter := StubNetworkManager{}
	testStatus := AgentState{}
	testId, err := uuid.Parse("00000000-0000-0000-0000-000000000000")
	check(err)
	expectedState := AgentState{
		Ip:          "10.0.0.1",
		ExecutionId: testId,
		State:       "Starting",
	}
	testStatus.InitState(ipGetter)

	eq := reflect.DeepEqual(expectedState, testStatus)
	if eq == false {
		t.Errorf("Object does not match expected state: %s", testStatus)
	}
}

func TestSetState(t *testing.T) {
	ipGetter := StubNetworkManager{}
	testStatus := AgentState{}
	testId, err := uuid.Parse("00000000-0000-0000-0000-000000000000")
	check(err)
	expectedState := AgentState{
		Ip:          "10.0.0.1",
		ExecutionId: testId,
		State:       "foobar",
	}
	testStatus.InitState(ipGetter)
	testStatus.SetStatus("foobar")

	eq := reflect.DeepEqual(expectedState, testStatus)
	if eq == false {
		t.Errorf("Object does not match expected state: %s", testStatus)
	}
}

func TestSetBuilding(t *testing.T) {
	ipGetter := StubNetworkManager{}
	testStatus := AgentState{}
	testId, err := uuid.Parse("00000000-0000-0000-0000-000000000000")
	check(err)
	expectedState := AgentState{
		Ip:          "10.0.0.1",
		ExecutionId: testId,
		State:       "Starting",
		Building:    "foobar",
	}
	testStatus.InitState(ipGetter)
	testStatus.SetBuilding("foobar")

	eq := reflect.DeepEqual(expectedState, testStatus)
	if eq == false {
		t.Errorf("Object does not match exspected state: %s", testStatus)
	}
}

func TestAddDone(t *testing.T) {
	ipGetter := StubNetworkManager{}
	testStatus := AgentState{}
	testId, err := uuid.Parse("00000000-0000-0000-0000-000000000000")
	check(err)
	expectedState := AgentState{
		Ip:          "10.0.0.1",
		ExecutionId: testId,
		State:       "Starting",
		Done:        []string{"foobar"},
	}
	testStatus.InitState(ipGetter)
	testStatus.AddDone("foobar")

	eq := reflect.DeepEqual(expectedState, testStatus)
	if eq == false {
		t.Errorf("Object does not match exspected state: %s", testStatus)
	}
}

func TestAddArtefact(t *testing.T) {
	ipGetter := StubNetworkManager{}
	testStatus := AgentState{}
	testId, err := uuid.Parse("00000000-0000-0000-0000-000000000000")
	check(err)
	testArtefact := Artefact{
		Name: "foobar",
		Type: "file",
	}
	expectedState := AgentState{
		Ip:          "10.0.0.1",
		ExecutionId: testId,
		State:       "Starting",
		Artefacts:   []Artefact{testArtefact},
	}
	testStatus.InitState(ipGetter)
	testStatus.AddArtefact(testArtefact)

	eq := reflect.DeepEqual(expectedState, testStatus)
	if eq == false {
		t.Errorf("Object does not match expected state: %s", testStatus)
	}
}

//
//func TestGetExecutionId_new(t *testing.T) {
//	_, err := getExecutionId("new")
//	if err != nil {
//		t.Errorf("Function raised an error: %s", err)
//	}
//}
//
//func TestGetExecutionId_parse_fail(t *testing.T) {
//	_, err := getExecutionId("blat")
//	if err == nil {
//		t.Errorf("Function did not raise expected error")
//	}
//}
//
//func TestGetExecutionId_parse_success(t *testing.T) {
//	_, err := getExecutionId("123e4567-e89b-12d3-a456-426655440000")
//	if err != nil {
//		t.Errorf("Function raised an error: %s", err)
//	}
//}
//
//func TestGetPipelineConfig(t *testing.T) {
//	expectedResult := Pipeline{
//		Name:                     "test-pipeline",
//		ExecutorBackend:          "local",
//		Dockerfile: 			  "test",
//		MiniciBinaryPath:         "testpath",
//		WebPort:                  8080,
//	}
//	result, err := getPipelineConfig("test-manifest.json")
//
//	if err != nil {
//		t.Errorf("Function raised an error: %s", err)
//	}
//
//	if result != expectedResult {
//		t.Errorf("Result did not match expected result")
//	}
//}
//
//func TestGetPipelineConfig_fail_nofile(t *testing.T) {
//	_, err := getPipelineConfig("nofile.json")
//
//	if err == nil {
//		t.Errorf("Function did not raise expected error")
//	}
//}
//
//func TestGetPipelineConfig_fail_malformed_json(t *testing.T) {
//	_, err := getPipelineConfig("test-manifest-malformed.json")
//
//	if err == nil {
//		t.Errorf("Function did not raise expected error")
//	}
//}
//
//func TestGetPipelineConfig_fail_bad_types(t *testing.T) {
//	_, err := getPipelineConfig("test-manifest-bad-types.json")
//
//	if err == nil {
//		t.Errorf("Function did not raise expected error")
//	}
//}

// OLD TESTS

//func TestGetManifest(t *testing.T) {
//	manifest_file := "test-manifest.json"
//	expected_result := Manifest{
//		"test-pipeline",
//		"test-workflow",
//		"docker",
//		"Dockerfile",
//	}
//	result := getManifest(manifest_file)
//
//	if result != expected_result {
//		t.Errorf("Manifest did not match expected result")
//	}
//}
//
//func TestGetWorkflow(t *testing.T) {
//	workflow_name := "test-workflow"
//
//	// todo: Must be a correct way of building an array of maps
//	test_builders := []map[string]string{}
//	test_builder := map[string]string{
//		"make": "output.dat",
//	}
//	test_builders = append(test_builders, test_builder)
//
//	expected_result := Workflow{
//		workflow_name,
//		test_builders,
//	}
//	result := getWorkflow(workflow_name)
//
//	eq := reflect.DeepEqual(result, expected_result)
//	if eq == false {
//		t.Errorf("Workflow did not match expected result")
//	}
//}
//
//func TestGetQualifiedFilename_absolute(t *testing.T) {
//	cwd, err := os.Getwd()
//	check(err)
//	testFilename := filepath.Join(cwd, "miniCIAgent_test.go")
//	expected_result := testFilename
//	result := getQualifiedFilename(testFilename)
//
//	if result != expected_result {
//		t.Errorf("Filename was incorrect")
//	}
//}
//
//func TestGetQualifiedFilename_relative(t *testing.T) {
//	cwd, err := os.Getwd()
//	check(err)
//	testFilename := "miniCIAgent_test.go"
//	expected_result := filepath.Join(cwd, testFilename)
//	result := getQualifiedFilename(testFilename)
//
//	if result != expected_result {
//		t.Errorf("Filename was incorrect")
//	}
//}
//
//// todo: How to test an error was generated?
////func TestGetQualifiedFilename_failure(t *testing.T) {
////    cwd, err := os.Getwd()
////    check(err)
////    testFilename := "i-dont-exist.nope"
////    expected_result := filepath.Join(cwd, testFilename)
////    result := getQualifiedFilename(testFilename)
////
////    if result != expected_result {
////        t.Errorf("Filename was incorrect")
////    }
////}
