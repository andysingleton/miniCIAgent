package main

import (
	"testing"
)

func TestGetPipelineConfig(t *testing.T) {
	expectedResult := Pipeline{
		Name:             "test-pipeline",
		ExecutorBackend:  "local",
		Dockerfile:       "test",
		MiniciBinaryPath: "testpath",
		WebPort:          8080,
	}
	result, err := readPipelineManifest("test-manifest.json")

	if err != nil {
		t.Errorf("Function raised an error: #{err}")
	}

	if result != expectedResult {
		t.Errorf("Result did not match expected result")
	}
}

func TestGetPipelineConfig_fail_nofile(t *testing.T) {
	_, err := readPipelineManifest("nofile.json")

	if err == nil {
		t.Errorf("Function did not raise expected error")
	}
}

func TestGetPipelineConfig_fail_malformed_json(t *testing.T) {
	_, err := readPipelineManifest("test-manifest-malformed.json")

	if err == nil {
		t.Errorf("Function did not raise expected error")
	}
}

func TestGetPipelineConfig_fail_bad_types(t *testing.T) {
	_, err := readPipelineManifest("test-manifest-bad-types.json")

	if err == nil {
		t.Errorf("Function did not raise expected error")
	}
}
