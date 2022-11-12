package jokit

import (
	"strings"
	"testing"
)

func Test_parseQueryParams(t *testing.T) {
	tt := []struct {
		description    string
		params         map[string]interface{}
		expectedString string
	}{
		{
			description: "should return a parsed string",
			params: map[string]interface{}{
				"a": "b",
				"c": "d",
				"e": "f",
			},
			expectedString: "?a=b&c=d&e=f",
		},
		{
			description:    "should return a default params if none is provided",
			params:         map[string]interface{}{},
			expectedString: defaultQueryParams,
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(st *testing.T) {
			got := parseQueryParams(tc.params)
			if !strings.EqualFold(got, tc.expectedString) {
				st.Fatalf("Expected %s got::%s", tc.expectedString, got)
			}
		})
	}
}
