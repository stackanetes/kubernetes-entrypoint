package env

import (
	"os"
	"strings"
)

func SplitEnvToList(env string, s ...string) (envList []string) {
	separator := ","
	if len(s) > 0 {
		separator = s[0]
	}
	e := os.Getenv(env)
	if e == "" {
		return envList
	}

	envList = strings.Split(e, separator)
	return envList
}
