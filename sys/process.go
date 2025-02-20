package sys

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type ProcessMap map[ShortID]Process

type Process struct {
	Pid  int
	Port int
	Path string
}

func Fill(old, new ProcessMap) ProcessMap {
	result := make(ProcessMap)
	for id, newProcess := range new {
		oldProcess, existedBefore := old[id]
		if existedBefore && newProcess.Pid == 0 {
			newProcess.Pid = oldProcess.Pid
			newProcess.Port = oldProcess.Port
		}
		result[id] = newProcess
	}
	return result
}

func GracefulShutdown(pm *ProcessMap) {
	serverStop := make(chan os.Signal, 1)
	signal.Notify(serverStop, syscall.SIGTERM)
	signal.Notify(serverStop, syscall.SIGQUIT)
	signal.Notify(serverStop, syscall.SIGINT)
	signal.Notify(serverStop, syscall.SIGKILL)
	<-serverStop
	slog.Debug("Server stopping")
	err := CollectGarbage(*pm, make(ProcessMap), true)
	if err != nil {
		slog.Error("Failed to collect garbage", "error", err)
	}
}

func CollectGarbage(old, new ProcessMap, keepFolders bool) error {
	killProcesses := make([]Process, 0)
	for id, oldProcess := range old {
		newProcess, exists := new[id]
		if !exists || newProcess.Path != oldProcess.Path {
			killProcesses = append(killProcesses, oldProcess)
		}
	}
	var err error
	for _, process := range killProcesses {
		if process.Pid != 0 {
			err = KillProcess(process.Pid)
		}
	}
	if err != nil {
		return err
	}
	if !keepFolders {
		err := ClearStorage(new)
		if err != nil {
			return err
		}
	}
	return nil
}
