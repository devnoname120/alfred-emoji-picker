package usage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"testing"
	"time"
)

func TestIncrementNormalizesEquivalentEmojiVariants(t *testing.T) {
	t.Setenv("alfred_workflow_data", t.TempDir())

	if err := Increment("🖐️"); err != nil {
		t.Fatalf("increment variant with VS16: %v", err)
	}
	if err := Increment("🖐"); err != nil {
		t.Fatalf("increment variant without VS16: %v", err)
	}

	stats, err := Load()
	if err != nil {
		t.Fatalf("load stats: %v", err)
	}

	if got := stats.Count("🖐️"); got != 2 {
		t.Fatalf("expected normalized count 2, got %d", got)
	}

	if got := stats[NormalizeEmoji("🖐️")]; got != 2 {
		t.Fatalf("expected stored normalized count 2, got %d", got)
	}
}

func TestLoadMissingFileReturnsEmptyStats(t *testing.T) {
	t.Setenv("alfred_workflow_data", t.TempDir())

	stats, err := Load()
	if err != nil {
		t.Fatalf("load stats: %v", err)
	}

	if len(stats) != 0 {
		t.Fatalf("expected empty stats, got %d entries", len(stats))
	}
}

func TestSaveCreatesExpectedFile(t *testing.T) {
	dataDir := t.TempDir()
	t.Setenv("alfred_workflow_data", dataDir)

	if err := Save(Stats{"🙂": 3}); err != nil {
		t.Fatalf("save stats: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dataDir, statsFileName)); err != nil {
		t.Fatalf("stat saved file: %v", err)
	}
}

func TestResetRemovesStatsFile(t *testing.T) {
	dataDir := t.TempDir()
	t.Setenv("alfred_workflow_data", dataDir)

	if err := Save(Stats{"🙂": 3}); err != nil {
		t.Fatalf("save stats: %v", err)
	}
	lockBefore, err := os.Stat(filepath.Join(dataDir, statsFileName+".lock"))
	if err != nil {
		t.Fatal(err)
	}

	if err := Reset(); err != nil {
		t.Fatalf("reset stats: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dataDir, statsFileName)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected stats file to be removed, got err=%v", err)
	}
	lockAfter, err := os.Stat(filepath.Join(dataDir, statsFileName+".lock"))
	if err != nil || !os.SameFile(lockBefore, lockAfter) {
		t.Fatalf("reset must preserve the lock file: %v", err)
	}
}

func TestConcurrentProcessesPreserveCountsAndCompleteSnapshots(t *testing.T) {
	dataDir := t.TempDir()
	t.Setenv("alfred_workflow_data", dataDir)
	// A larger snapshot makes partial writes observable to concurrent readers.
	initial := Stats{"🙂": 0}
	for i := 0; i < 1000; i++ {
		initial[fmt.Sprintf("entry-%04d", i)] = i
	}
	if err := Save(initial); err != nil {
		t.Fatal(err)
	}

	const workers, increments = 6, 40
	start := filepath.Join(dataDir, "start")
	t.Setenv("TEST_USAGE_START", start)
	type process struct {
		done   <-chan error
		output *bytes.Buffer
	}
	processes := make([]process, workers)
	for i := range processes {
		processes[i].done, processes[i].output = startUsageProcess(t, "increment", increments)
	}

	stop := make(chan struct{})
	readerDone := make(chan error, 1)
	go func() {
		for {
			stats, err := Load()
			if err != nil {
				readerDone <- fmt.Errorf("read concurrent snapshot: %w", err)
				return
			}
			if len(stats) != len(initial) || stats["entry-0999"] != 999 {
				readerDone <- fmt.Errorf("incomplete snapshot: %d entries", len(stats))
				return
			}
			select {
			case <-stop:
				readerDone <- nil
				return
			default:
			}
		}
	}()
	t.Cleanup(func() {
		close(stop)
		if err := <-readerDone; err != nil {
			t.Error(err)
		}
	})
	if err := os.WriteFile(start, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	for _, p := range processes {
		if err := <-p.done; err != nil {
			t.Fatalf("increment process: %v\n%s", err, p.output)
		}
	}
	stats, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if got := stats.Count("🙂"); got != workers*increments {
		t.Fatalf("lost increments: got %d, want %d", got, workers*increments)
	}
}

func TestMutationsWaitForSharedProcessLock(t *testing.T) {
	for _, operation := range []string{"increment", "save", "reset"} {
		t.Run(operation, func(t *testing.T) {
			dataDir := t.TempDir()
			t.Setenv("alfred_workflow_data", dataDir)
			if err := Save(Stats{"🙂": 3}); err != nil {
				t.Fatal(err)
			}
			lock, err := os.OpenFile(filepath.Join(dataDir, statsFileName+".lock"), os.O_RDWR, 0)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := lock.Close(); err != nil {
					t.Errorf("close lock file: %v", err)
				}
			}()
			if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
				t.Fatal(err)
			}
			ready := filepath.Join(dataDir, "ready")
			t.Setenv("TEST_USAGE_READY", ready)
			done, output := startUsageProcess(t, operation, 1)
			deadline := time.Now().Add(5 * time.Second)
			for {
				if _, err := os.Stat(ready); err == nil {
					break
				}
				if time.Now().After(deadline) {
					t.Fatal("child process did not become ready")
				}
				time.Sleep(time.Millisecond)
			}
			select {
			case err := <-done:
				t.Fatalf("%s did not wait for the lock: %v\n%s", operation, err, output)
			case <-time.After(100 * time.Millisecond):
			}
			stats, err := Load()
			if err != nil || stats.Count("🙂") != 3 {
				t.Fatalf("stats changed while locked: %v, %v", stats, err)
			}
			if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_UN); err != nil {
				t.Fatal(err)
			}
			if err := <-done; err != nil {
				t.Fatalf("%s process: %v\n%s", operation, err, output)
			}
			stats, err = Load()
			want := map[string]int{"increment": 4, "save": 7, "reset": 0}[operation]
			if err != nil || stats.Count("🙂") != want {
				t.Fatalf("stats after %s: %v, %v; want count %d", operation, stats, err, want)
			}
		})
	}
}

func startUsageProcess(t *testing.T, operation string, count int) (<-chan error, *bytes.Buffer) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	t.Cleanup(cancel)
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestUsageProcess$")
	cmd.Env = append(os.Environ(), "TEST_USAGE_OPERATION="+operation, "TEST_USAGE_COUNT="+strconv.Itoa(count))
	output := &bytes.Buffer{}
	cmd.Stdout, cmd.Stderr = output, output
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	return done, output
}

func TestUsageProcess(t *testing.T) {
	operation := os.Getenv("TEST_USAGE_OPERATION")
	if operation == "" {
		return
	}
	if ready := os.Getenv("TEST_USAGE_READY"); ready != "" {
		if err := os.WriteFile(ready, nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if start := os.Getenv("TEST_USAGE_START"); start != "" {
		for {
			if _, err := os.Stat(start); err == nil {
				break
			}
			time.Sleep(time.Millisecond)
		}
	}
	count, err := strconv.Atoi(os.Getenv("TEST_USAGE_COUNT"))
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < count; i++ {
		switch operation {
		case "increment":
			err = Increment("🙂")
		case "save":
			err = Save(Stats{"🙂": 7})
		case "reset":
			err = Reset()
		default:
			t.Fatalf("unknown operation %q", operation)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
}
