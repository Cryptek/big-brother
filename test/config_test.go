package test

import (
	"big-brother/internal/config"
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Test loading a valid config file
	cfg, err := config.LoadConfig("test_config.yaml")
	fmt.Printf("%+v\n", cfg)
	if err != nil {
		t.Errorf("LoadConfig failed for valid config: %v", err)
	}

	// Add assertions to verify the loaded config
	if cfg.WaitTime != 1 {
		t.Errorf("Expected WaitTime to be 10, but got %d", cfg.WaitTime)
	}
	if len(cfg.Services) != 7 {
		t.Errorf("Expected 2 services, but got %d", len(cfg.Services))
	}
	// ... (Add more assertions as needed)

	// Test loading an invalid config file
	_, err = config.LoadConfig("invalid_config.yaml")
	if err == nil {
		t.Error("LoadConfig should have failed for invalid config")
	}
}
