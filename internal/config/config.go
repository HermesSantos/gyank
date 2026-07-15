package config

import "time"

// Config holds application settings. Load expands this later (file/env);
// for now Default is the single source of truth.
type Config struct {
	AppName      string
	PollInterval time.Duration
	WindowWidth  float32
	WindowHeight float32
}

func Default() Config {
	return Config{
		AppName:      "Gyank",
		PollInterval: 500 * time.Millisecond,
		WindowWidth:  400,
		WindowHeight: 200,
	}
}

// Load returns the active configuration.
// Persistence (file/env) can be added here without changing callers.
func Load() Config {
	return Default()
}
