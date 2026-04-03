package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrepareACPCommandStagesPIHomeForWrapper(t *testing.T) {
	commonRoot := t.TempDir()
	etcDir := filepath.Join(commonRoot, "etc")
	if err := os.MkdirAll(etcDir, 0o755); err != nil {
		t.Fatalf("mkdir etc: %v", err)
	}
	if err := os.WriteFile(filepath.Join(etcDir, "pi-settings.xproxy.json"), []byte("{\"defaultModel\":\"openai://gpt-5\"}\n"), 0o644); err != nil {
		t.Fatalf("write settings: %v", err)
	}
	if err := os.WriteFile(filepath.Join(etcDir, "pi-models.xproxy.json"), []byte("{\"providers\":{\"xproxy\":{\"models\":[{\"id\":\"openai://gpt-5\",\"name\":\"OpenAI GPT-5\"}]}}}\n"), 0o644); err != nil {
		t.Fatalf("write models: %v", err)
	}
	wrapperDir := filepath.Join(commonRoot, "pi-container")
	if err := os.MkdirAll(wrapperDir, 0o755); err != nil {
		t.Fatalf("mkdir wrapper dir: %v", err)
	}
	wrapperPath := filepath.Join(wrapperDir, "acp-podman.sh")
	if err := os.WriteFile(wrapperPath, []byte("#!/usr/bin/env bash\n"), 0o755); err != nil {
		t.Fatalf("write wrapper: %v", err)
	}

	prepared, err := prepareACPCommand(wrapperPath, t.TempDir(), []string{"A=B"})
	if err != nil {
		t.Fatalf("prepareACPCommand error = %v", err)
	}
	defer func() {
		if err := prepared.cleanup(); err != nil && !os.IsNotExist(err) {
			t.Fatalf("cleanup error = %v", err)
		}
	}()

	if prepared.sessionACPPath != "/home/user" {
		t.Fatalf("sessionACPPath = %q, want /home/user", prepared.sessionACPPath)
	}
	if prepared.commandPath != wrapperPath {
		t.Fatalf("commandPath = %q, want %q", prepared.commandPath, wrapperPath)
	}
	homeDir := ""
	for _, entry := range prepared.env {
		if strings.HasPrefix(entry, "PI_CONTAINER_HOME_DIR=") {
			homeDir = strings.TrimPrefix(entry, "PI_CONTAINER_HOME_DIR=")
			break
		}
	}
	if homeDir == "" {
		t.Fatal("prepared env missing PI_CONTAINER_HOME_DIR")
	}
	if _, err := os.Stat(filepath.Join(homeDir, ".pi", "agent", "settings.json")); err != nil {
		t.Fatalf("staged settings missing: %v", err)
	}
}
