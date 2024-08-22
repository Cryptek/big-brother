package utils

import (
	"big-brother/internal/models"
	"errors"
	"fmt"
)

func ValidateConfigAndBuildDependencyTree(cfg *models.Config) error {
	if err := validateConfig(cfg); err != nil {
		return err
	}

	graph, err := createDependencyGraph(cfg)
	if err != nil {
		return err
	}

	if isCyclic(graph) {
		return errors.New("cyclic dependency detected in config")
	}

	sortedServices, err := topologicalSort(graph, cfg)
	if err != nil {
		return err
	}

	// Clear the DependencyTree and only add root nodes
	var rootNodes []*models.Service
	for _, service := range sortedServices {
		if service.DependsOn == "" {
			rootNodes = append(rootNodes, service)
		}
	}
	cfg.DependencyTree = rootNodes

	// Populate Dependents and Dependencies for each service
	for _, service := range cfg.DependencyTree {
		populateDependents(service, sortedServices)
	}

	return nil
}

func populateDependents(parent *models.Service, services []*models.Service) {
	for _, service := range services {
		if service.DependsOn == parent.Name {
			parent.Dependents = append(parent.Dependents, service)
			populateDependents(service, services) // Recursively populate the child nodes
		}
	}
}

func validateConfig(cfg *models.Config) error {
	// Check for duplicate service names
	serviceNames := make(map[string]bool)
	for _, service := range cfg.Services {
		if _, exists := serviceNames[service.Name]; exists {
			return fmt.Errorf("duplicate service name: %s", service.Name)
		}
		serviceNames[service.Name] = true

		// Check for duplicate process names within a service
		processNames := make(map[string]bool)
		for _, process := range service.Processes {
			if _, exists := processNames[process.Name]; exists {
				return fmt.Errorf("duplicate process name: %s in service: %s", process.Name, service.Name)
			}
			processNames[process.Name] = true
		}
	}
	return nil
}

func createDependencyGraph(cfg *models.Config) (map[string][]string, error) {
	graph := make(map[string][]string)
	for _, service := range cfg.Services {
		if _, exists := graph[service.Name]; !exists {
			graph[service.Name] = []string{}
		}
		if service.DependsOn != "" {
			graph[service.DependsOn] = append(graph[service.DependsOn], service.Name)
		}
	}
	return graph, nil
}

func topologicalSort(graph map[string][]string, cfg *models.Config) ([]*models.Service, error) {
	visited := make(map[string]bool)
	includedInTree := make(map[string]bool)
	var stack []*models.Service

	for node := range graph {
		if !visited[node] {
			if err := topologicalSortUtil(graph, node, visited, includedInTree, &stack, cfg); err != nil {
				return nil, err
			}
		}
	}
	return stack, nil
}

func topologicalSortUtil(graph map[string][]string, node string, visited, includedInTree map[string]bool, stack *[]*models.Service, cfg *models.Config) error {
	visited[node] = true

	// Add the current node's dependents first
	for _, dependentName := range graph[node] {
		if !visited[dependentName] {
			if err := topologicalSortUtil(graph, dependentName, visited, includedInTree, stack, cfg); err != nil {
				return err
			}
		}
	}

	// Only add to stack if not already included
	if !includedInTree[node] {
		for _, service := range cfg.Services {
			if service.Name == node {
				*stack = append(*stack, &service)
				includedInTree[node] = true
				break
			}
		}
	}

	return nil
}

func isCyclic(graph map[string][]string) bool {
	visited := make(map[string]bool)
	recursionStack := make(map[string]bool)

	for node := range graph {
		if !visited[node] {
			if isCyclicUtil(graph, node, visited, recursionStack) {
				return true
			}
		}
	}
	return false
}

func isCyclicUtil(graph map[string][]string, node string, visited, recursionStack map[string]bool) bool {
	visited[node] = true
	recursionStack[node] = true

	for _, neighbor := range graph[node] {
		if !visited[neighbor] {
			if isCyclicUtil(graph, neighbor, visited, recursionStack) {
				return true
			}
		} else if recursionStack[neighbor] {
			return true
		}
	}

	recursionStack[node] = false
	return false
}

func GetRootNodes(services []*models.Service) []*models.Service {
	var rootNodes []*models.Service
	for _, service := range services {
		if len(service.Dependencies) == 0 {
			rootNodes = append(rootNodes, service)
		}
	}
	return rootNodes
}

func GetLeafNodes(services []*models.Service) []*models.Service {
	var leafNodes []*models.Service
	for _, service := range services {
		if len(service.Dependents) == 0 {
			leafNodes = append(leafNodes, service)
		}
	}
	return leafNodes
}

func FindServiceByName(cfg *models.Config, serviceName string) (*models.Service, error) {
	for _, service := range cfg.Services {
		if service.Name == serviceName {
			return &service, nil
		}
	}
	return nil, fmt.Errorf("service not found: %s", serviceName)
}

func FindProcessByName(service *models.Service, processName string) (*models.Process, error) {
	for _, process := range service.Processes {
		if process.Name == processName {
			return &process, nil
		}
	}
	return nil, fmt.Errorf("process not found: %s in service: %s", processName, service.Name)
}

func PrintDependencyTree(services []*models.Service, prefix string, isLast bool) {
	for i, service := range services {
		var branchSymbol string
		if i == len(services)-1 {
			branchSymbol = "└───"
		} else {
			branchSymbol = "├───"
		}

		fmt.Printf("%s%s %s\n", prefix, branchSymbol, service.Name)

		newPrefix := prefix
		if isLast {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}

		PrintDependencyTree(service.Dependents, newPrefix, true)
	}
}
