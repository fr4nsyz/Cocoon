package isolation

import (
	"regexp"
	"strings"
)

var SensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(api[_-]?key|apikey)['"]?\s*[:=]\s*['"]?[\w-]{20,}`),
	regexp.MustCompile(`(?i)secret['"]?\s*[:=]\s*['"]?[\w-]{20,}`),
	regexp.MustCompile(`(?i)password['"]?\s*[:=]\s*['"]?[\w-]{8,}`),
	regexp.MustCompile(`(?i)token['"]?\s*[:=]\s*['"]?[\w-]{20,}`),
	regexp.MustCompile(`(?i)private[_-]?key['"]?\s*[:=]\s*['"]?[\w-]{20,}`),
	regexp.MustCompile(`sk-[a-zA-Z0-9]{20,}`),
}

func ScanEnvForSecrets(env []string) map[string]string {
	secrets := make(map[string]string)

	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		for _, pattern := range SensitivePatterns {
			if pattern.MatchString(value) {
				secrets[key] = "REDACTED"
				break
			}
		}
	}

	return secrets
}

func IsPathAllowed(path string, projectDir string) bool {
	normalizedPath := strings.TrimSpace(path)
	if normalizedPath == "" {
		return false
	}

	if strings.HasPrefix(normalizedPath, "/sandbox") ||
		strings.HasPrefix(normalizedPath, "./") ||
		strings.HasPrefix(normalizedPath, "../") {
		return true
	}

	return false
}
