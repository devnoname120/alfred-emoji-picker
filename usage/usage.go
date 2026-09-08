package usage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const statsFileName = "emoji-usage.json"

type Stats map[string]int

func Load() (Stats, error) {
	path, err := statsFilePath()
	if err != nil {
		return nil, err
	}
	return load(path)
}

func load(path string) (Stats, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Stats{}, nil
	}
	if err != nil {
		return nil, err
	}

	stats := Stats{}
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, err
	}
	if stats == nil {
		stats = Stats{}
	}

	return stats, nil
}

func Increment(emojiChar string) error {
	return withStatsLock(func(path string) error {
		stats, err := load(path)
		if err != nil {
			return err
		}
		stats[NormalizeEmoji(emojiChar)]++
		return save(path, stats)
	})
}

func Save(stats Stats) error {
	return withStatsLock(func(path string) error {
		return save(path, stats)
	})
}

func withStatsLock(update func(string) error) error {
	path, err := statsFilePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	// Keep this file in place: replacing or removing it would let processes
	// lock different inodes while updating the same stats file.
	lock, err := os.OpenFile(path+".lock", os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return err
	}
	defer lock.Close() // Closing the descriptor also releases the flock.
	for {
		err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX)
		if !errors.Is(err, syscall.EINTR) {
			break
		}
	}
	if err != nil {
		return err
	}
	return update(path)
}

func save(path string, stats Stats) error {
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}

	// Rename a complete file from the same directory so unlocked readers see
	// either the previous snapshot or the new one, never partially written JSON.
	tmp, err := os.CreateTemp(filepath.Dir(path), ".emoji-usage-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()
	if err := tmp.Chmod(0o644); err != nil {
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

func Reset() error {
	return withStatsLock(func(path string) error {
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		return nil
	})
}

func (s Stats) Count(emojiChar string) int {
	if s == nil {
		return 0
	}

	return s[NormalizeEmoji(emojiChar)]
}

func statsFilePath() (string, error) {
	dataDir := os.Getenv("alfred_workflow_data")
	if dataDir == "" {
		return "", errors.New("alfred_workflow_data is not set")
	}
	return filepath.Join(dataDir, statsFileName), nil
}

func NormalizeEmoji(e string) string {
	return strings.ReplaceAll(e, "\uFE0F", "")
}
