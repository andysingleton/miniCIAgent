package main

import (
	"github.com/google/uuid"
	"testing"
)

func TestGetExecutionId_new(t *testing.T) {
	err := setExecutionId("new")
	if err != nil {
		t.Errorf("Function raised an error: #{err}")
	}
}

func TestGetExecutionId_parse_fail(t *testing.T) {
	err := setExecutionId("blat")
	if err == nil {
		t.Errorf("Function did not raise expected error")
	}
}

func TestGetExecutionId_parse_success(t *testing.T) {
	testId := "123e4567-e89b-12d3-a456-426655440000"
	err := setExecutionId(testId)
	if err != nil {
		t.Errorf("Function raised an error: #{err}")
	}
	testUUID, _ := uuid.Parse(testId)
	if testUUID != executionId {
		t.Errorf("Execution ID is incorrect: %s", executionId)
	}
}

func TestGetPipelineConfig(t *testing.T) {
	expectedResult := Pipeline{
		Name:             "test-pipeline",
		ExecutorBackend:  "local",
		Dockerfile:       "test",
		MiniciBinaryPath: "testpath",
		WebPort:          8080,
	}
	result, err := getPipelineConfig("test-manifest.json")

	if err != nil {
		t.Errorf("Function raised an error: #{err}")
	}

	if result != expectedResult {
		t.Errorf("Result did not match expected result")
	}
}

func TestGetPipelineConfig_fail_nofile(t *testing.T) {
	_, err := getPipelineConfig("nofile.json")

	if err == nil {
		t.Errorf("Function did not raise expected error")
	}
}

func TestGetPipelineConfig_fail_malformed_json(t *testing.T) {
	_, err := getPipelineConfig("test-manifest-malformed.json")

	if err == nil {
		t.Errorf("Function did not raise expected error")
	}
}

func TestGetPipelineConfig_fail_bad_types(t *testing.T) {
	_, err := getPipelineConfig("test-manifest-bad-types.json")

	if err == nil {
		t.Errorf("Function did not raise expected error")
	}
}
