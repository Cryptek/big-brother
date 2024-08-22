package app

import (
	"big-brother/internal/config"
	"big-brother/internal/executor"
	"big-brother/internal/logger"
	"big-brother/internal/models"
	"big-brother/internal/utils"
	"fmt"
	"sync"
	"time"
)

type App struct {
	config      *models.Config
	Executor    *executor.Executor
	logger      *logger.Logger
	threadCount int
	ignoreCheck bool
}

func NewApp(configFilePath string, threadCount int, ignoreCheck bool, logger *logger.Logger) *App {
	cfg, err := config.LoadConfig(configFilePath)
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}

	// Validate config and build dependency tree
	if err := utils.ValidateConfigAndBuildDependencyTree(cfg); err != nil {
		logger.Fatalf("Config validation or dependency tree building failed: %v", err)
	}

	if logger.Verbose {
		logger.Infof("Constructed Config is : %s", cfg)
		logger.Info("Dependency Tree:")
		utils.PrintDependencyTree(cfg.DependencyTree, "", true)
	}

	return &App{
		config:      cfg,
		Executor:    executor.NewExecutor(logger, cfg.WaitTime),
		logger:      logger,
		threadCount: min(threadCount, 192),
		ignoreCheck: ignoreCheck,
	}
}

func (a *App) StartAll() {
	a.logger.Info("Starting all services...")

	// Get the root nodes of the dependency tree
	rootNodes := utils.GetRootNodes(a.config.DependencyTree)

	if a.threadCount > 1 {
		a.processTreeParallel(rootNodes, func(s *models.Service) error {
			return a.startService(s)
		})
	} else {
		a.processTreeSequential(rootNodes, func(s *models.Service) error {
			return a.startService(s)
		})
	}

	a.logger.Info("All services started successfully.")
}

func (a *App) StopAll() {
	a.logger.Info("Stopping all services...")

	// Get the leaf nodes of the dependency tree (reverse order for stopping)
	leafNodes := utils.GetLeafNodes(a.config.DependencyTree)

	if a.threadCount > 1 {
		a.processTreeParallel(leafNodes, a.stopService)
	} else {
		a.processTreeSequential(leafNodes, a.stopService)
	}

	a.logger.Info("All services stopped successfully.")
}

func (a *App) StartService(serviceName string) {
	a.logger.Infof("Starting service: %s", serviceName)

	service, err := utils.FindServiceByName(a.config, serviceName)
	if err != nil {
		a.logger.Fatalf("Error finding service: %v", err)
	}

	// Check dependencies if ignoreCheck is false
	if !a.ignoreCheck {
		for _, dependency := range service.Dependencies {
			isRunning, err := a.isServiceRunning(dependency.Name)
			if err != nil {
				a.logger.Fatalf("Error checking dependency status: %v", err)
			}
			if !isRunning {
				a.logger.Fatalf("Dependency %s is not running. Cannot start %s.", dependency.Name, service.Name)
			}
		}
	}

	if err := a.Executor.StartService(service); err != nil {
		a.logger.Fatalf("Error starting service: %v", err)
	}
}

func (a *App) startService(service *models.Service) error {
	a.logger.Infof("Starting service: %s", service.Name)

	for _, process := range service.Processes {
		a.logger.Infof("Starting process: %s on host: %s", process.Name, process.HostName)
		_, err := a.Executor.ExecuteCommand(process.StartCmd, process.HostName)
		if err != nil {
			return err
		}

		// Wait for the process to start
		time.Sleep(time.Duration(a.config.WaitTime) * time.Second)

		// Check if the process is running
		isRunning, err := a.Executor.CheckProcess(&process)
		if err != nil {
			return err
		}
		if !isRunning {
			return fmt.Errorf("process %s on host %s failed to start", process.Name, process.HostName)
		}
	}

	a.logger.Infof("Service %s started successfully.", service.Name)
	return nil
}

func (a *App) StopService(serviceName string) {
	a.logger.Infof("Stopping service: %s", serviceName)

	service, err := utils.FindServiceByName(a.config, serviceName)
	if err != nil {
		a.logger.Fatalf("Error finding service: %v", err)
	}

	if err := a.Executor.StopService(service); err != nil {
		a.logger.Fatalf("Error stopping service: %v", err)
	}
}

