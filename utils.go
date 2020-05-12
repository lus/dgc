package dgc

import "strings"

// stringHasPrefix checks whether or not the string contains one of the given prefixes and returns the string without the prefix
func stringHasPrefix(str string, prefixes []string, ignoreCase bool) (bool, string) {
	if ignoreCase {
		str = strings.ToLower(str)
	}

	for _, prefix := range prefixes {
		if ignoreCase {
			prefix = strings.ToLower(prefix)
		}
		if strings.HasPrefix(str, prefix) {
			return true, strings.TrimSpace(strings.TrimPrefix(str, prefix))
		}
	}
	return false, ""
}

// equals provides a simple method to check whether or not 2 strings are equal
func equals(str1, str2 string, ignoreCase bool) bool {
	if !ignoreCase {
		return str1 == str2
	}
	return strings.ToLower(str1) == strings.ToLower(str2)
}
