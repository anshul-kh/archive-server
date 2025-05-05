package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type LangHandler struct {
	exts         map[string]string
	runner       map[string]string
	runner_flags map[string]string
	dir_path     string
}

func NewLangHandler() *LangHandler {
	return &LangHandler{
		exts: map[string]string{
			"js":     "js",
			"golang": "go",
			"python": "py",
		},
		runner: map[string]string{
			"js":     "node",
			"golang": "go",
			"python": "python3",
		},
		runner_flags: map[string]string{
			"golang": "run",
			"js":     "",
			"python": "",
		},
		dir_path: "scripts",
	}
}

func (langh *LangHandler) CreateFile(job_id, lang, code string) error {
	file_path := filepath.Join(langh.dir_path, fmt.Sprintf("%s.%s", job_id, langh.exts[lang]))

	err := os.MkdirAll(langh.dir_path, 0755)

	if err != nil {
		return err
	}

	err = os.WriteFile(file_path, []byte(code), 0644)

	if err != nil {
		return err
	}

	return nil
}
