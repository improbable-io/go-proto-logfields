// Copyright (c) Improbable Worlds Ltd, All Rights Reserved

package logfieldstest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testOut = "test_out"

func buildProtoPathArgs(protoFile string) []string {
	var args []string
	for _, path := range strings.Split(os.Getenv("GOPATH"), ":") {
		if path == "" {
			continue
		}
		args = append(args, fmt.Sprintf("--proto_path=%v/src", path))
		args = append(args, fmt.Sprintf("--proto_path=%v/src/github.com/improbable-io/go-proto-logfields/", path))
		args = append(args, fmt.Sprintf("--proto_path=%v/src/github.com/improbable-io/go-proto-logfields/vendor/github.com/gogo/protobuf/protobuf", path))
	}
	args = append(args, fmt.Sprintf("--proto_path=%v", path.Dir(protoFile)))
	return args
}

func buildPlainProtocCommand(protoFile string) *exec.Cmd {
	args := buildProtoPathArgs(protoFile)
	args = append(args, "--gogo_out="+testOut)
	args = append(args, protoFile)
	return exec.Command("protoc", args...)
}

func buildLogfieldsProtocCommand(protoFile string) *exec.Cmd {
	args := buildProtoPathArgs(protoFile)
	args = append(args, "--gogo_out="+testOut)
	args = append(args, "--gologfields_out="+testOut)
	args = append(args, protoFile)
	return exec.Command("protoc", args...)
}

func requireSuccess(t *testing.T, cmd *exec.Cmd) {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	assert.Equal(t, "", string(stderr.Bytes()))
	require.NoError(t, err)
}

func buildTestFunc(protoPath string, expectSuccess bool) func(*testing.T) {
	return func(t *testing.T) {
		// Test that protoc ignoring logfields works, i.e. the proto file is valid
		requireSuccess(t, buildPlainProtocCommand(protoPath))
		// Test protoc with the logfields plugin
		if expectSuccess {
			requireSuccess(t, buildLogfieldsProtocCommand(protoPath))
		} else {
			require.Error(t, buildLogfieldsProtocCommand(protoPath).Run())
		}
	}
}

func buildTestFuncsForDir(t *testing.T, dir string, expectSuccess bool) []testing.InternalTest {
	var tests []testing.InternalTest
	files, err := ioutil.ReadDir(dir)
	require.NoError(t, err)
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".proto") {
			continue
		}
		tests = append(tests, testing.InternalTest{
			Name: file.Name(),
			F:    buildTestFunc(path.Join(dir, file.Name()), expectSuccess),
		})
	}
	require.NotEmpty(t, tests)
	return tests
}

func TestExpectedErrors(t *testing.T) {
	os.Mkdir(testOut, 0770)
	// Verify that the commands work for good protos
	tests := buildTestFuncsForDir(t, ".", true)
	// Test that the logfields plugin fails for bad protos
	tests = append(tests, buildTestFuncsForDir(t, "bad_protos", false)...)
	// Run tests
	if !testing.RunTests(func(_, _ string) (bool, error) { return true, nil }, tests) {
		t.Fail()
	}
}
