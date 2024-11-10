package envconfig

import (
	"os"
	"strconv"
)

func Get(key, defaultString string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultString
	}
	return v
}

func GetInt(key string, defaultInt int) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultInt
	}

	retInt, err := strconv.ParseInt(v, 10, 32)
	if nil != err {
		return defaultInt
	}
	return int(retInt)
}
