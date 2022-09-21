package env

import (
	"os"
)

func Environment() string {
	if env := os.Getenv("ENV"); len(env) > 0 {
		return env
	}

	return "dev"
}
