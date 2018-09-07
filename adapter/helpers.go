package adapter

import "os"

func getEnv(env string, def string) string {
	envVal := os.Getenv(env)
	if envVal != "" {
		return envVal
	}
	return def
}
