package url

import "strings"

// ParseQueryString parses a query string and returns a map of key-value pairs.
//
//nolint:mnd // This function is simple enough.
func ParseQueryString(queryString string) map[string]string {
	params := make(map[string]string)

	pairs := strings.Split(queryString, "&")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			params[kv[0]] = kv[1]
		}
	}

	return params
}
