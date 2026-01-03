package proxy

import (
	"fmt"
	"net/url"
	"regexp"

	"strings"
)

// normalize target string into "host:port"
func ParseTarget(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("target cannot be empty")
	}

	// Case 1: Just a port number "8080"
	if matched, _ := regexp.MatchString(`^\d+$`, input); matched {
		return fmt.Sprintf("localhost:%s", input), nil
	}

	// Case 2: Port with leading colon ":8080"
	if strings.HasPrefix(input, ":") {
		return "localhost" + input, nil
	}

	// Case 3: Parse as URL to handle http/https defaults
	var u *url.URL
	var err error

	hasScheme := strings.Contains(input, "://")

	if hasScheme {
		u, err = url.Parse(input)
	} else {
		u, err = url.Parse("http://" + input)
	}

	if err != nil {
		return "", fmt.Errorf("invalid target format: %w", err)
	}

	host := u.Hostname()
	port := u.Port()

	// If the original input had https, default port is 443
	if port == "" {
		if u.Scheme == "https" || strings.HasPrefix(input, "https://") {
			port = "443"
		} else {
			port = "80"
		}
	}

	if host == "" {
		host = "localhost"
	}

	return fmt.Sprintf("%s:%s", host, port), nil
}
