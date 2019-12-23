package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetManifest(t *testing.T) {
	manifest_file := "test-manifest.json"
	expected_result := Manifest{
		"test-pipeline",
		"test-workflow",
		"docker",
		"Dockerfile",
	}
	result := getManifest(manifest_file)

	if result != expected_result {
		t.Errorf("Manifest did not match expected result")
	}
}

func TestGetWorkflow(t *testing.T) {
	workflow_name := "test-workflow"

	// todo: Must be a correct way of building an array of maps
	test_builders := []map[string]string{}
	test_builder := map[string]string{
		"make": "output.dat",
	}
	test_builders = append(test_builders, test_builder)

	expected_result := Workflow{
		workflow_name,
		test_builders,
	}
	result := getWorkflow(workflow_name)

	eq := reflect.DeepEqual(result, expected_result)
	if eq == false {
		t.Errorf("Workflow did not match expected result")
	}
}

func TestGetQualifiedFilename_absolute(t *testing.T) {
	cwd, err := os.Getwd()
	check(err)
	testFilename := filepath.Join(cwd, "miniCI_test.go")
	expected_result := testFilename
	result := getQualifiedFilename(testFilename)

	if result != expected_result {
		t.Errorf("Filename was incorrect")
	}
}

func TestGetQualifiedFilename_relative(t *testing.T) {
	cwd, err := os.Getwd()
	check(err)
	testFilename := "miniCI_test.go"
	expected_result := filepath.Join(cwd, testFilename)
	result := getQualifiedFilename(testFilename)

	if result != expected_result {
		t.Errorf("Filename was incorrect")
	}
}

// todo: How to test an error was generated?
//func TestGetQualifiedFilename_failure(t *testing.T) {
//    cwd, err := os.Getwd()
//    check(err)
//    testFilename := "i-dont-exist.nope"
//    expected_result := filepath.Join(cwd, testFilename)
//    result := getQualifiedFilename(testFilename)
//
//    if result != expected_result {
//        t.Errorf("Filename was incorrect")
//    }
//}
