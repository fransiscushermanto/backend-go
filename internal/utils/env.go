package utils

import "sync"

var (
	appEnv string
	envMu  sync.RWMutex
)

// SetEnv sets the application environment (e.g., "development", "production").
func SetEnv(env string) {
	envMu.Lock()
	defer envMu.Unlock()
	appEnv = env
}

// GetEnv retrieves the application environment.
func GetEnv() string {
	envMu.RLock()
	defer envMu.RUnlock()
	return appEnv
}

// IsProduction checks if the environment is production.
func IsProduction() bool {
	return GetEnv() == "production"
}

// IsDevelopment checks if the environment is development.
func IsDevelopment() bool {
	return GetEnv() == "development"
}
