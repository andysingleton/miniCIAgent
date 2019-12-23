package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"path"
)

func BuildContainer(filename string, image_name string) {
	cmdStr := fmt.Sprintf("docker build . -f %s -t \"%s\"", filename, image_name)
	out, err := exec.Command("/bin/sh", "-c", cmdStr).Output()
	check(err)

	fmt.Printf("%s", out)
}

func LaunchContainer(image_name string, workflow_name string, manifest_file string) {
	containerName := fmt.Sprintf("miniCI-%s-%d", workflow_name, rand.Intn(1000))
	binaryPath := getQualifiedFilename("miniCI")
	manifestPath := getQualifiedFilename(manifest_file)
	strippedManifestName := path.Base(manifest_file)
	workflowPath := getQualifiedFilename(workflow_name)

	cmdStr := fmt.Sprintf("/usr/bin/docker run --rm -t "+
		"--name %s "+
		"-v %s:/app/miniCI:ro "+
		"-v %s:/app/%s:ro "+
		"-v %s:/app/%s "+
		"%s /app/miniCI -subprocess=yes",
		containerName, binaryPath, manifestPath, strippedManifestName, workflowPath, workflow_name, image_name)

	out, err := exec.Command("/bin/bash", "-c", cmdStr).Output()
	check(err)
	fmt.Printf("%s", out)
}
