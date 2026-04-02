package isolation

import (
	"strings"
)

var AllowedHosts = []string{
	"registry.npmjs.org",
	"npmjs.org",
	"yarnpkg.com",
	"pypi.org",
	"pip.confederation.tech",
	"files.pythonhosted.org",
	"github.com",
	"raw.githubusercontent.com",
}

func IsNetworkAllowed(host string) bool {
	host = strings.ToLower(strings.TrimSpace(host))
	if host == "" {
		return false
	}

	for _, allowed := range AllowedHosts {
		if strings.Contains(host, allowed) {
			return true
		}
	}

	return false
}

func GetAllowedHosts() []string {
	result := make([]string, len(AllowedHosts))
	copy(result, AllowedHosts)
	return result
}

func AddAllowedHost(host string) {
	host = strings.TrimSpace(host)
	if host != "" {
		AllowedHosts = append(AllowedHosts, host)
	}
}
