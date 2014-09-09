package mogger

import (
	"os"
	"testing"

	json "github.com/mohae/customjson"
)

var marshal = json.NewMarshalString()
var tableData, tableDataNoHeader [][]string

func init() {
	tableData =  [][]string{
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

	tableDataNoHeader =  [][]string{
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
	tests :=  []struct{
		name string
		value string
		expected *CSV
		expectedErr string
	}{
		{"NewCSV", "", &CSV{hasHeaderRow: true, destinationType: "bytes", table: [][]string{}}, ""},
	}

	for _, test := range tests {
		value := marshal.Get(NewCSV())
		if value != marshal.Get(test.expected) {
			t.Errorf("%s: expected %s, got %s", test.name, marshal.Get(test.expected), value)
		}
	}
}


func TestFileToTable(t *testing.T) {
	tests := []struct{
		name string
		filename string
		expected string
		expectedErr string
	}{
		{"test.csv", "tests/test.csv", "|Item|Description|Price||--|--|--||string| a string of indeterminate length| $9.99||towel| an intergalactic traveller's essential| $42.00|", ""},
		{"no filename", "", "", "open : no such file or directory"},
	}

	for _, test := range tests{
		c := NewCSV()
		err  := c.FileToTable(test.filename)	
		if err != nil {
			if err.Error() != test.expectedErr {
				t.Errorf("%s: expected %s, got %s", test.name, test.expectedErr, err.Error())
			}
		} else {
			if test.expectedErr != "" {
				t.Errorf("%s: expected %s, but no error was received.", test.name, test.expectedErr)
			} else {
				md := c.MD()
				if string(md) != test.expected {
					t.Errorf("%s: expected %s, got %s", test.name, test.expected, string(md))
				}
			}
		}
	}
}

func TestReadCSV(t *testing.T) {
	tests :=  []struct{
		name string
		value string
		expected [][]string
		expectedErr string
	}{
		{"Test Read CSV", "tests/test.csv", [][]string{
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
	tests := []struct{
		name string
		filename string
		expected [][]string
		expectedErr string
	}{
		{"invalid filename test", "tests/tests.csv", [][]string{}, "open tests/tests.csv: no such file or directory"},
                {"no filename test", "", [][]string{}, "open : no such file or directory"},
		{"valid csv filename test", "tests/test.csv", [][]string{
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


func TestToMD(t *testing.T) {
	tests :=  []struct{
		name string
		hasHeader bool
		table [][]string
		expected string
	}{
		{"TOMD w header", true, tableData, `|Item|Description|Price||--|--|--||string| a string of indeterminate length| $9.99||towel| an intergalactic traveller's essential| $42.00|`},
		{"TOMD w/o header", false, tableDataNoHeader, `|--|--|--||book|Has the words "don't panic" in large, friendly letters on the cover|Price||string| a string of indeterminate length| $9.99||towel| an intergalactic traveller's essential| $42.00|`},
	}

	for _, test := range tests {
		c := NewCSV()
		c.hasHeaderRow = test.hasHeader
		c.table = test.table
		c.toMD()
		if  string(c.md) != test.expected {
			t.Errorf("%s: expected %s, got %s", test.name, test.expected, string(c.md))
		}
	}
}

func TestRowToMD(t *testing.T) {
	tests :=  []struct{
		name string
		value []string
		expected string
	}{
		{"test 2 cols", []string{"col1", "col2"}, "|col1|col2|"},
		{"test 3 cols", []string{"col1", "col2", "col3"}, "|col1|col2|col3|"},
		{"test 4 cols", []string{"col1", "col2", "col3", "col4"}, "|col1|col2|col3|col4|"},
		{"test 5 cols", []string{"col1", "col2", "col3", "col4", "col5"}, "|col1|col2|col3|col4|col5|"},
	}


	for _, test := range tests {
		c := NewCSV()
		c.rowToMD(test.value)
		if string(c.md) != test.expected {
			t.Errorf("%s: expected %s, got %s", test.name, test.expected, string(c.md))
		}
	}
}


func TestAddHeader(t *testing.T) {
	tests :=  []struct{
		name string
		headerRow bool
		data [][]string
		expected string
	}{
		{"w header row", false, tableData, "|--|--|--|"},
		{"without header row", true, tableData, "|Item|Description|Price||--|--|--|"},
	}


	for _, test := range tests {
		c := NewCSV()
		c.hasHeaderRow =  test.headerRow
		c.table = test.data		
		c.addHeader()
		if string(c.md) != test.expected {
			t.Errorf("%s: expected %s, got %s", test.name, test.expected, string(c.md))
		}
	}
}

func TestAppendHeaderSeparatorRow(t *testing.T) {
	tests :=  []struct{
		name string
		cols int
		expected string
	}{
		{"3 cols", 3, "|--|--|--|"},
		{"4 cols", 4, "|--|--|--|--|"},
	}


	for _, test := range tests {
		c := NewCSV()		
		c.appendHeaderSeparatorRow(test.cols)
		if string(c.md) != test.expected {
			t.Errorf("%s: expected %s, got %s", test.name, test.expected, string(c.md))
		}
	}
}

func TestAppendColumnSeparator(t *testing.T) {
	tests :=  []struct{
		name string
		expected string
	}{
		{"append column separator", "|"},
	}

	
	for _, test := range tests {
		c := NewCSV()
		c.appendColumnSeparator()
		if string(c.md) != test.expected {
			t.Errorf("%s: expected %s, got %s", test.name, test.expected, string(c.md))
		}
	}
}
