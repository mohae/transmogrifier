package transmogrifier

import (
	_ "fmt"
	_ "io"
	"log"
	"os"
	"path/filepath"
	_ "strconv"
	"strings"
)

// MDTable format representations.
var (
	// Pipe is the MD column separator
	mdPipe []byte = []byte("|")

	// LeftJustify is the MD for left justification of columns.
	mdLeftJustify []byte = []byte(":---")

	// RightJustify is the Md for right justification of columns,
	mdRightJustify []byte = []byte("---:")
	mdCentered     []byte = []byte(":---:")
	mdDontJustify  []byte = []byte("---")
)

// MDTable is a struct for representing and working with markdown tables
type MDTable struct {
	// data source, if applicable.
	source resource
	// sink, if applicable
	sink resource
	// FormatSource: the location and name of the source file to use. It
	// can either be explicitely set, or TOMD will look for it as
	// `source.fmt`, for `source.csv`.
	formatSource string
	// whether the formatSource was autoset or not.
	formatSourceAutoset bool
	// useFormat: whether there's a format to use with the CSV or not. For
	// files, this is usually a file, with the same name and path as the
	// source, using the 'fmt' extension. This can also be set explicitely.
	// 'useFormat' == false implies 'hasHeaderRow' == true.
	useFormat bool
	// formatType:	the type of format to use. By default, this is in sync
	//		with the source type, but it can be set independently.
	// Supported:
	//	file	The format information is in a format file. By default,
	//		this is the source filename with the `.fmt` file
	//		extension, instead of the original extension. This can
	//		be set independently too.
	//	default Any setting other than another supported type will be
	//		interpreted as using the default, which is to manually
	//		set the different format information you wish to use
	//		in the marshal using their Setters.
	formatType string
	// ColumnAlignment contains the alignment information, if any, for each
	// column.  This is supplied by the format.
	columnAlignment []string
	// ColumnEmphasis contains the emphasis information, if any. for each column.
	// This is supplied by the format.
	columnEmphasis []string
	md             []byte
}

func NewMDTable() *MDTable {
	return &MDTable{columnAlignment: []string{}, columnEmphasis: []string{}, Table: [][]string{}, md: []byte{}}
}

// FromCSV creates a md table from CSV.
func (m *MDTable) FromCSV(c *CSV) error {
	// Process the header first
	m.addHeader()
	// for each row of table data, process it.
	for _, row := range c.table {
		m.rowToMD(row)
	}
	return
}

// rowTomd takes a table row and returns the md version of it consistent
// with its configuration.
func (m *MDTable) rowToMD(cols []string) {
	m.appendColumnSeparator()
	for _, col := range cols {
		// TODO this is where column data decoration would occur
		// with templates
		bcol := []byte(col)
		m.md = append(m.md, bcol...)
		m.appendColumnSeparator()
	}
	// add a new line at the end of a row
	m.md = append(m.md, []byte("  \n")...)
}

// addHeader adds the table header row and the separator row that goes between
// the header row and the data.
func (m *MDTable) addHeader() {
	if m.HasHeaderRow {
		m.rowToMD(m.Table[0])
		//remove the first row
		m.Table = append(m.Table[1:])
	} else {
		if m.useFormat {
			m.rowToMD(m.HeaderRow)
		}
	}
	m.appendHeaderSeparatorRow(len(m.Table[0]))
	m.md = append(m.md, []byte("  \n")...)
	return
}

// appendHeaderSeparator adds the configured column  separator
func (m *MDTable) appendHeaderSeparatorRow(cols int) {
	m.appendColumnSeparator()
	for i := 0; i < cols; i++ {
		var separator []byte

		if m.useFormat {
			switch m.columnAlignment[i] {
			case "left", "l":
				separator = mdLeftJustify
			case "center", "c":
				separator = mdCentered
			case "right", "r":
				separator = mdRightJustify
			default:
				separator = mdDontJustify
			}
		} else {
			separator = mdDontJustify
		}

		separator = append(separator, mdPipe...)

		m.md = append(m.md, separator...)
	}
	m.md = append(m.md, []byte("  ")...)
	return
}

// appendColumnSeparator appends a pip to the md array
func (m *MDTable) appendColumnSeparator() {
	m.md = append(m.md, mdPipe...)
}

// FormatFromFile loads the format file specified.
func (m *MDTable) formatFromFile() error {
	// not really considering this an error that stops things, just one
	// that requirs error level logging. Is this right?
	if m.formatType != "file" {
		log.Printf("formatFromFile: nothing to do, formatType was %s, expected file", m.formatType)
		return nil
	}
	// if formatSource isn't set, nothing todo
	if m.formatSource == "" {
		log.Printf("formatFromFile: nothing to do, formatSource was not set", m.formatType)
		return nil
	}
	// Read from the format file
	table, err := ReadCSVFile(m.formatSource)
	if err != nil {
		log.Print(err)
		return err
	}
	//Row 0 is the header information
	m.HeaderRow = table[0]
	m.columnAlignment = table[1]
	m.columnEmphasis = table[2]
	return nil
}

// Write saves the md table as a source.md; the original extension of source is replaced
// by ,md, markdown, for the markdown output.
func (m *MDTable) Write() (n int, err error) {
	// figure out the output filename using source
	dname, fname := filepath.Split(m.source.Path)
	fparts := strings.Split(fname, ".")
	if len(fparts) == 1 {
		fname = fparts[0] + ".md"
	} else {
		// the last part is the extention
		fparts[len(fparts)-1] = "md"
		fname = strings.Join(fparts, ".") // first rejoin the name, using the md ext
	}
	f, err := os.OpenFile(filepath.Join(dname, fname), os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC, 0640)
	if err != nil {
		return 0, err
	}
	n, err = f.Write(m.md)
	if err != nil {
		return n, err
	}
	err = f.Close()
	if err != nil {
		return n, err
	}
	return n, nil
}
