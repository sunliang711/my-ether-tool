package utils

import "regexp"

var re = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

func IsValidAddress(address string) bool {
	return re.MatchString(address)
}
