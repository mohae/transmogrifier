package transmogrifier

import (
	"strings"
	"testing"
)

func TestSetSource(t *testing.T) {
	tests := []struct {
		source         string
		expectedFormat FormatType
		expectedErr    string
	}{
		{"", FmtUnsupported, "source string was empty"},
		{"test", FmtUnsupported, "unable to determine format of \"test\""},
		{"test.yaml", FmtUnsupported, "unsupported format for \"test.yaml\": \"yaml\""},
		{"test.CSV", FmtCSV, ""},
		{"test.csv", FmtCSV, ""},
		{"test.MD", FmtMD, ""},
		{"test.md", FmtMD, ""},
	}
	md := NewMDTable()
	for i, test := range tests {
		err := md.SetSource(test.source)
		if err != nil {
			if err.Error() != test.expectedErr {
				t.Errorf("%d: expected %q; got %q", i, test.expectedErr, err.Error())
			}
			continue
		}
		if test.expectedErr != "" {
			t.Errorf("%d: expected error %q: got none", i, test.expectedErr)
			continue
		}
		if test.expectedFormat != md.sourceFormat {
			t.Errorf("%d: expected format to be %q: got %q", i, test.expectedFormat.String(), md.sourceFormat.String())
		}
	}
}

func TestSetDest(t *testing.T) {
	tests := []struct {
		dest         string
		expected     string
		expectedPath string
		expectedName string
	}{
		{"", "", "", ""},
		{"test.md", "test.md", "", "test.md"},
		{"path/to/test.md", "path/to/test.md", "path/to", "test.md"},
	}
	for i, test := range tests {
		md := NewMDTable()
		md.SetDest(test.dest)
		if md.dest.String() != test.expected {
			t.Errorf("%d: expected %q got %q", i, test.expected, md.String())
		}
		if md.dest.Path != test.expectedPath {
			t.Errorf("%d: expected %q got %q", i, test.expectedPath, md.dest.Path)
		}
		if md.dest.Name != test.expectedName {
			t.Errorf("%d: expected %q got %q", i, test.expectedName, md.dest.Name)
		}
	}
}

/*
	tests := []struct{
		useFormat bool
		columnAlignment []string
		expected string
	}{
		{false, []string{}, ""},
		{false, []string{}, ""},
		{false, []string{}, ""},
		{false, []string{}, ""},
		{false, []string{}, ""},
	}
	for i, test := range tests {
		md := NewMDTable()
		md.useFormat = test.useFormat
		md.SetColumnAlignment(test,columnAlingment)

	}
}
*/

func TestBools(t *testing.T) {
	md := NewMDTable()
	md.SetHasColumnNames(true)
	if !md.hasColumnNames {
		t.Error("Expected hasColumnNames to be 'true', was 'false'")
	}
	md.SetHasColumnNames(false)
	if md.hasColumnNames {
		t.Errorf("Expected hasColumNames to be 'false', was 'true'")
	}
	md.SetUseFormat(true)
	if !md.useFormat {
		t.Error("Expected useFormat to be 'true', was 'false'")
	}
	if !md.hasColumnNames {
		t.Error("Expected hasColumnNames to be 'true', was 'false'")
	}
	md.SetUseFormat(false)
	if md.useFormat {
		t.Errorf("Expected useFormat to be 'false', was 'true'")
	}
	if md.hasColumnNames {
		t.Errorf("Expected hasColumNames to be 'false', was 'true'")
	}
}

