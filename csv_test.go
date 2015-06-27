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
		{"NewCSV", "", &CSV{source: resource{}, sink: resource{}, hasHeaderRow: true, headerRow: []string{}, table: [][]string{}}, ""},
	}

	for _, test := range tests {
		value := marshal.Get(NewCSV())
		if value != marshal.Get(test.expected) {
			t.Errorf("%s: expected %s, got %s", test.name, marshal.Get(test.expected), value)
		}
	}
}

func TestReadCSV(t *testing.T) {
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
		file, err := os.Open(test.value)
		if err != nil {
			if err.Error() != test.expectedErr {
				t.Errorf("%s: expected %s, got %s", test.name, test.expectedErr, err.Error())
			}
		} else {

			value, err := ReadCSV(file)
			if err != nil {
				if test.expectedErr != "" {
					t.Errorf("%s: expected an error: %s, but no error was received", test.name, test.expectedErr)
				}
			} else {
				if marshal.Get(value) != marshal.Get(test.expected) {
					t.Errorf("%s: expected %s, got %s", test.name, marshal.Get(test.expected), marshal.Get(value))
				}
			}
		}
		file.Close()
	}
}

func TestReadCSVFile(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		expected    [][]string
		expectedErr string
	}{
		{"invalid filename test", "test_files/tests.csv", [][]string{}, "open test_files/tests.csv: no such file or directory"},
		{"no filename test", "", [][]string{}, "open : no such file or directory"},
		{"valid csv filename test", "test_files/test.csv", [][]string{
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
		res, err := ReadCSVFile(test.filename)
		if err != nil {
			if err.Error() != test.expectedErr {
				t.Errorf("%s expected %s, got %s\n", test.name, test.expectedErr, err.Error())
			}
		} else {
			if test.expectedErr != "" {
				t.Errorf("%s expected an error with %s, no error was returned.\n", test.name, test.expectedErr)
			} else {
				for i, row := range res {
					for j, col := range row {
						if col != test.expected[i][j] {
							t.Error("For ", i, j, "expected: ", test.expected[i][j], " got", col)
						}
					}
				}
			}
		}
	}
}
