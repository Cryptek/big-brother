package models

type Config struct {
	WaitTime       int       `yaml:"wait_time"`
	Services       []Service `yaml:"services"`
	DependencyTree []*Service
}

type Service struct {
	Name         string    `yaml:"name"`
	DependsOn    string    `yaml:"depends_on"`
	Processes    []Process `yaml:"processes"`
	Dependents   []*Service
	Dependencies []*Service
}

type Process struct {
	Name      string `yaml:"name"`
	HostName  string `yaml:"host_name"`
	StartCmd  string `yaml:"start_cmd"`
	StopCmd   string `yaml:"stop_cmd"`
	StatusCmd string `yaml:"status_cmd"`
}

type CheckResult struct {
	ServiceName string `json:"service_name"`
	ProcessName string `json:"process_name"`
	HostName    string `json:"host_name"`
	IsRunning   bool   `json:"is_running"`
}
