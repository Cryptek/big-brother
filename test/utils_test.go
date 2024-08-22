package test

import (
	"big-brother/internal/models"
	"big-brother/internal/utils"
	"testing"
)

// Test for ValidateConfigAndBuildDependencyTree
func TestValidateConfigAndBuildDependencyTree(t *testing.T) {
	// Valid configuration with dependencies
	validConfig := &models.Config{
		Services: []models.Service{
			{Name: "service1"},
			{Name: "service2", DependsOn: "service1"},
		},
	}

	err := utils.ValidateConfigAndBuildDependencyTree(validConfig)
	if err != nil {
		t.Errorf("Valid configuration should not return an error, but got: %v", err)
	}

	// Invalid configuration with cyclic dependencies
	invalidConfig := &models.Config{
		Services: []models.Service{
			{Name: "service1", DependsOn: "service2"},
			{Name: "service2", DependsOn: "service1"},
		},
	}

	err = utils.ValidateConfigAndBuildDependencyTree(invalidConfig)
	if err == nil {
		t.Error("Invalid configuration with cyclic dependencies should return an error")
	}

	// Invalid configuration with duplicate service names
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

	// Test with a config that has dependencies
	dependencyConfig := &models.Config{
		Services: []models.Service{
			{Name: "service1", DependsOn: "service2"},
			{Name: "service2"},
		},
	}

	err = utils.ValidateConfigAndBuildDependencyTree(dependencyConfig)
	if err != nil {
		t.Errorf("ValidateConfigAndBuildDependencyTree failed for dependency config: %v", err)
	}

	// Check if the dependency tree is constructed correctly
	expectedTree := []*models.Service{
		{Name: "service2", Dependents: []*models.Service{{Name: "service1"}}},
	}
	//TODO: Fix equals checking. Why is it failing?
	if !(&dependencyConfig.DependencyTree == &expectedTree) {
		t.Errorf("Incorrect dependency tree construction.\nExpected: %+v\nGot: %+v", expectedTree, dependencyConfig.DependencyTree)
	}

}

// Test for GetRootNodes
func TestGetRootNodes(t *testing.T) {
	service1 := &models.Service{Name: "service1"}
	service2 := &models.Service{Name: "service2"}
	service3 := &models.Service{Name: "service3"}

	service2.Dependencies = append(service2.Dependencies, service1)
	service3.Dependencies = append(service3.Dependencies, service1)
	service1.Dependents = append(service1.Dependents, service2, service3)

	services := []*models.Service{service1, service2, service3}

	rootNodes := utils.GetRootNodes(services)

	if len(rootNodes) != 1 || rootNodes[0].Name != "service1" {
		t.Errorf("Expected one root node 'service1', but got: %v", getServiceNames(rootNodes))
	}

	// Additional test case with multiple root nodes
	service2.Dependencies = nil // Remove service1 as a dependency
	service3.Dependencies = []*models.Service{service1}
	service1.Dependents = []*models.Service{service3}

	rootNodes = utils.GetRootNodes(services)

	if len(rootNodes) != 2 || !containsService(rootNodes, "service1") || !containsService(rootNodes, "service2") {
		t.Errorf("Expected root nodes 'service1' and 'service2', but got: %v", getServiceNames(rootNodes))
	}
}

// Helper function to extract service names from a slice of services
func getServiceNames(services []*models.Service) []string {
	names := make([]string, len(services))
	for i, service := range services {
		names[i] = service.Name
	}
	return names
}

// Helper function to check if a service name is in the slice of services
func containsService(services []*models.Service, name string) bool {
	for _, service := range services {
		if service.Name == name {
			return true
		}
	}
	return false
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
