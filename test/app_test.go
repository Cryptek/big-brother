package test

//
//import (
//	"big-brother/internal/app"
//	"big-brother/internal/logger"
//	"big-brother/internal/models"
//	"errors"
//	"testing"
//)
//
//type mockExecutor struct {
//	startServiceFunc func(*models.Service) error
//	stopServiceFunc  func(*models.Service) error
//	checkServiceFunc func(*models.Service) []models.CheckResult
//	checkProcessFunc func(string, string) (bool, error)
//}
//
//func (m *mockExecutor) StartService(service *models.Service) error {
//	if m.startServiceFunc != nil {
//		return m.startServiceFunc(service)
//	}
//	return nil
//}
//
//func (m *mockExecutor) StopService(service *models.Service) error {
//	if m.stopServiceFunc != nil {
//		return m.stopServiceFunc(service)
//	}
//	return nil
//}
//
//func (m *mockExecutor) CheckService(service *models.Service) []models.CheckResult {
//	if m.checkServiceFunc != nil {
//		return m.checkServiceFunc(service)
//	}
//	return nil
//}
//
//func (m *mockExecutor) CheckProcess(serviceName, processName string) (bool, error) {
//	if m.checkProcessFunc != nil {
//		return m.checkProcessFunc(serviceName, processName)
//	}
//	return false, nil
//}
//
//type mockLogger struct {
//	*logger.Logger // Embed logger.Logger
//	fatalfFunc     func(string, ...interface{})
//}
//
//func (m *mockLogger) Info(msg string) {}
//
//func (m *mockLogger) Infof(format string, v ...interface{}) {}
//
//func (m *mockLogger) Error(msg string) {}
//
//func (m *mockLogger) Errorf(format string, v ...interface{}) {}
//
//func (m *mockLogger) Fatal(msg string) {
//	if m.fatalfFunc != nil {
//		m.fatalfFunc(msg)
//	}
//}
//
//func (m *mockLogger) Fatalf(format string, v ...interface{}) {
//	if m.fatalfFunc != nil {
//		m.fatalfFunc(format, v...)
//	}
//}
//
//func TestApp_StartAll(t *testing.T) {
//	mockExec := &mockExecutor{
//		startServiceFunc: func(service *models.Service) error {
//			// Simulate successful service start
//			return nil
//		},
//	}
//	mockLog := &mockLogger{}
//
//	// Create an App instance using NewApp
//	appInstance := app.NewApp("test_config.yaml", 1, false, mockLog)
//	appInstance.Executor = mockExec // Inject the mock executor
//
//	err := appInstance.StartAll() // Use appInstance to call StartAll
//	if err != nil {
//		t.Errorf("StartAll failed: %v", err)
//	}
//
//	// Add assertions to verify the mock executor's startServiceFunc was called for each service
//	// You may need to track which services were started in the mockExecutor
//}
//
//func TestApp_StopAll(t *testing.T) {
//	// Implement similar to StartAll, ensuring to check that StopService is called
//}
//
//func TestApp_StartService(t *testing.T) {
//	mockExec := &mockExecutor{
//		startServiceFunc: func(service *models.Service) error {
//			if service.Name == "service2" {
//				return errors.New("failed to start service2")
//			}
//			return nil
//		},
//		checkProcessFunc: func(serviceName, processName string) (bool, error) {
//			// Simulate all processes running
//			return true, nil
//		},
//	}
//	mockLog := &mockLogger{
//		fatalfFunc: func(format string, v ...interface{}) {
//			// Capture the fatal error message
//			t.Logf("Fatal error: "+format, v...)
//		},
//	}
//
//	app := app.NewApp("test_config.yaml", 1, false, mockLog)
//	app.Executor = mockExec // Inject the mock executor
//
//	// Test starting a service with dependencies
//	err := app.StartService("service1") // Assuming service1 exists
//	if err != nil {
//		t.Errorf("Expected service1 to start successfully, got error: %v", err)
//	}
//
//	// Test starting a service that fails to start
//	err = app.StartService("service2")
//	if err == nil {
//		t.Error("Expected error when starting service2, but got none")
//	}
//}
//
//// ... (Add more tests for StopService, CheckAll, CheckService, CheckProcess, etc.)
