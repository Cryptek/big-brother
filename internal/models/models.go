package models

import (
	"fmt"
	"strings"
)

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

func (s *Service) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Service(Name=%s", s.Name))
	if s.DependsOn != "" {
		sb.WriteString(fmt.Sprintf(", DependsOn=%s", s.DependsOn))
	}
	if len(s.Dependencies) > 0 {
		var deps []string
		for _, dep := range s.Dependencies {
			deps = append(deps, dep.Name)
		}
		sb.WriteString(fmt.Sprintf(", Dependencies=%v", deps))
	}
	if len(s.Dependents) > 0 {
		var deps []string
		for _, dep := range s.Dependents {
			deps = append(deps, dep.Name)
		}
		sb.WriteString(fmt.Sprintf(", Dependents=%v", deps))
	}
	sb.WriteString(")")
	return sb.String()
}

func (c *Config) String() string {
	var sb strings.Builder
	sb.WriteString("Config(WaitTime=")
	sb.WriteString(fmt.Sprintf("%d", c.WaitTime))
	sb.WriteString(", Services=[")
	for i, service := range c.Services {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(service.String())
	}
	sb.WriteString("], DependencyTree=[")
	for i, service := range c.DependencyTree {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(service.String())
	}
	sb.WriteString("])")
	return sb.String()
}
