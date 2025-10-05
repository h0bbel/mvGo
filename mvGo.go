package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Rule struct {
	Pattern     string
	Destination string
}

type SyslogConfig struct {
	Enabled bool   `json:"enabled"`
	Network string `json:"network"`
	Address string `json:"address"`
}

type Config struct {
	WatchDir     string       `json:"watch_dir"`
	PollInterval int          `json:"poll_interval"`
	StateFile    string       `json:"state_file"`
	LogFile      string       `json:"log_file"`
	DuplicateDir string       `json:"duplicate_dir"`
	Syslog       SyslogConfig `json:"syslog"`
}

type State struct {
	Processed map[string]struct{} `json:"processed"`
}

var (
	logger *log.Logger
	state  State
	cfg    *Config
)

func loadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	if c.PollInterval <= 0 {
		c.PollInterval = 5
	}
	if c.StateFile == "" {
		c.StateFile = "mvGo.state.json"
	}
	if c.LogFile == "" {
		c.LogFile = "mvGo.log"
	}
	return &c, nil
}

func loadRules(path string) ([]Rule, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	var rules []Rule
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue
		}
		pattern := strings.TrimSpace(parts[0])
		dest := strings.TrimSpace(parts[1])
		rules = append(rules, Rule{Pattern: pattern, Destination: dest})
	}
	return rules, nil
}

func initLogger() {
	logFile, err := os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		os.Exit(1)
	}

	outputs := []io.Writer{os.Stdout, logFile}
	logger = newLogger(outputs, cfg.Syslog)
}

func loadState() {
	state = State{Processed: make(map[string]struct{})}
	data, err := ioutil.ReadFile(cfg.StateFile)
	if err == nil {
		_ = json.Unmarshal(data, &state)
	}
}

func saveState() {
	data, _ := json.MarshalIndent(state, "", "  ")
	_ = ioutil.WriteFile(cfg.StateFile, data, 0644)
}

func matchRule(file string, rules []Rule) *Rule {
	for _, rule := range rules {
		ok, err := filepath.Match(rule.Pattern, filepath.Base(file))
		if err == nil && ok {
			return &rule
		}
	}
	return nil
}

func moveFile(src, dstDir string) error {
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}
	dst := filepath.Join(dstDir, filepath.Base(src))
	return os.Rename(src, dst)
}

func main() {
	configFile := "mvGo.json"
	rulesFile := "mvGo.rules"

	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}
	if len(os.Args) > 2 {
		rulesFile = os.Args[2]
	}

	var err error
	cfg, err = loadConfig(configFile)
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	rules, err := loadRules(rulesFile)
	if err != nil {
		fmt.Println("Error loading rules:", err)
		os.Exit(1)
	}

	initLogger()
	loadState()
	defer saveState()

	logger.Println("Starting mvGo watcher on:", cfg.WatchDir)

	for {
		filepath.WalkDir(cfg.WatchDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}

			base := filepath.Base(path)

			if _, exists := state.Processed[path]; exists {
				if cfg.DuplicateDir != "" {
					dst := filepath.Join(cfg.DuplicateDir, base)
					logger.Printf("File already processed: \"%s\", moving to duplicate destination: \"%s\"\n", base, dst)
					if err := moveFile(path, cfg.DuplicateDir); err != nil {
						logger.Println("Error moving duplicate file:", err)
					}
				} else {
					logger.Printf("File already processed: \"%s\", no duplicate destination configured\n", base)
				}
				return nil
			}

			rule := matchRule(path, rules)
			if rule != nil {
				logger.Printf("Moving \"%s\" -> \"%s\"\n", base, rule.Destination)
				if err := moveFile(path, rule.Destination); err != nil {
					logger.Println("Error moving file:", err)
				} else {
					state.Processed[path] = struct{}{}
					saveState()
				}
			} else {
				logger.Printf("File found but no matching rule: \"%s\"\n", base)
			}

			return nil
		})
		time.Sleep(time.Duration(cfg.PollInterval) * time.Second)
	}
}
