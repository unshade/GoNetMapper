package internal

import (
	"regexp"
	"strconv"
)

func IsValidIP(ip string) bool {
	regex := `^(\d{1,3}\.){3}\d{1,3}$`
	return regexp.MustCompile(regex).MatchString(ip)
}

func IsValidPort(port string) bool {
	i, err := strconv.Atoi(port)
	if err != nil {
		return false
	}
	return i > 0 && i < 65536
}
