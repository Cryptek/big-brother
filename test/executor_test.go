package test

import (
	"big-brother/internal/executor"
	"big-brother/internal/logger"
	"big-brother/internal/models"
	"testing"
)

func TestExecutor_ExecuteCommand(t *testing.T) {
	log := logger.NewLogger(false) // Initialize the actual logger
	executor := executor.NewExecutor(log, 1)

	// Test executing a valid command (assuming 'echo' exists)
	output, err := executor.ExecuteCommand("echo hello", "localhost")
	if err != nil {
		t.Errorf("ExecuteCommand failed for valid command: %v", err)
	}
	if output != "hello\n" {
		t.Errorf("Unexpected output: %s", output)
	}

	// Test executing an invalid command
	_, err = executor.ExecuteCommand("invalid_command", "localhost")
	if err == nil {
		t.Error("ExecuteCommand should have failed for invalid command")
	}
}

func TestExecutor_StartService(t *testing.T) {
	log := logger.NewLogger(true) // Initialize the actual logger
	executor := executor.NewExecutor(log, 1)

	service := &models.Service{
		Name: "test_service",
		Processes: []models.Process{
			{
				Name:      "process1",
				HostName:  "localhost",
				StartCmd:  "echo 'starting process1'",
				StopCmd:   "echo 'stopping process1'",
				StatusCmd: "echo 'All Good!'",
			},
		},
	}

	err := executor.StartService(service)
	if err != nil {
		t.Errorf("StartService failed: %v", err)
	}

	// Optionally, add assertions to verify the process status
	// (though this may require more complex setup to verify running processes)
}
