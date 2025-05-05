package main

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const (
	PENDING = '#'
	SUCCESS = '+'
	ERROR   = '-'
)

type ExecutionMapper struct {
	JobSet    map[string]bool
	ResultSet map[string]string
}

func NewExecutionMapper() *ExecutionMapper {
	return &ExecutionMapper{
		JobSet:    make(map[string]bool),
		ResultSet: make(map[string]string),
	}
}

func (em *ExecutionMapper) CreateJob(lang, code string) (string, error) {
	job_id := uuid.New().String()
	em.JobSet[job_id] = false

	err := GlobalLangHandler.CreateFile(job_id, lang, code)
	if err != nil {
		return "", err
	}

	GlobalCommander.RunScript(job_id, GlobalLangHandler.exts[lang], GlobalLangHandler.dir_path, GlobalLangHandler.runner[lang], GlobalLangHandler.runner_flags[lang])

	return job_id, nil
}

func (em *ExecutionMapper) UpdateJobStatus(job_id string, status bool) error {
	em.JobSet[job_id] = status
	return nil
}

func (em *ExecutionMapper) GetJobResult(job_id string) (string, error) {
	job_status, ok := em.JobSet[job_id]

	if !ok {
		return string(ERROR), errors.New("invalid Job Id/ Job Id do not exist")
	}

	if !job_status {
		return string(PENDING), nil
	}

	job_result, ok := em.ResultSet[job_id]

	if !ok {
		return string(ERROR), errors.New("invalid Job Id/ Job Id do not exist in result set")
	}

	return fmt.Sprintf("%c %s", SUCCESS, job_result), nil
}

func (em *ExecutionMapper) AppendResult(job_id string, result string) error {
	em.ResultSet[job_id] = result
	em.JobSet[job_id] = true
	return nil
}