func TestTransmogrifyStringTable(t *testing.T) {
	rows := [][]string{[]string{
		"name", "id", "desc", "price",
	},
		[]string{
			"towel", "42001", "A hitchhicker's essential: don't leave home without your travel towel.", "19.99",
		},
		[]string{
			"pan-galactic gargle blaster", "42002", "Zaphod's favorite drink. Have one, or some!", "10.00",
		},
		[]string{
			"dinner: restaurant at the end of the universe", "42042", "Good for 1 steak dinner while you watch the universe's favorite table-side show.", "99.99",
		},
	}
	header := []string{
		"name", "id", "desc", "price",
	}
	alignment := []string{
		"", "", "", "",
	}
	emphasis := []string{
		"bold", "italic", "", "strikethrough",
	}
	tests := []struct {
		hasColumnNames bool
		useFormat      bool
		skipRow0       bool
		expected       string
		expectedErr    string
	}{
		{false, false, false,
			`|  
|  
|name|id|desc|price|  
|towel|42001|A hitchhicker's essential: don't leave home without your travel towel.|19.99|  
|pan-galactic gargle blaster|42002|Zaphod's favorite drink. Have one, or some!|10.00|  
|dinner: restaurant at the end of the universe|42042|Good for 1 steak dinner while you watch the universe's favorite table-side show.|99.99|  
`,
			""},
		{true, false, false,
			`|name|id|desc|price|  
|---|---|---|---|  
|towel|42001|A hitchhicker's essential: don't leave home without your travel towel.|19.99|  
|pan-galactic gargle blaster|42002|Zaphod's favorite drink. Have one, or some!|10.00|  
|dinner: restaurant at the end of the universe|42042|Good for 1 steak dinner while you watch the universe's favorite table-side show.|99.99|  
`,
			""},

		{false, true, false,
			`|name|id|desc|price|  
|---|---|---|---|  
|__towel__|_42001_|A hitchhicker's essential: don't leave home without your travel towel.|~~19.99~~|  
|__pan-galactic gargle blaster__|_42002_|Zaphod's favorite drink. Have one, or some!|~~10.00~~|  
|__dinner: restaurant at the end of the universe__|_42042_|Good for 1 steak dinner while you watch the universe's favorite table-side show.|~~99.99~~|  
`,
			""},
		{true, true, false,
			`|name|id|desc|price|  
|---|---|---|---|  
|__towel__|_42001_|A hitchhicker's essential: don't leave home without your travel towel.|~~19.99~~|  
|__pan-galactic gargle blaster__|_42002_|Zaphod's favorite drink. Have one, or some!|~~10.00~~|  
|__dinner: restaurant at the end of the universe__|_42042_|Good for 1 steak dinner while you watch the universe's favorite table-side show.|~~99.99~~|  
`,
			""},

		{false, false, true,
			`|  
|  
|towel|42001|A hitchhicker's essential: don't leave home without your travel towel.|19.99|  
|pan-galactic gargle blaster|42002|Zaphod's favorite drink. Have one, or some!|10.00|  
|dinner: restaurant at the end of the universe|42042|Good for 1 steak dinner while you watch the universe's favorite table-side show.|99.99|  
`,
			""},
		{true, false, true,
			`|towel|42001|A hitchhicker's essential: don't leave home without your travel towel.|19.99|  
|---|---|---|---|  
|pan-galactic gargle blaster|42002|Zaphod's favorite drink. Have one, or some!|10.00|  
|dinner: restaurant at the end of the universe|42042|Good for 1 steak dinner while you watch the universe's favorite table-side show.|99.99|  
`,
			""},

		{false, true, true,
			`|towel|42001|A hitchhicker's essential: don't leave home without your travel towel.|19.99|  
|---|---|---|---|  
|__pan-galactic gargle blaster__|_42002_|Zaphod's favorite drink. Have one, or some!|~~10.00~~|  
|__dinner: restaurant at the end of the universe__|_42042_|Good for 1 steak dinner while you watch the universe's favorite table-side show.|~~99.99~~|  
`,
			""},
		{true, true, true,
			`|towel|42001|A hitchhicker's essential: don't leave home without your travel towel.|19.99|  
|---|---|---|---|  
|__pan-galactic gargle blaster__|_42002_|Zaphod's favorite drink. Have one, or some!|~~10.00~~|  
|__dinner: restaurant at the end of the universe__|_42042_|Good for 1 steak dinner while you watch the universe's favorite table-side show.|~~99.99~~|  
`,
			""},
	}
	for i, test := range tests {
		md := NewMDTable()
		if test.hasColumnNames {
			md.SetHasColumnNames(test.hasColumnNames)
			md.SetColumnNames(header)

		}
		if test.useFormat {
			md.SetUseFormat(test.useFormat)
			md.SetColumnNames(header)
			md.SetColumnAlignment(alignment)
			md.SetColumnEmphasis(emphasis)
		}
		var err error
		if test.skipRow0 {
			err = md.TransmogrifyStringTable(rows[1:])
		} else {
			err = md.TransmogrifyStringTable(rows)
		}
		if err != nil {
			t.Errorf("%d: expected no error, got %q", i, err)
			continue
		} else {
			if md.String() != test.expected {
				t.Errorf("%d: expected\n%s\ngot\n%s", i, strings.Replace(strings.Replace(test.expected, " ", "+", -1), "\n", "-", -1), strings.Replace(strings.Replace(md.String(), " ", "+", -1), "\n", "-", -1))
			}
		}
	}
}

