package template

import (
	"os"
	"runtime"
)

func GenReplaceFunc(m map[string]string) func(string) string {
	return func(key string) string {
		switch key {
		case "os":
			return runtime.GOOS
		case "arch":
			return runtime.GOARCH
		}
		if v, ok := m[key]; ok {
			return v
		}
		return os.Getenv(key)
	}
}
