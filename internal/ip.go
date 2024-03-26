package internal

import "regexp"

func IsValidIP(ip string) bool {
	regex := `^(\d{1,3}\.){3}\d{1,3}$`
	return regexp.MustCompile(regex).MatchString(ip)
}
