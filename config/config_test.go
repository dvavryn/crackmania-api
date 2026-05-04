package config_test

import (
	"config"
	"testing"
)

func TestGetConfig(t *testing.T) {
	_, err := config.GetConfig()
	if err != nil {
		t.Fatal(err.Error())
	}
}
