package utils

import (
	"big-brother/internal/models"
	"errors"
	"fmt"
)

func ValidateConfigAndBuildDependencyTree(cfg *models.Config) error {
	// 1. Validate config (check for duplicate service names, etc.)
	if err := validateConfig(cfg); err != nil {
		return err
	}

	// 2. Create dependency graph
	graph, err := createDependencyGraph(cfg)
	if err != nil {
		return err
	}

	// 3. Check for cyclic dependencies
	if isCyclic(graph) {
		return errors.New("cyclic dependency detected in config")
	}

	// 4. Perform topological sort to get the dependency tree
	sortedServices, err := topologicalSort(graph, cfg)
	if err != nil {
		return err
	}

	// 5. Update the config with the dependency tree
	cfg.DependencyTree = sortedServices

	// 6. Populate Dependents and Dependencies for each service
	for _, service := range cfg.DependencyTree {
		for _, dependencyName := range graph[service.Name] {
			dependencyService, _ := FindServiceByName(cfg, dependencyName)
			service.Dependents = append(service.Dependents, dependencyService)
			dependencyService.Dependencies = append(dependencyService.Dependencies, service)
		}
	}

	return nil
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

	// ... (Add more validation checks as needed)

	return nil
}

func createDependencyGraph(cfg *models.Config) (map[string][]string, error) {
	graph := make(map[string][]string)
	for _, service := range cfg.Services {
		graph[service.Name] = []string{}
		if service.DependsOn != "" {
			graph[service.Name] = append(graph[service.Name], service.DependsOn)
		}
	}
	return graph, nil
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

func topologicalSort(graph map[string][]string, cfg *models.Config) ([]*models.Service, error) { // Pass cfg here
	visited := make(map[string]bool)
	var stack []*models.Service

	for node := range graph {
		if !visited[node] {
			if err := topologicalSortUtil(graph, node, visited, &stack, cfg); err != nil {
				return nil, err
			}
		}
	}

	// Reverse the stack to get the correct order
	for i, j := 0, len(stack)-1; i < j; i, j = i+1, j-1 {
		stack[i], stack[j] = stack[j], stack[i]
	}

	return stack, nil
}

func topologicalSortUtil(graph map[string][]string, node string, visited map[string]bool, stack *[]*models.Service, cfg *models.Config) error {
	visited[node] = true

	for _, neighbor := range graph[node] {
		if !visited[neighbor] {
			if err := topologicalSortUtil(graph, neighbor, visited, stack, cfg); err != nil {
				return err
			}
		}
	}

	// Find the service in the original config
	for _, service := range cfg.Services {
		if service.Name == node {
			*stack = append(*stack, &service)
			break
		}
	}

	return nil
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
		if isLast {
			if i == len(services)-1 {
				branchSymbol = "└───"
			} else {
				branchSymbol = "├───"
			}
		} else {
			if i == len(services)-1 {
				branchSymbol = "└──"
			} else {
				branchSymbol = "├──"
			}
		}

		fmt.Printf("%s%s %s\n", prefix, branchSymbol, service.Name)

		newPrefix := prefix
		if isLast {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}

		PrintDependencyTree(service.Dependents, newPrefix, i == len(services)-1)
	}
}
