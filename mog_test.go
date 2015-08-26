package transmogrifier

import "testing"

func TestNewResource(t *testing.T) {
	tests := []struct {
		path     string
		format   FormatType
		typ      ResourceType
		expected resource
	}{
		{"", FmtUnsupported, UnsupportedResource, resource{Format: FmtUnsupported, Type: UnsupportedResource}},
		{"test.txt", FmtUnsupported, File, resource{Name: "test.txt", Format: FmtUnsupported, Type: File}},
		{"path/test.csv", FmtCSV, File, resource{Name: "test.csv", Path: "path", Format: FmtCSV, Type: File}},
		{"another/path/test", FmtMDTable, File, resource{Name: "test", Path: "another/path", Format: FmtMDTable, Type: File}},
	}
	for i, test := range tests {
		r := NewResource(test.path, test.format, test.typ)
		if r != test.expected {
			t.Errorf("%d: expected %#v, got %#v", i, test.expected, r)
		}
	}
}
