package config

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

type Config struct {

	TokensFile string `json:"tokens_file"`

	ProxyFile string `json:"proxy_file"`

	LogDir string `json:"log_dir"`

	MaxConcurrent int `json:"max_concurrent"`

	MinBoostDelay time.Duration `json:"min_boost_delay"`
	MaxBoostDelay time.Duration `json:"max_boost_delay"`

	CheckDelay time.Duration `json:"check_delay"`

	MaxRetries int `json:"max_retries"`

	BaseURL string `json:"base_url"`

	UserAgent string `json:"user_agent"`
}

func Default() *Config {
	return &Config{
		TokensFile:    "tokens.txt",
		ProxyFile:     "",
		LogDir:        "logs",
		MaxConcurrent: 10,
		MinBoostDelay: 1500 * time.Millisecond,
		MaxBoostDelay: 3500 * time.Millisecond,
		CheckDelay:    800 * time.Millisecond,
		MaxRetries:    3,
		BaseURL:       "https://discord.com/api/v10",
		UserAgent:     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
}

func (c *Config) UnmarshalJSON(b []byte) error {
	type Alias Config
	aux := &struct {
		MinBoostDelay json.RawMessage `json:"min_boost_delay"`
		MaxBoostDelay json.RawMessage `json:"max_boost_delay"`
		CheckDelay    json.RawMessage `json:"check_delay"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(b, aux); err != nil {
		return err
	}

	if len(aux.MinBoostDelay) > 0 {
		c.MinBoostDelay = parseDurationRaw(aux.MinBoostDelay, 1500*time.Millisecond)
	}
	if len(aux.MaxBoostDelay) > 0 {
		c.MaxBoostDelay = parseDurationRaw(aux.MaxBoostDelay, 3500*time.Millisecond)
	}
	if len(aux.CheckDelay) > 0 {
		c.CheckDelay = parseDurationRaw(aux.CheckDelay, 800*time.Millisecond)
	}

	return nil
}

func (c *Config) MarshalJSON() ([]byte, error) {
	type Alias Config
	return json.Marshal(&struct {
		MinBoostDelay string `json:"min_boost_delay"`
		MaxBoostDelay string `json:"max_boost_delay"`
		CheckDelay    string `json:"check_delay"`
		*Alias
	}{
		MinBoostDelay: c.MinBoostDelay.String(),
		MaxBoostDelay: c.MaxBoostDelay.String(),
		CheckDelay:    c.CheckDelay.String(),
		Alias:         (*Alias)(c),
	})
}

func parseDurationRaw(raw json.RawMessage, defaultVal time.Duration) time.Duration {
	s := strings.TrimSpace(string(raw))
	if s == "" || s == "null" {
		return defaultVal
	}

	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
		str := strings.Trim(s, "\"")
		if d, err := time.ParseDuration(str); err == nil {
			return d
		}
	}

	var num float64
	if err := json.Unmarshal(raw, &num); err == nil {
		if num > 1e7 {

			return time.Duration(num)
		}

		return time.Duration(num) * time.Millisecond
	}

	return defaultVal
}

func FromFile(path string) (*Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
