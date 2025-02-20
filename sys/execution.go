package sys

import (
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"
)

var Shell = "bash"

var IndependentProcessAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}

func init() {
	if runtime.GOOS == "windows" {
		Shell = "cmd"
	}
}

func (s Service) build(path string) error {
	cmd := exec.Command(Shell, s.BuildCmd)
	cmd.Dir = path
	if cmd.Err != nil {
		return cmd.Err
	}
	err := cmd.Run()
	if err != nil {
		return err
	}
	slog.Debug("Built service", "id", s.ID, "path", path)
	return nil
}

func (s Service) run(path string, port int) (error, int) {
	cmd := exec.Command("/bin/sh", "-c", s.RunCmd)
	if cmd.Err != nil {
		return cmd.Err, cmd.Process.Pid
	}
	cmd.SysProcAttr = IndependentProcessAttr
	cmd.Env = append(os.Environ(), "PORT="+strconv.Itoa(port))
	cmd.Dir = path
	err := cmd.Start()
	if err != nil {
		return err, cmd.Process.Pid
	}
	slog.Debug("Started service", "id", s.ID, "path", path, "pid", cmd.Process.Pid, "port", port)
	return nil, cmd.Process.Pid
}

func KillProcess(pid int) error {
	pgid, err := syscall.Getpgid(pid)
	if err != nil {
		return err
	}
	err = syscall.Kill(-pgid, 15)
	slog.Debug("Killed process", "pid", pid, "group", pgid)
	if err != nil {
		return err
	}
	return nil
}
