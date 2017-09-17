package main

import (
	"bytes"
	"github.com/golang/glog"
	"strings"
	"os/exec"
	"errors"
	"fmt"
)

type flyCli struct {
	target string
}

type pipeline struct {
	name   string
	team   string
	paused bool
	public bool
}

func (fly *flyCli) Pipelines() ([]*pipeline, error) {
	cmd := fmt.Sprintf("fly --target=%s pipelines --all", fly.target)
	stdout, _, err := executeCmd(cmd, "")
	if err != nil {
		return nil, err
	}
	return parseFlyPipelinesOutput(stdout), nil
}

// pipeline commands
func (fly *flyCli) GetPipeline(name string) (string, error) {
	cmd := fmt.Sprintf("fly --target=%s get-pipeline --pipeline=%s", fly.target, name)
	out, _, err := executeCmd(cmd, "")
	return out, err
}

func (fly *flyCli) RenamePipeline(oldName string, newName string) error {
	cmd := fmt.Sprintf("fly --target=%s rename-pipeline --old-name=%s --new-name=%s", fly.target, oldName, newName)
	_, _, err := executeCmd(cmd, "")
	return err
}

func (fly *flyCli) SetPipeline(name string, config string) error {
	// todo: figure out a better, crossplatform way to do this (instead of writing to /dev/stdin....)
	cmd := fmt.Sprintf("fly --target=%s set-pipeline --pipeline=%s --non-interactive --config /dev/stdin", fly.target, name)
	_, _, err := executeCmd(cmd, config)
	return err
}

func (fly *flyCli) DestroyPipeline(name string) error {
	cmd := fmt.Sprintf("fly --target=%s destroy-pipeline --pipeline=%s --non-interactive", fly.target, name)
	_, _, err := executeCmd(cmd, "")
	return err
}

func (fly *flyCli) HidePipeline(name string) error {
	cmd := fmt.Sprintf("fly --target=%s hide-pipeline --pipeline=%s", fly.target, name)
	_, _, err := executeCmd(cmd, "")
	return err
}

func (fly *flyCli) ExposePipeline(name string) error {
	cmd := fmt.Sprintf("fly --target=%s expose-pipeline --pipeline=%s", fly.target, name)
	_, _, err := executeCmd(cmd, "")
	return err
}

func (fly *flyCli) PausePipeline(name string) error {
	cmd := fmt.Sprintf("fly --target=%s pause-pipeline --pipeline=%s", fly.target, name)
	_, _, err := executeCmd(cmd, "")
	return err
}

func (fly *flyCli) UnpausePipeline(name string) error {
	cmd := fmt.Sprintf("fly --target=%s unpause-pipeline --pipeline=%s", fly.target, name)
	_, _, err := executeCmd(cmd, "")
	return err
}


func parseFlyPipelinesOutput(pipelineTable string) []*pipeline {
	buf := []*pipeline{}
	for _, row := range strings.Split(pipelineTable, "\n") {
		if row == "" {
			continue // skip empty rows (like beginning/end of output i guess)
		}
		fields := strings.Fields(row)
		if len(fields) != 4 {
			glog.Errorf("Parse pipeline fail, list output does not have 4 fields: %s", row)
			continue
		}
		p := &pipeline{
			name:   fields[0],
			team:   fields[1],
			paused: fields[2] == "yes",
			public: fields[3] == "yes",
		}
		buf = append(buf, p)
	}
	return buf
}

// handy wrapper to execute a CLI command and return the result
func executeCmd(cmdStr string, stdin string) (stdout string, stderr string, err error) {
	args := strings.Split(cmdStr, " ")
	cmd := exec.Command(args[0], args[1:]...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	cmd.Stdin = strings.NewReader(stdin)
	err = cmd.Run()
	stdout = outBuf.String()
	stderr = errBuf.String()
	if err != nil {
		err = errors.New(err.Error() + stderr + stdout)
		glog.ErrorDepth(1, err)
	} else {
		glog.V(4).Infof("Exec success: %s\n%s%s", cmdStr, stdout, stderr)
	}
	return
}
