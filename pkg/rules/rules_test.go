package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vshn/waf-tool/pkg/model"
)

const (
	ruleID = 10101010
)

var ruleTests = []struct {
	testName string
	alerts   []model.ModsecAlert
	rule     string
}{
	{
		"Some path",
		[]model.ModsecAlert{{URI: "/some/path", ID: 10}},
		`
SecRule REQUEST_URI "@strmatch /some/path" \
  "phase:2,nolog,id:10101010,\
  ctl:ruleRemoveById=10"
`,
	},
	{
		"Root path",
		[]model.ModsecAlert{{URI: "/", ID: 9010}},
		`
SecRule REQUEST_URI "@strmatch /" \
  "phase:2,nolog,id:10101010,\
  ctl:ruleRemoveById=9010"
`,
	},
	{
		"Combine multiple alerts by path",
		[]model.ModsecAlert{
			{URI: "/path", ID: 9010, RuleTemplate: "# This is the default template"},
			{URI: "/path", ID: 9011, RuleTemplate: "# Some other template"},
			{URI: "/path", ID: 9012},
			{URI: "/path", ID: 942430, Description: `ModSecurity: Warning. Pattern match "((?:[~!@#\\$%\\^&\\*\\(\\)\\-\\+=\\{\\}\\[\\]\\|:;\"'\xc2\xb4\xe2\x80\x99\xe2\x80\x98<>][^~!@#\\$%\\^&\\*\\(\\)\\-\\+=\\{\\}\\[\\]\\|:;\"'\xc2\xb4\xe2\x80\x99\xe2\x80\x98<>]*?){6})" at ARGS:variables.`},
		},
		`
# Path /path
# This is the default template
# Some other template
SecRule REQUEST_URI "@strmatch /path" \
  "phase:2,nolog,id:10101010,\
  ctl:ruleRemoveById=9010,\
  ctl:ruleRemoveById=9011,\
  ctl:ruleRemoveById=9012,\
  ctl:ruleRemoveTargetById=942430;ARGS:variables"
`,
	},
	{
		"Multiple paths",
		[]model.ModsecAlert{
			{URI: "/path/one", ID: 9010, RuleTemplate: "# This is the default template"},
			{URI: "/path/two", ID: 9011, RuleTemplate: "# Some other template"},
			{URI: "/path/three", ID: 9012},
			{URI: "/some/path", ID: 942430, Description: "ModSecurity: Warning. Pattern match W{4} at ARGS:query.", RuleTemplate: "# ModSec Rule Exclusion: 942430 : Restricted SQL Character Anomaly Detection (args): # of special characters exceeded (12) (severity:  WARNING) PL2"},
		},
		`
# This is the default template
SecRule REQUEST_URI "@strmatch /path/one" \
  "phase:2,nolog,id:10101010,\
  ctl:ruleRemoveById=9010"

SecRule REQUEST_URI "@strmatch /path/three" \
  "phase:2,nolog,id:10101011,\
  ctl:ruleRemoveById=9012"

# Some other template
SecRule REQUEST_URI "@strmatch /path/two" \
  "phase:2,nolog,id:10101012,\
  ctl:ruleRemoveById=9011"

# ModSec Rule Exclusion: 942430 : Restricted SQL Character Anomaly Detection (args): # of special characters exceeded (12) (severity:  WARNING) PL2
SecRule REQUEST_URI "@strmatch /some/path" \
  "phase:2,nolog,id:10101013,\
  ctl:ruleRemoveTargetById=942430;ARGS:query"
`,
	},
	{
		"With parameter",
		[]model.ModsecAlert{
			{URI: "/some/path", ID: 942430, Description: "ModSecurity: Warning. Pattern match W{4} at ARGS:query."},
			{URI: "/some/path", ID: 942431, Description: `ModSecurity: Warning. Pattern match "((?:[~!@#\\$%\\^&\\*\\(\\)\\-\\+=\\{\\}\\[\\]\\|:;\"'\xc2\xb4\xe2\x80\x99\xe2\x80\x98<>][^~!@#\\$%\\^&\\*\\(\\)\\-\\+=\\{\\}\\[\\]\\|:;\"'\xc2\xb4\xe2\x80\x99\xe2\x80\x98<>]*?){6})" at ARGS:identificationRedirectURL.`},
		},
		`
# Path /some/path
SecRule REQUEST_URI "@strmatch /some/path" \
  "phase:2,nolog,id:10101010,\
  ctl:ruleRemoveTargetById=942430;ARGS:query,\
  ctl:ruleRemoveTargetById=942431;ARGS:identificationRedirectURL"
`,
	},
	{
		"Non matched parameter",
		[]model.ModsecAlert{
			{URI: "/", ID: 921180, Description: `ModSecurity: Warning. Pattern match "TX:paramcounter_(.*)" at TX:paramcounter_ARGS_NAMES:prospectSingle.contactMethods.contactMethods.value.`},
		},
		`
SecRule REQUEST_URI "@strmatch /" \
  "phase:2,nolog,id:10101010,\
  ctl:ruleRemoveById=921180"
`,
	},
}

func TestCreateByIDExclusion(t *testing.T) {
	for _, test := range ruleTests {
		t.Run(test.testName, func(t *testing.T) {
			rule, err := CreateByIDExclusion(test.alerts, ruleID)
			assert.NoError(t, err)
			assert.Equal(t, test.rule, rule)
		})
	}
}

func TestEmptyAlerts(t *testing.T) {
	_, err := CreateByIDExclusion([]model.ModsecAlert{}, 1)
	assert.Error(t, err)
}