func (a *App) stopService(service *models.Service) error {
	a.logger.Infof("Stopping service: %s", service.Name)

	for _, process := range service.Processes {
		a.logger.Infof("Stopping process: %s on host: %s", process.Name, process.HostName)
		_, err := a.Executor.ExecuteCommand(process.StopCmd, process.HostName)
		if err != nil {
			return err
		}

		// Wait for the process to stop
		time.Sleep(time.Duration(a.config.WaitTime) * time.Second)

		// Check if the process is stopped
		isRunning, err := a.Executor.CheckProcess(&process)
		if err != nil {
			return err
		}
		if isRunning {
			return fmt.Errorf("process %s on host %s failed to stop", process.Name, process.HostName)
		}
	}

	a.logger.Infof("Service %s stopped successfully.", service.Name)
	return nil
}

func (a *App) CheckAll() []models.CheckResult {
	a.logger.Info("Checking all services...")
	var allResults []models.CheckResult

	for _, service := range a.config.Services {
		results := a.Executor.CheckService(&service)
		allResults = append(allResults, results...)
	}

	return allResults
}

func (a *App) CheckService(serviceName string) []models.CheckResult {
	a.logger.Infof("Checking service: %s", serviceName)

	service, err := utils.FindServiceByName(a.config, serviceName)
	if err != nil {
		a.logger.Fatalf("Error finding service: %v", err)
	}

	return a.Executor.CheckService(service)
}

func (a *App) CheckProcess(serviceName, processName string) []models.CheckResult {
	a.logger.Infof("Checking process: %s in service: %s", processName, serviceName)

	service, err := utils.FindServiceByName(a.config, serviceName)
	if err != nil {
		a.logger.Fatalf("Error finding service: %v", err)
	}

	process, err := utils.FindProcessByName(service, processName)
	if err != nil {
		a.logger.Fatalf("Error finding process: %v", err)
	}

	isRunning, err := a.Executor.CheckProcess(process)
	if err != nil {
		a.logger.Fatalf("Error checking process: %v", err)
	}

	return []models.CheckResult{
		{
			ServiceName: serviceName,
			ProcessName: processName,
			HostName:    process.HostName,
			IsRunning:   isRunning,
		},
	}
}

func (a *App) StartProcess(serviceName, processName string) {
	a.logger.Infof("Starting process: %s in service: %s", processName, serviceName)

	service, err := utils.FindServiceByName(a.config, serviceName)
	if err != nil {
		a.logger.Fatalf("Error finding service: %v", err)
	}

	process, err := utils.FindProcessByName(service, processName)
	if err != nil {
		a.logger.Fatalf("Error finding process: %v", err)
	}

	// Don't wait to check start when starting only individual process
	a.logger.Infof("Starting process: %s on host: %s", process.Name, process.HostName)
	_, err = a.Executor.ExecuteCommand(process.StartCmd, process.HostName)
	if err != nil {
		a.logger.Fatalf("Error starting process: %v", err)
	}
}

func (a *App) StopProcess(serviceName, processName string) {
	a.logger.Infof("Stopping process: %s in service: %s", processName, serviceName)

	service, err := utils.FindServiceByName(a.config, serviceName)
	if err != nil {
		a.logger.Fatalf("Error finding service: %v", err)
	}

	process, err := utils.FindProcessByName(service, processName)
	if err != nil {
		a.logger.Fatalf("Error finding process: %v", err)
	}

	//Don't wait to check stop when stopping only individual process
	a.logger.Infof("Stopping process: %s on host: %s", process.Name, process.HostName)
	_, err = a.Executor.ExecuteCommand(process.StopCmd, process.HostName)
	if err != nil {
		a.logger.Fatalf("Error stopping process: %v", err)
	}
}

func (a *App) processTreeParallel(nodes []*models.Service, action func(*models.Service) error) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, a.threadCount)

	for _, node := range nodes {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(service *models.Service) {
			defer wg.Done()
			defer func() { <-semaphore }()

			if err := action(service); err != nil {
				a.logger.Fatalf("Error processing service %s: %v", service.Name, err)
			}

			a.processTreeParallel(service.Dependents, action)
		}(node)
	}

	wg.Wait()
}

func (a *App) processTreeSequential(nodes []*models.Service, action func(*models.Service) error) {
	for _, node := range nodes {
		if err := action(node); err != nil {
			a.logger.Fatalf("Error processing service %s: %v", node.Name, err)
		}

		a.processTreeSequential(node.Dependents, action)
	}
}

func (a *App) isServiceRunning(serviceName string) (bool, error) {
	results := a.CheckService(serviceName)
	for _, result := range results {
		if result.IsRunning {
			return true, nil
		}
	}
	return false, nil
}
