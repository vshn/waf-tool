package model

import (
	"regexp"
)

var (
	alertParameter = regexp.MustCompile(`^ModSecurity: Warning. Pattern match.+ (at|against) "?([^T][^X].+)"?\.$`)
)

// ExtractParameter extracts the paramter of an alert, if any
func (alert ModsecAlert) ExtractParameter() (string, bool) {
	if params := alertParameter.FindStringSubmatch(alert.Description); params != nil {
		if len(params) == 3 {
			return params[2], true
		}
	}
	return "", false
}
