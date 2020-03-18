package rules

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vshn/waf-tool/pkg/model"
)

const (
	ruleID   = 10101010
	testdata = "testdata/"
)

var ruleTests = []struct {
	testName string
	alerts   []model.ModsecAlert
	ruleFile string
}{
	{
		"Some path",
		[]model.ModsecAlert{{URI: "/some/path", ID: 10}},
		"rule1.conf",
	},
	{
		"Root path",
		[]model.ModsecAlert{{URI: "/", ID: 9010}},
		"rule2.conf",
	},
	{
		"Combine multiple alerts by path",
		[]model.ModsecAlert{
			{URI: "/path", ID: 9010, RuleTemplate: "# This is the default template"},
			{URI: "/path", ID: 9011, RuleTemplate: "# Some other template"},
			{URI: "/path", ID: 9012},
			{URI: "/path", ID: 942430, Description: `ModSecurity: Warning. Pattern match "((?:[~!@#\\$%\\^&\\*\\(\\)\\-\\+=\\{\\}\\[\\]\\|:;\"'\xc2\xb4\xe2\x80\x99\xe2\x80\x98<>][^~!@#\\$%\\^&\\*\\(\\)\\-\\+=\\{\\}\\[\\]\\|:;\"'\xc2\xb4\xe2\x80\x99\xe2\x80\x98<>]*?){6})" at ARGS:variables.`},
		},
		"rule3.conf",
	},
	{
		"Multiple paths",
		[]model.ModsecAlert{
			{URI: "/path/one", ID: 9010, RuleTemplate: "# This is the default template"},
			{URI: "/path/two", ID: 9011, RuleTemplate: "# Some other template"},
			{URI: "/path/three", ID: 9012},
			{URI: "/some/path", ID: 942430, Description: "ModSecurity: Warning. Pattern match W{4} at ARGS:query.", RuleTemplate: "# ModSec Rule Exclusion: 942430 : Restricted SQL Character Anomaly Detection (args): # of special characters exceeded (12) (severity:  WARNING) PL2"},
		},
		"rule4.conf",
	},
	{
		"With parameter",
		[]model.ModsecAlert{
			{URI: "/some/path", ID: 942430, Description: "ModSecurity: Warning. Pattern match W{4} at ARGS:query."},
			{URI: "/some/path", ID: 942431, Description: `ModSecurity: Warning. Pattern match "((?:[~!@#\\$%\\^&\\*\\(\\)\\-\\+=\\{\\}\\[\\]\\|:;\"'\xc2\xb4\xe2\x80\x99\xe2\x80\x98<>][^~!@#\\$%\\^&\\*\\(\\)\\-\\+=\\{\\}\\[\\]\\|:;\"'\xc2\xb4\xe2\x80\x99\xe2\x80\x98<>]*?){6})" at ARGS:identificationRedirectURL.`},
		},
		"rule5.conf",
	},
	{
		"Non matched parameter",
		[]model.ModsecAlert{
			{URI: "/", ID: 921180, Description: `ModSecurity: Warning. Pattern match "TX:paramcounter_(.*)" at TX:paramcounter_ARGS_NAMES:prospectSingle.contactMethods.contactMethods.value.`},
		},
		"rule6.conf",
	},
}

func TestCreateByIDExclusion(t *testing.T) {
	for _, test := range ruleTests {
		t.Run(test.testName, func(t *testing.T) {
			rule, err := CreateByIDExclusion(test.alerts, ruleID)
			assert.NoError(t, err)
			ruleBytes, err := ioutil.ReadFile(testdata + test.ruleFile)
			assert.NoError(t, err)
			ruleString := string(ruleBytes)
			assert.Equal(t, ruleString, rule)
		})
	}
}

func TestEmptyAlerts(t *testing.T) {
	_, err := CreateByIDExclusion([]model.ModsecAlert{}, 1)
	assert.Error(t, err)
}
