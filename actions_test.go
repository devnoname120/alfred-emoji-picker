package main

import (
	"bytes"
	"encoding/xml"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Exercise the actual binary and the scripts shipped to Alfred, including their
// exit status and stdout contract. All history is confined to temporary paths.
func TestWorkflowActions(t *testing.T) {
	dir := t.TempDir()
	binary := filepath.Join(dir, "alfred-emoji-picker")
	if output, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, output)
	}
	env := []string{}
	for _, value := range os.Environ() {
		if !strings.HasPrefix(value, "alfred_") {
			env = append(env, value)
		}
	}
	run := func(dataDir string, args ...string) (string, string, error) {
		t.Helper()
		cmd := exec.Command(binary, args...)
		cmd.Env = append(append([]string{}, env...), "alfred_workflow_data="+dataDir)
		var stdout, stderr bytes.Buffer
		cmd.Stdout, cmd.Stderr = &stdout, &stderr
		err := cmd.Run()
		return stdout.String(), stderr.String(), err
	}

	t.Run("terminal record and reset", func(t *testing.T) {
		dataDir := filepath.Join(t.TempDir(), "history")
		stdout, stderr, err := run(dataDir, "--record", "🙂")
		if err != nil || stdout != "" || stderr != "" {
			t.Fatalf("record: stdout=%q stderr=%q err=%v", stdout, stderr, err)
		}
		statsPath := filepath.Join(dataDir, "emoji-usage.json")
		data, err := os.ReadFile(statsPath)
		if err != nil || !bytes.Contains(data, []byte("🙂")) {
			t.Fatalf("recorded history: %s, err=%v", data, err)
		}
		// Reset also recovers corrupt history and is safe to repeat.
		if err := os.WriteFile(statsPath, []byte("{"), 0o644); err != nil {
			t.Fatal(err)
		}
		for i := 0; i < 2; i++ {
			stdout, stderr, err = run(dataDir, "--reset-frequent")
			if err != nil || stdout != "" || stderr != "" {
				t.Fatalf("reset: stdout=%q stderr=%q err=%v", stdout, stderr, err)
			}
			if _, err := os.Stat(statsPath); !os.IsNotExist(err) {
				t.Fatalf("history remains after reset: %v", err)
			}
		}
	})

	t.Run("action errors use stderr", func(t *testing.T) {
		for _, args := range [][]string{{"--record", "🙂"}, {"--reset-frequent"}, {"--record"}, {"--reset-frequent", "extra"}} {
			stdout, stderr, err := run("", args...)
			if err == nil || stdout != "" || stderr == "" || strings.Contains(stderr, "panic:") {
				t.Fatalf("%v: stdout=%q stderr=%q err=%v", args, stdout, stderr, err)
			}
		}
	})

	scripts := recordActionScripts(t)
	t.Run("scripts suppress helper stdout", func(t *testing.T) {
		helperDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(helperDir, "alfred-emoji-picker"), []byte("#!/bin/bash\nprintf 'unexpected helper output'\nexit 1\n"), 0o755); err != nil {
			t.Fatal(err)
		}
		for i, script := range scripts {
			cmd := exec.Command("/bin/bash", "-c", script, "action", "🙂")
			cmd.Dir = helperDir
			cmd.Env = env
			if stdout, err := cmd.Output(); err != nil || string(stdout) != "🙂" {
				t.Fatalf("script %d: stdout=%q err=%v", i, stdout, err)
			}
		}
	})
	for _, scenario := range []string{"valid", "malformed history", "unwritable storage"} {
		t.Run(scenario, func(t *testing.T) {
			dataDir := t.TempDir()
			switch scenario {
			case "malformed history":
				if err := os.WriteFile(filepath.Join(dataDir, "emoji-usage.json"), []byte("{"), 0o644); err != nil {
					t.Fatal(err)
				}
			case "unwritable storage":
				// A regular file cannot be used as a directory, even as root.
				dataDir = filepath.Join(dataDir, "file")
				if err := os.WriteFile(dataDir, nil, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			for i, script := range scripts {
				cmd := exec.Command("/bin/bash", "-c", script, "action", "🙂")
				cmd.Dir = dir
				cmd.Env = append(append([]string{}, env...),
					"alfred_workflow_data="+dataDir,
					"alfred_workflow_bundleid=com.github.devnoname120.alfred-emoji-picker",
					"alfred_workflow_cache="+t.TempDir())
				var stdout, stderr bytes.Buffer
				cmd.Stdout, cmd.Stderr = &stdout, &stderr
				if err := cmd.Run(); err != nil || stdout.String() != "🙂" {
					t.Fatalf("script %d: stdout=%q stderr=%q err=%v", i, stdout.String(), stderr.String(), err)
				}
				if scenario != "valid" && stderr.Len() == 0 {
					t.Fatalf("script %d: missing error diagnostic", i)
				}
			}
			if scenario != "valid" {
				stdout, stderr, err := run(dataDir, "--record", "🙂")
				if err == nil || stdout != "" || stderr == "" {
					t.Fatalf("record: stdout=%q stderr=%q err=%v", stdout, stderr, err)
				}
			}
		})
	}
}

func recordActionScripts(t *testing.T) []string {
	t.Helper()
	data, err := os.ReadFile("info.plist")
	if err != nil {
		t.Fatal(err)
	}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var scripts []string
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if start, ok := token.(xml.StartElement); ok && start.Name.Local == "string" {
			var value string
			if err := decoder.DecodeElement(&value, &start); err != nil {
				t.Fatal(err)
			}
			if strings.Contains(value, "./alfred-emoji-picker --record") {
				scripts = append(scripts, value)
			}
		}
	}
	if len(scripts) != 2 {
		t.Fatalf("expected paste and copy record scripts, got %d", len(scripts))
	}
	return scripts
}
