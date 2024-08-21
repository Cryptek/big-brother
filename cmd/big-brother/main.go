package main

import (
	"big-brother/internal/app"
	"big-brother/internal/logger"
	"big-brother/internal/models"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	service := flag.String("s", "", "Service to start/stop/check")
	process := flag.String("p", "", "Process to start/stop/check (only with -s)")
	verbose := flag.Bool("v", false, "Enable verbose logging")
	jsonOutput := flag.Bool("j", false, "Enable JSON output for check")
	configFilePath := flag.String("c", "config/config.yaml", "Config file path")
	ignoreCheck := flag.Bool("ic", false, "Ignore dependency checks")
	threadCount := flag.Int("t", 1, "Number of threads for parallel processing")

	flag.Parse()

	command := flag.Arg(0)
	if command == "" {
		fmt.Println("Usage: big-brother [start|stop|check] [options]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Initialize logger
	logger := logger.NewLogger(*verbose)

	// Create app instance
	app := app.NewApp(*configFilePath, *threadCount, *ignoreCheck, logger)

	switch command {
	case "start":
		if *service == "" {
			app.StartAll()
		} else if *process == "" {
			app.StartService(*service)
		} else {
			app.StartProcess(*service, *process)
		}
	case "stop":
		if *service == "" {
			app.StopAll()
		} else if *process == "" {
			app.StopService(*service)
		} else {
			app.StopProcess(*service, *process)
		}
	case "check":
		if *service == "" {
			result := app.CheckAll()
			if *jsonOutput {
				jsonBytes, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					logger.Fatalf("Error marshaling JSON: %v", err)
				}
				fmt.Println(string(jsonBytes))
			} else {
				printCheckResultTable(result)
			}
		} else if *process == "" {
			result := app.CheckService(*service)
			if *jsonOutput {
				// ... (similar to above, marshal and print JSON)
			} else {
				printCheckResultTable(result)
			}
		} else {
			result := app.CheckProcess(*service, *process)
			if *jsonOutput {
				// ... (similar to above, marshal and print JSON)
			} else {
				printCheckResultTable(result)
			}
		}
	default:
		fmt.Println("Invalid command. Use start, stop, or check.")
		os.Exit(1)
	}
}

func printCheckResultTable(results []models.CheckResult) {

	// Sort the results
	sort.Slice(results, func(i, j int) bool {
		if results[i].ServiceName != results[j].ServiceName {
			return results[i].ServiceName < results[j].ServiceName
		} else if results[i].ProcessName != results[j].ProcessName {
			return results[i].ProcessName < results[j].ProcessName
		} else {
			// Convert boolean status to string for comparison
			statusI := "Not Running"
			if results[i].IsRunning {
				statusI = "Running"
			}
			statusJ := "Not Running"
			if results[j].IsRunning {
				statusJ = "Running"
			}
			return statusI < statusJ
		}
	})

	// Define column widths for better formatting
	const (
		serviceNameWidth = 35
		processNameWidth = 10
		hostNameWidth    = 25
		statusWidth      = 10
	)

	// Print header row
	fmt.Printf("%-*s %-*s %-*s %s\n",
		serviceNameWidth, "Service",
		processNameWidth, "Process",
		hostNameWidth, "Host",
		"Status")
	fmt.Println(strings.Repeat("-", serviceNameWidth+processNameWidth+hostNameWidth+statusWidth+3)) // Separator

	// Print each result row with truncation
	for _, result := range results {
		status := "Not Running"
		if result.IsRunning {
			status = "Running"
		}

		// Truncate values if they exceed column width
		serviceName := truncateString(result.ServiceName, serviceNameWidth)
		processName := truncateString(result.ProcessName, processNameWidth)
		hostName := truncateString(result.HostName, hostNameWidth)

		fmt.Printf("%-*s %-*s %-*s %s\n",
			serviceNameWidth, serviceName,
			processNameWidth, processName,
			hostNameWidth, hostName,
			status)
	}
}

func truncateString(str string, maxWidth int) string {
	if len(str) > maxWidth {
		return str[:maxWidth-3] + "..."
	}
	return str
}
