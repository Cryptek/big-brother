package executor

import (
	"big-brother/internal/logger"
	"big-brother/internal/models"
	"big-brother/internal/utils"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Executor struct {
	logger *logger.Logger
}

func NewExecutor(logger *logger.Logger) *Executor {
	return &Executor{logger: logger}
}

func (e *Executor) ExecuteCommand(command string, hostName string) (string, error) {
	// For now, assume all commands are local
	// You'll need to implement remote execution logic if needed
	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing command '%s' on host '%s': %w, output: %s", command, hostName, err, string(output))
	}

	return string(output), nil
}

func (e *Executor) StartService(service *models.Service) error {
	e.logger.Infof("Starting service: %s", service.Name)

	for _, process := range service.Processes {
		e.logger.Infof("Starting process: %s on host: %s", process.Name, process.HostName)
		_, err := e.ExecuteCommand(process.StartCmd, process.HostName)
		if err != nil {
			return err
		}

		// Wait for the service to start
		time.Sleep(time.Duration(e.logger.Config.WaitTime) * time.Second)

		// Check if the process is running
		isRunning, err := e.CheckProcess(service.Name, process.Name)
		if err != nil {
			return err
		}
		if !isRunning {
			return fmt.Errorf("process %s on host %s failed to start", process.Name, process.HostName)
		}
	}

	e.logger.Infof("Service %s started successfully.", service.Name)
	return nil
}

func (e *Executor) StopService(service *models.Service) error {
	e.logger.Infof("Stopping service: %s", service.Name)

	for _, process := range service.Processes {
		e.logger.Infof("Stopping process: %s on host: %s", process.Name, process.HostName)
		_, err := e.ExecuteCommand(process.StopCmd, process.HostName)
		if err != nil {
			return err
		}

		// Wait for the service to stop
		time.Sleep(time.Duration(e.logger.Config.WaitTime) * time.Second)

		// Check if the process is stopped
		isRunning, err := e.CheckProcess(service.Name, process.Name)
		if err != nil {
			return err
		}
		if isRunning {
			return fmt.Errorf("process %s on host %s failed to stop", process.Name, process.HostName)
		}
	}

	e.logger.Infof("Service %s stopped successfully.", service.Name)
	return nil
}

func (e *Executor) CheckService(service *models.Service) []models.CheckResult {
	var results []models.CheckResult

	for _, process := range service.Processes {
		isRunning, err := e.CheckProcess(service.Name, process.Name)
		if err != nil {
			e.logger.Errorf("Error checking process %s on host %s: %v", process.Name, process.HostName, err)
			isRunning = false // Assume not running in case of error
		}

		results = append(results, models.CheckResult{
			ServiceName: service.Name,
			ProcessName: process.Name,
			HostName:    process.HostName,
			IsRunning:   isRunning,
		})
	}

	return results
}

func (e *Executor) CheckProcess(serviceName, processName string) (bool, error) {
	service, err := utils.FindServiceByName(e.logger.Config, serviceName)
	if err != nil {
		return false, err
	}

	process, err := utils.FindProcessByName(service, processName)
	if err != nil {
		return false, err
	}

	output, err := e.ExecuteCommand(process.StatusCmd, process.HostName)
	if err != nil {
		return false, err
	}

	// You'll need to adjust this logic based on the actual output of your status commands
	return strings.Contains(output, "running") || strings.Contains(output, "active"), nil
}
