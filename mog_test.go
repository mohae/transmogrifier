package transmogrifier

import "testing"

func TestNewResource(t *testing.T) {
	tests := []struct {
		path     string
		expected resource
	}{
		{"", resource{}},
		{"test.txt", resource{Name: "test.txt"}},
		{"path/test.txt", resource{Name: "test.txt", Path: "path"}},
		{"another/path/test", resource{Name: "test", Path: "another/path"}},
	}
	for _, test := range tests {
		r := NewResource(test.path)
		if r != test.expected {
			t.Errorf("Expected %+v, got %+v", test.expected, r)
		}
	}
}
