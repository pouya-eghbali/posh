package env

import (
	"os"
)

func GetEnvironmentVariable(name string) string {
	return os.Getenv(name)
}
