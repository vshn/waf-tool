package rules

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/vshn/waf-tool/pkg/model"
)

const (
	// NewLine character
	NewLine = '\n'
)

// CreateByIDExclusion creates a rule exclusion combined by path.
// If no parameters are found, the whole rule will be removed.
func CreateByIDExclusion(alerts []model.ModsecAlert, ruleBaseID int) (string, error) {
	if len(alerts) == 0 {
		return "", errors.New("no ModSecurity alerts")
	}
	pathMap := map[string][]model.ModsecAlert{}
	// Combine alerts by path
	for _, alert := range alerts {
		if alertsForURL, ok := pathMap[alert.URI]; ok {
			pathMap[alert.URI] = append(alertsForURL, alert)
		} else {
			pathMap[alert.URI] = []model.ModsecAlert{alert}
		}
	}
	var builder strings.Builder

	paths := make([]string, 0, len(pathMap))
	for path := range pathMap {
		paths = append(paths, path)
	}
	// Sort by path to have a stable order of generated rules
	sort.Strings(paths)
	for _, path := range paths {
		alerts := pathMap[path]
		if len(alerts) > 1 {
			fmt.Fprintf(&builder, "\n# Path %s", path)
		}
		for _, alert := range alerts {
			if len(alert.RuleTemplate) > 0 {
				builder.WriteRune(NewLine)
				builder.WriteString(alert.RuleTemplate)
			}
		}

		fmt.Fprintf(&builder, `
SecRule REQUEST_URI "@strmatch %s" \
  "phase:2,nolog,id:%d,\
`, path, ruleBaseID)

		for i, alert := range alerts {
			rule := fmt.Sprintf("  ctl:ruleRemoveById=%d", alert.ID)
			if param, ok := alert.ExtractParameter(); ok {
				rule = fmt.Sprintf("  ctl:ruleRemoveTargetById=%d;%s", alert.ID, param)
			}
			builder.WriteString(rule)
			if i < len(alerts)-1 {
				builder.WriteString(`,\`)
				builder.WriteRune(NewLine)
			} else {
				builder.WriteRune('"')
				builder.WriteRune(NewLine)
			}
		}
		ruleBaseID++
	}
	return builder.String(), nil
}
