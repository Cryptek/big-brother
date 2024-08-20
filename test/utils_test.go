package test

import (
	"big-brother/internal/models"
	"big-brother/internal/utils"
	"testing"
)

// ... (Include the TestIsCyclic from the previous response)

func TestValidateConfigAndBuildDependencyTree(t *testing.T) {
	// Test with a valid config
	validConfig := &models.Config{
		Services: []models.Service{
			{Name: "service1"},
			{Name: "service2", DependsOn: "service1"},
		},
	}

	err := utils.ValidateConfigAndBuildDependencyTree(validConfig)
	if err != nil {
		t.Errorf("ValidateConfigAndBuildDependencyTree failed for valid config: %v", err)
	}

	// Test with a config containing duplicate service names
	duplicateServiceNameConfig := &models.Config{
		Services: []models.Service{
			{Name: "service1"},
			{Name: "service1"}, // Duplicate
		},
	}

	err = utils.ValidateConfigAndBuildDependencyTree(duplicateServiceNameConfig)
	if err == nil {
		t.Error("ValidateConfigAndBuildDependencyTree should have failed for duplicate service name")
	}

	// ... (Add more tests for other invalid config scenarios and cyclic dependencies)
}

func TestGetRootNodes(t *testing.T) {
	services := []*models.Service{
		{Name: "service1"},
		{Name: "service2", DependsOn: "service1"},
		{Name: "service3", DependsOn: "service1"},
	}

	rootNodes := utils.GetRootNodes(services)

	if len(rootNodes) != 1 || rootNodes[0].Name != "service1" {
		t.Errorf("Expected one root node 'service1', but got: %v", rootNodes)
	}
}

func TestFindServiceByName(t *testing.T) {
	config := &models.Config{
		Services: []models.Service{
			{Name: "service1"},
			{Name: "service2"},
		},
	}

	service, err := utils.FindServiceByName(config, "service1")
	if err != nil {
		t.Errorf("FindServiceByName failed for existing service: %v", err)
	}
	if service.Name != "service1" {
		t.Errorf("Expected service name 'service1', but got: %s", service.Name)
	}

	_, err = utils.FindServiceByName(config, "nonexistent_service")
	if err == nil {
		t.Error("FindServiceByName should have failed for nonexistent service")
	}
}

func TestFindProcessByName(t *testing.T) {
	service := &models.Service{
		Name: "test_service",
		Processes: []models.Process{
			{Name: "process1"},
			{Name: "process2"},
		},
	}

	process, err := utils.FindProcessByName(service, "process1")
	if err != nil {
		t.Errorf("FindProcessByName failed for existing process: %v", err)
	}
	if process.Name != "process1" {
		t.Errorf("Expected process name 'process1', but got: %s", process.Name)
	}

	_, err = utils.FindProcessByName(service, "nonexistent_process")
	if err == nil {
		t.Error("FindProcessByName should have failed for nonexistent process")
	}
}
