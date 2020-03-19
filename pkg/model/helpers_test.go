package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var ruleTests = []struct {
	testName string
	alert    ModsecAlert
	param    string
}{
	{
		"Empty description",
		ModsecAlert{Description: ""},
		"",
	},
	{
		"No param",
		ModsecAlert{Description: "ModSecurity: Warning. Pattern match"},
		"",
	},
	{
		"ARG param",
		ModsecAlert{Description: `ModSecurity: Warning. Pattern match "\\W{4}" at ARGS:identificationRedirectURL.`},
		"ARGS:identificationRedirectURL",
	},
	{
		"COOKIE param",
		ModsecAlert{Description: `ModSecurity: Warning. Pattern match "(?:[\"'][\\s\\d]*?[^\\w\\s]+\\W*?\\d\\W*?.*?[\"'\\d])" at REQUEST_COOKIES:_pk_ref.13.3b4d.`},
		"REQUEST_COOKIES:_pk_ref.13.3b4d",
	},
	{
		"ARG param 2",
		ModsecAlert{Description: `ModSecurity: Warning. Pattern match "((?:[~!@#\\$%\\^&\\*\\(\\)\\-\\+=\\{\\}\\[\\]\\|:;\"'\xc2\xb4\xe2\x80\x99\xe2\x80\x98<>][^~!@#\\$%\\^&\\*\\(\\)\\-\\+=\\{\\}\\[\\]\\|:;\"'\xc2\xb4\xe2\x80\x99\xe2\x80\x98<>]*?){6})" at ARGS:query.`},
		"ARGS:query",
	},
	{
		"TX param",
		ModsecAlert{Description: `ModSecurity: Warning. Pattern match "TX:paramcounter_(.*)" at TX:paramcounter_ARGS_NAMES:prospectSingle.contactMethods.contactMethods.value.`},
		"",
	},
	{
		"TX param 2",
		ModsecAlert{Description: `ModSecurity: Warning. Pattern match "TX:paramcounter_(.*)" at TX:paramcounter_ARGS_NAMES:orders.products.products.id.`},
		"",
	},
	{
		"No TX param",
		ModsecAlert{Description: "ModSecurity: Warning. Pattern match W{4} at ARGS:TXvariable."},
		"ARGS:TXvariable",
	},
}

func TestExtractParameter(t *testing.T) {
	for _, test := range ruleTests {
		t.Run(test.testName, func(t *testing.T) {
			param, ok := test.alert.ExtractParameter()
			if len(test.param) > 0 {
				assert.True(t, ok)
				assert.Equal(t, test.param, param)
			} else {
				assert.False(t, ok)
			}
		})
	}
}
