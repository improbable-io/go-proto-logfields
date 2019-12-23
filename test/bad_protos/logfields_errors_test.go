// Copyright (c) Improbable Worlds Ltd, All Rights Reserved

package logfieldstest

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	descriptorPath           = flag.String("descriptor_path", "../../deps/include/google/protobuf/descriptor.proto", "Path to the descriptor.proto file.")
	duplicatePath            = flag.String("duplicate_names_path", "duplicate_logfield_names.proto", "Path to the proto file with duplicate field names uses")
	logfieldPath             = flag.String("logfield_path", "../../logfield.proto", "Path to the logfield proto definition")
	protocPath               = flag.String("protoc", "protoc", "Name / path of the protoc compiler.")
	protocGenGoPath          = flag.String("protoc-gen-go", "", "Name / path of the protoc Go compiler plugin.")
	protocGenGologfieldsPath = flag.String("protoc-gen-gologfields", "", "Name / path of the protoc Go logfields compiler plugin.")
	repeatedPath             = flag.String("repeated_path", "repeated_logfield.proto", "Path to the proto file with repeated logfield uses")
)

func protocPathEnv(t *testing.T) []string {
	var extraPaths []string
	if *protocGenGoPath != "" {
		genGoPath, err := filepath.Abs(*protocGenGoPath)
		require.NoError(t, err)
		extraPaths = append(extraPaths, filepath.Dir(genGoPath))
	}
	if *protocGenGologfieldsPath != "" {
		genGologfieldsPath, err := filepath.Abs(*protocGenGologfieldsPath)
		require.NoError(t, err)
		extraPaths = append(extraPaths, filepath.Dir(genGologfieldsPath))
	}
	extraPaths = append(extraPaths, os.Getenv("PATH"))
	return extraPaths
}

func protocArgs(protoFile string, outputPath string) []string {
	// Bazel will pass multiple paths to the -descriptor_path flag. We only need the first.
	realDescriptorPath := strings.Split(*descriptorPath, " ")[0]
	pathElements := strings.Split(realDescriptorPath, string(filepath.Separator))
	pathElements = pathElements[:len(pathElements)-3]
	return []string{
		"--proto_path=" + filepath.Dir(*logfieldPath),
		"--proto_path=" + filepath.Dir(protoFile),
		"--proto_path=" + strings.Join(pathElements, "/"),
		"--go_out=" + outputPath,
	}
}

func buildPlainProtocCommand(t *testing.T, protoFile string, outputPath string) *exec.Cmd {
	cmd := exec.Command(*protocPath, append(protocArgs(protoFile, outputPath), "./"+protoFile)...)
	cmd.Env = []string{fmt.Sprintf("PATH=%s", strings.Join(protocPathEnv(t), ":"))}
	return cmd
}

func buildLogfieldsProtocCommand(t *testing.T, protoFile string, outputPath string) *exec.Cmd {
	cmd := exec.Command(*protocPath, append(protocArgs(protoFile, outputPath), "--gologfields_out="+outputPath, "./"+protoFile)...)
	cmd.Env = []string{fmt.Sprintf("PATH=%s", strings.Join(protocPathEnv(t), ":"))}
	return cmd
}

func logEnvironment(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)
	t.Logf(
		`Test setup:
- current dir:            %s
- protoc:                 %s
- protoc-gen-go:          %s
- protoc-gen-gologfields: %s

- descriptor.proto: %s
- logfield.proto:   %s

- duplicate_logfield_names.proto: %s
- repeated_logfield:              %s`,
		cwd,
		*protocPath,
		*protocGenGoPath,
		*protocGenGologfieldsPath,
		*descriptorPath,
		*logfieldPath,
		*duplicatePath,
		*repeatedPath,
	)
}

func TestPluginFailsOnFilesWithErrors(t *testing.T) {
	testcases := map[string]string{
		"DuplicateLogfieldNames": *duplicatePath,
		"RepeatedLogField":       *repeatedPath,
	}

	logEnvironment(t)

	tmpDir, err := ioutil.TempDir("", "go-proto-logfields")
	require.NoError(t, err)

	for name, path := range testcases {
		protoFile := path
		t.Run(name, func(t *testing.T) {
			var errBytes bytes.Buffer

			cmd := buildPlainProtocCommand(t, protoFile, tmpDir)
			cmd.Stderr = &errBytes
			out, err := cmd.Output()
			t.Log(cmd.Env)
			t.Logf("%s %s:\n%s", cmd.Path, strings.Join(cmd.Args[1:], " "), out)
			require.Empty(t, errBytes.String())
			require.NoError(t, err)

			cmd = buildLogfieldsProtocCommand(t, protoFile, tmpDir)
			out, err = cmd.Output()
			t.Log(cmd.Env)
			t.Logf("%s %s:\n%s", cmd.Path, strings.Join(cmd.Args[1:], " "), out)
			require.Error(t, err)
		})
	}
}