func TestAppendHeaderSeparatorRow(t *testing.T) {
	tests := []struct {
		cols      []string
		alignment []string
		expected  string
	}{
		{cols: []string{"", "", ""}, alignment: []string{"", "", ""}, expected: "|---|---|---|  \n"},
		{[]string{"col1", "col2", "col3"}, []string{"", "", ""}, "|---|---|---|  \n"},
		{[]string{"col1", "col2", "col3"}, []string{"l", "r", "c"}, "|:---|---:|:---:|  \n"},
		{[]string{"col1", "col2", "col3"}, []string{"c", "l", "r"}, "|:---:|:---|---:|  \n"},
		{[]string{"col1", "col2", "col3"}, []string{"r", "c", "l"}, "|---:|:---:|:---|  \n"},
		{[]string{"col1", "col2", "col3"}, []string{"l", "", "l"}, "|:---|---|:---|  \n"},
		{[]string{"col1", "col2", "col3"}, []string{"c", "", "c"}, "|:---:|---|:---:|  \n"},
		{[]string{"col1", "col2", "col3"}, []string{"r", "", "r"}, "|---:|---|---:|  \n"},
		{[]string{"col1", "col2", "col3"}, []string{"l", "l", ""}, "|:---|:---|---|  \n"},
		{[]string{"col1", "col2", "col3"}, []string{"c", "c", ""}, "|:---:|:---:|---|  \n"},
		{[]string{"col1", "col2", "col3"}, []string{"r", "r", ""}, "|---:|---:|---|  \n"},
		{[]string{"col1", "col2", "col3"}, []string{"", "l", "l"}, "|---|:---|:---|  \n"},
		{[]string{"col1", "col2", "col3"}, []string{"", "c", "c"}, "|---|:---:|:---:|  \n"},
		{[]string{"col1", "col2", "col3"}, []string{"", "r", "r"}, "|---|---:|---:|  \n"},
	}

	for i, test := range tests {
		md := NewMDTable()
		md.useFormat = true
		md.SetColumnNames(test.cols)
		md.SetColumnAlignment(test.alignment)
		md.appendHeaderSeparatorRow()
		if md.String() != test.expected {
			t.Errorf("%d: expected %q got %q", i, test.expected, md.String())
		}
	}
}

func TestMDFilenameFrom(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"", ""},
		{"test", "test.md"},
		{"test.csv", "test.md"},
		{"path/to/test.csv", "path/to/test.md"},
		{"test.file.csv", "test.file.md"},
	}
	for i, test := range tests {
		dest := mdFilenameFrom(test.name)
		if dest != test.expected {
			t.Errorf("%d: expected %q, got %q", i, test.expected, dest)
		}
	}
}

func TestFormatFromFile(t *testing.T) {
	tests := []struct {
		name                    string
		expectedColumnNames     []string
		expectedColumnAlignment []string
		expectedColumnEmphasis  []string
		expectedSetErr          string
		expectedFormatErr       string
	}{
		{"", []string{}, []string{}, []string{}, "unable to set format source: received empty string", ""},
		{"test_files/bad-test.fmt", []string{}, []string{}, []string{}, "", "insufficient format rows: expected at least 3, got 2"},
		{"test_files/test.fmt", []string{"Item", "Id", "Description", "Price"}, []string{"left", "", "centered", "right"}, []string{"bold", "italic", "strikethrough", ""}, "", ""},
	}
	for i, test := range tests {
		md := NewMDTable()
		err := md.SetFormatSource(test.name)
		if err != nil {
			if err.Error() != test.expectedSetErr {
				t.Errorf("%d: expected %q got %q", i, test.expectedSetErr, err.Error())
			}
			continue
		}
		if test.expectedSetErr != "" {
			t.Errorf("%d: expected %q, got no error", i, test.expectedSetErr)
			continue
		}
		err = md.formatFromFile()
		if err != nil {
			if err.Error() != test.expectedFormatErr {
				t.Errorf("%d: expected %q got %q", i, test.expectedFormatErr, err.Error())
			}
			continue
		}
		if test.expectedFormatErr != "" {
			t.Errorf("%d: expected %q, got no error", i, test.expectedFormatErr)
			continue
		}
		for k, v := range md.columnNames {
			if test.expectedColumnNames[k] != v {
				t.Errorf("%d-%d: expected %q, got %q", i, k, test.expectedColumnNames[k], v)
			}
		}
		for k, v := range md.columnAlignment {
			if test.expectedColumnAlignment[k] != v {
				t.Errorf("%d-%d: expected %q, got %q", i, k, test.expectedColumnAlignment[k], v)
			}
		}
		for k, v := range md.columnEmphasis {
			if test.expectedColumnEmphasis[k] != v {
				t.Errorf("%d-%d: expected %q, got %q", i, k, test.expectedColumnEmphasis[k], v)
			}
		}
	}
}
