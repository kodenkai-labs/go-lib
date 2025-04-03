package url_test

import (
	"testing"

	"github.com/kodenkai-labs/go-lib/url"
)

func Test_ParseQueryString(t *testing.T) {
	tests := []struct {
		name        string
		queryString string
		expected    map[string]string
	}{
		{
			name:        "Empty query string",
			queryString: "",
			expected:    map[string]string{},
		},
		{
			name:        "One pair",
			queryString: "key=value",
			expected:    map[string]string{"key": "value"},
		},
		{
			name:        "Two pairs",
			queryString: "key1=value1&key2=value2",
			expected:    map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			name:        "Two pairs with empty value",
			queryString: "key1=&key2=value2",
			expected:    map[string]string{"key1": "", "key2": "value2"},
		},
		{
			name:        "Two pairs with empty key",
			queryString: "=value1&key2=value2",
			expected:    map[string]string{"": "value1", "key2": "value2"},
		},
		{
			name:        "Two pairs with empty key and value",
			queryString: "=&key2=value2",
			expected:    map[string]string{"": "", "key2": "value2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := url.ParseQueryString(tt.queryString)
			if len(result) != len(tt.expected) {
				t.Fatalf("Expected %d pairs, got %d", len(tt.expected), len(result))
			}

			for key, value := range tt.expected {
				if result[key] != value {
					t.Fatalf("Expected value %s for key %s, got %s", value, key, result[key])
				}
			}
		})
	}
}
