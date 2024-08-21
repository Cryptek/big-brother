# Big Brother - Generic Application Monitoring and Management Tool

Big Brother is a Go-based tool for monitoring and managing (start/stop) applications across different environments. It
reads configuration from a YAML file, validates dependencies, and provides commands to start, stop, and check the status
of services and processes.

## Features

* **Configuration-driven:** Define your environment and services in a YAML file.
* **Dependency management:** Automatically handles service dependencies to ensure correct startup and shutdown order.
* **Start/Stop/Check:** Provides commands to start, stop, and check the status of services and processes.
* **Parallel processing:** Optionally process services in parallel for faster execution.
* **JSON output:** Get structured JSON output for the `check` command.
* **Verbose logging:** Enable detailed logging for troubleshooting.
* **Cross-platform:** Builds static binaries for different operating systems and architectures.
* **Testable:** Includes unit and integration tests.

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://your-repo-url/big-brother.git
   cd big-brother
   ```

2. **Build the binary:**

   ```bash
   make build
   ```

   This will create the `big-brother` binary in the `build` directory.

## Usage

```
big-brother [start|stop|check] [options]

Options:

-s, --service string     Service to start/stop/check
-p, --process string     Process to start/stop/check (only with -s)
-v, --verbose            Enable verbose logging
-j, --json               Enable JSON output for check
-c, --config string      Config file path (default "config/config.yaml")
-ic, --ignore-check      Ignore dependency checks
-t, --thread-count int   Number of threads for parallel processing (default 1)
```

**Examples:**

* **Start all services:**

  ```bash
  big-brother start
  ```

* **Stop a specific service:**

  ```bash
  big-brother stop -s service1
  ```

* **Check the status of all services and get JSON output:**

  ```bash
  big-brother check -j
  ```

## Configuration

Create a `config.yaml` file in the `config` directory with the following structure:

```yaml
wait_time: 10  # Wait time in seconds between service check after start/stop
services:
  - name: service1
    depends_on: service2
    processes:
      - name: process1
        host_name: host1
        start_cmd: "command_to_start_process1"
        stop_cmd: "command_to_stop_process1"
        status_cmd: "command_to_check_process1"
  - name: service2
    processes:
      - name: process2
        host_name: host2
        start_cmd: "command_to_start_process2"
        stop_cmd: "command_to_stop_process2"
        status_cmd: "command_to_check_process2"
```
