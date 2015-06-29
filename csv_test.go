package transmogrifier

import (
	"os"
	"testing"

	json "github.com/mohae/customjson"
)

var marshal = json.NewMarshalString()
var tableData, tableDataNoHeader [][]string

func init() {
	tableData = [][]string{
		[]string{
			"Item",
			"Description",
			"Price",
		},
		[]string{
			"string",
			" a string of indeterminate length",
			" $9.99",
		},
		[]string{
			"towel",
			" an intergalactic traveller's essential",
			" $42.00",
		}}

	tableDataNoHeader = [][]string{
		[]string{
			"book",
			"Has the words \"don't panic\" in large, friendly letters on the cover",
			"Price",
		},
		[]string{
			"string",
			" a string of indeterminate length",
			" $9.99",
		},
		[]string{
			"towel",
			" an intergalactic traveller's essential",
			" $42.00",
		}}

}

func TestNewCSV(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expected    *CSV
		expectedErr string
	}{
		{"NewCSV", "", &CSV{source: resource{}, sink: resource{}, hasHeader: true, headerRow: []string{}, rows: [][]string{}}, ""},
	}

	for _, test := range tests {
		value := marshal.Get(NewCSV())
		if value != marshal.Get(test.expected) {
			t.Errorf("%s: expected %s, got %s", test.name, marshal.Get(test.expected), value)
		}
	}
}

func TestNewCSVSource(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expected    *CSV
		expectedErr string
	}{
		{"NewCSV", "hello.csv", &CSV{source: resource{Name: "hello.csv"}, sink: resource{}, hasHeader: true, headerRow: []string{}, rows: [][]string{}}, ""},
		{"NewCSV", "test/hello.csv", &CSV{source: resource{Name: "hello.csv", Path: "test"}, sink: resource{}, hasHeader: true, headerRow: []string{}, rows: [][]string{}}, ""},
	}

	for _, test := range tests {
		value := marshal.Get(NewCSVSource(test.value))
		if value != marshal.Get(test.expected) {
			t.Errorf("%s: expected %s, got %s", test.name, marshal.Get(test.expected), value)
		}
	}
}
func TestRead(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expected    [][]string
		expectedErr string
	}{
		{"Test Read CSV", "test_files/test.csv", [][]string{
			[]string{
				"Item",
				"Description",
				"Price",
			},
			[]string{
				"string",
				" a string of indeterminate length",
				" $9.99",
			},
			[]string{
				"towel",
				" an intergalactic traveller's essential",
				" $42.00",
			},
		}, ""},
	}

	for _, test := range tests {
		c := NewCSV()
		c.hasHeader = false
		file, err := os.Open(test.value)
		if err != nil {
			if err.Error() != test.expectedErr {
				t.Errorf("%s: expected %s, got %s", test.name, test.expectedErr, err.Error())
			}
			goto close
		}
		err = c.Read(file)
		if err != nil {
			if test.expectedErr != "" {
				t.Errorf("%s: expected an error: %s, but no error was received", test.name, test.expectedErr)
			}
			goto close
		}
		if marshal.Get(c.rows) != marshal.Get(test.expected) {
			t.Errorf("%s: expected %s, got %s", test.name, marshal.Get(test.expected), marshal.Get(c.rows))
		}
	close:
		file.Close()
	}
}

func TestReadFile(t *testing.T) {
	tests := []struct {
		name           string
		hasHeader      bool
		filename       string
		expectedHeader []string
		expectedRows   [][]string
		expectedErr    string
	}{
		{"invalid filename test", false, "test_files/tests.csv", []string{}, [][]string{}, "open test_files/tests.csv: no such file or directory"},
		{"no filename test", false, "", []string{}, [][]string{}, "no source was specified"},
		{"valid csv filename test", false, "test_files/test.csv", []string{}, [][]string{
			[]string{
				"Item",
				"Description",
				"Price",
			},
			[]string{
				"string",
				" a string of indeterminate length",
				" $9.99",
			},
			[]string{
				"towel",
				" an intergalactic traveller's essential",
				" $42.00",
			},
		}, ""},
		{"valid csv filename test", true, "test_files/test.csv", []string{
			"Item",
			"Description",
			"Price",
		}, [][]string{
			[]string{
				"string",
				" a string of indeterminate length",
				" $9.99",
			},
			[]string{
				"towel",
				" an intergalactic traveller's essential",
				" $42.00",
			},
		}, ""},
	}

	for _, test := range tests {
		c := NewCSV()
		c.hasHeader = test.hasHeader
		err := c.ReadFile(test.filename)
		if err != nil {
			if err.Error() != test.expectedErr {
				t.Errorf("%s expected %s, got %s\n", test.name, test.expectedErr, err.Error())
			}
			continue
		}
		if test.expectedErr != "" {
			t.Errorf("%s expected an error with %s, no error was returned.\n", test.name, test.expectedErr)
			continue
		}
		for i, col := range c.headerRow {
			if col != test.expectedHeader[i] {
				t.Error("header col %d: expected %s, got %s", test.expectedHeader[i], col)
			}
		}
		for i, row := range c.rows {
			for j, col := range row {
				if col != test.expectedRows[i][j] {
					t.Error("For ", i, j, "expected: ", test.expectedRows[i][j], " got", col)
				}
			}
		}
	}
}
