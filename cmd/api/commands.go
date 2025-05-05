package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Commander struct {
	cont_name      string
	conatiner_list map[string]string
}

func NewCommander() *Commander {
	return &Commander{
		conatiner_list: make(map[string]string),
		cont_name:      "archive",
	}
}

func (cmdr *Commander) GetListOfContainers() {
	command := exec.Command("docker", "ps", "-a", "--format", "{{.Names}}")

	out, err := command.CombinedOutput()

	if err != nil {
		fmt.Printf("Something Went Wrong, Make Sure Your Docker Deamon Is Running \n %s", err.Error())
		os.Exit(0)
	}

	names := strings.Split(string(out), "\n")

	for indx, name := range names {
		cmdr.conatiner_list[name] = fmt.Sprintf("%d", indx)
	}
}

func (cmdr *Commander) InitServer() {
	_, ok := cmdr.conatiner_list[cmdr.cont_name]

	if ok {
		cmd := exec.Command("docker", "start", cmdr.cont_name)
		out, err := cmd.CombinedOutput()

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		fmt.Print(string(out))
		return
	} else {

		command := exec.Command("docker", "run", "-d", "--name", cmdr.cont_name, "ubuntu", "sleep", "infinity")

		out, err := command.CombinedOutput()

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		fmt.Print(string(out))

		command = exec.Command("docker", "exec", cmdr.cont_name, "apt", "update")

		out, err = command.CombinedOutput()

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		fmt.Print(string(out))

		command = exec.Command("docker", "exec", cmdr.cont_name, "apt", "install", "nodejs", "golang", "python3", "-y")

		out, err = command.CombinedOutput()

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		log.Print(string(out))

	}

}

func (cmdr *Commander) RunScript(job_id, ext, dir_path, runner, runner_flags string) {
	go func() {
		filePath := fmt.Sprintf("%s/%s.%s", dir_path, job_id, ext)
		command := exec.Command("docker", "cp", filePath, fmt.Sprintf("%s:/", cmdr.cont_name))

		_, err := command.CombinedOutput()

		if err != nil {
			GlobalJobMapper.ResultSet[job_id] = fmt.Sprintf("%cERR: Error While Creating Script", ERROR)
			GlobalJobMapper.JobSet[job_id] = true
			return
		}

		_ = os.Remove(filePath)

		var out []byte

		if len(runner_flags) <= 0 {
			command = exec.Command("docker", "exec", cmdr.cont_name, runner, fmt.Sprintf("%s.%s", job_id, ext))
			out, err = command.CombinedOutput()
		} else {
			command = exec.Command("docker", "exec", cmdr.cont_name, runner, runner_flags, fmt.Sprintf("%s.%s", job_id, ext))
			out, err = command.CombinedOutput()
		}

		if err != nil {
			GlobalJobMapper.ResultSet[job_id] = fmt.Sprintf("%c %s \n %s", ERROR, out, err)
			GlobalJobMapper.JobSet[job_id] = true
			return
		}

		GlobalJobMapper.ResultSet[job_id] = fmt.Sprintf("%c %s", SUCCESS, string(out))
		GlobalJobMapper.JobSet[job_id] = true

		command = exec.Command("docker", "exec", cmdr.cont_name, "rm", "-rf", fmt.Sprintf("%s.%s", job_id, ext))
		_, _ = command.CombinedOutput()
	}()
}
