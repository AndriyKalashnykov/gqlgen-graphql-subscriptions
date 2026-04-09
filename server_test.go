package main

import (
	"testing"
)

func TestGetEnv_Default(t *testing.T) {
	t.Setenv("TEST_ENV_VAR", "")

	result := getEnv("TEST_ENV_VAR", "default-value")
	if result != "default-value" {
		t.Errorf("expected 'default-value', got %s", result)
	}
}

func TestGetEnv_Set(t *testing.T) {
	t.Setenv("TEST_ENV_VAR", "custom-value")

	result := getEnv("TEST_ENV_VAR", "default-value")
	if result != "custom-value" {
		t.Errorf("expected 'custom-value', got %s", result)
	}
}

func TestGetEnv_Empty(t *testing.T) {
	t.Setenv("TEST_ENV_VAR", "")

	result := getEnv("TEST_ENV_VAR", "default-value")
	if result != "default-value" {
		t.Errorf("expected 'default-value' for empty env, got %s", result)
	}
}
