package test

//
//import (
//	"big-brother/internal/executor"
//	"big-brother/internal/models"
//	"testing"
//)
//
//func TestExecutor_ExecuteCommand(t *testing.T) {
//	mockLog := &mockLogger{}
//	executor := executor.NewExecutor(mockLog)
//
//	// Test executing a valid command (assuming 'echo' exists)
//	output, err := executor.ExecuteCommand("echo hello", "localhost")
//	if err != nil {
//		t.Errorf("ExecuteCommand failed for valid command: %v", err)
//	}
//	if output != "hello\n" {
//		t.Errorf("Unexpected output: %s", output)
//	}
//
//	// Test executing an invalid command
//	_, err = executor.ExecuteCommand("invalid_command", "localhost")
//	if err == nil {
//		t.Error("ExecuteCommand should have failed for invalid command")
//	}
//}
//
//func TestExecutor_StartService(t *testing.T) {
//	mockLog := &mockLogger{}
//	executor := executor.NewExecutor(mockLog)
//
//	service := &models.Service{
//		Name: "test_service",
//		Processes: []models.Process{
//			{
//				Name:     "process1",
//				HostName: "localhost",
//				StartCmd: "echo 'starting process1'",
//			},
//		},
//	}
//
//	err := executor.StartService(service)
//	if err != nil {
//		t.Errorf("StartService failed: %v", err)
//	}
//
//	// ... (Add assertions to verify the mock logger's output and the process status)
//}
//
//// ... (Add more tests for StopService, CheckService, CheckProcess)
