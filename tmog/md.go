package mogger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	_ "strconv"
	"strings"

	"github.com/mohae/transmogrifier/mog"
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
	// producer/consumer information.
	//mogger
	Source mog.Resource
	// hasHeaderRows: whether the csv data includes a header row as its
	// first row. If the csv data does not include header data, the header
	// data must be provided via template, e.g. false implies
	// 'useFormat' == true. True does not have any implications on using
	// the format file.
	HasHeaderRow bool

	// headerRow contains the header row information. This is when a format
	// has been supplied, the header row information is set.
	HeaderRow []string

	// columnAlignment contains the alignment information for each column
	// in the table. This is supplied by the format
	ColumnAlignment []string

	// columnEmphasis contains the emphasis information, if any. for each
	// column. This is supplied by the format.
	ColumnEmphasis []string

	// formatSource: the location and name of the source file to use. It
	// can either be explicitely set, or TOMD will look for it as
	// `source.fmt`, for `source.csv`.
	FormatSource string

	// whether for formatSource was autoset or not.
	FormatSourceAutoset bool

	// useFormat: whether there's a format to use with the CSV or not. For
	// files, this is usually a file, with the same name and path as the
	// source, using the 'fmt' extension. This can also be set explicitely.
	// 'useFormat' == false implies 'hasHeaderRow' == true.
	UseFormat bool

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
	FormatType string

	// table is the parsed csv data
	Table [][]string

	// md holds the md representation of the csv data
	MD []byte
}

func NewMDTable() *MDTable {
	return &MDTable{HasHeaderRow: true, HeaderRow: []string{}, ColumnAlignment: []string{}, ColumnEmphasis: []string{}, Table: [][]string{}, MD: []byte{}}
}

// TMogCSVTableReader takes an io.Reader, which contains a CSV Table, and
// tranmogrifies it to a MD table, which is returned, unless an error occurrs.
func (m *MDTable) TMogCSVTableReader(r io.Reader) ([]byte, error) {
	var err error
	m.Table, err = mog.ReadCSV(r)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	//Now convert the data to md
	m.toMD()
	return m.MD, nil
}

// TmogCSVTable marshals the CSV within the struct transmogrified to a MD
// table. The results are stored in m.md.
func (m *MDTable) TmogCSVTable() error {
	log.Printf("MarshalTable enter, source: %s", m.Source.Path)
	var err error
	// Try to read the source
	m.Table, err = mog.ReadCSVFile(m.Source.Path)
	if err != nil {
		log.Print(err)
		return err
	}

	var formatName string
	// otherwise see if  HasFormat
	if m.UseFormat {
		//		c.setFormatFile()
		if m.FormatType == "file" {
			//derive the format filename
			filename := filepath.Base(m.Source.Path)
			if filename == "." {
				err = fmt.Errorf("unable to determine format filename")
				log.Print(err)
				return err
			}

			dir := filepath.Dir(m.Source.Path)
			parts := strings.Split(filename, ".")
			formatName = parts[0] + ".fmt"
			if dir != "." {
				formatName = dir + formatName
			}
		}
	}

	if m.UseFormat {
		err := m.formatFromFile()
		if err != nil {
			log.Print(err)
			return err
		}
	}

	// Now convert the data to md
	m.toMD()
	return nil
}

// toMD converts the table to markdown
func (m *MDTable) toMD() {
	// Process the header first
	m.addHeader()

	// for each row of table data, process it.
	for _, row := range m.Table {
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
		m.MD = append(m.MD, bcol...)
		m.appendColumnSeparator()
	}
	// add a new line at the end of a row
	m.MD = append(m.MD, []byte("  \n")...)

}

// addHeader adds the table header row and the separator row that goes between
// the header row and the data.
func (m *MDTable) addHeader() {
	if m.HasHeaderRow {
		m.rowToMD(m.Table[0])
		//remove the first row
		m.Table = append(m.Table[1:])
	} else {
		if m.UseFormat {
			m.rowToMD(m.HeaderRow)
		}
	}

	m.appendHeaderSeparatorRow(len(m.Table[0]))
	m.MD = append(m.MD, []byte("  \n")...)
	return
}

// appendHeaderSeparator adds the configured column  separator
func (m *MDTable) appendHeaderSeparatorRow(cols int) {
	m.appendColumnSeparator()

	for i := 0; i < cols; i++ {
		var separator []byte

		if m.UseFormat {
			switch m.ColumnAlignment[i] {
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

		m.MD = append(m.MD, separator...)
	}
	m.MD = append(m.MD, []byte("  ")...)
	return

}

// appendColumnSeparator appends a pip to the md array
func (m *MDTable) appendColumnSeparator() {
	m.MD = append(m.MD, mdPipe...)
}

// FormatFromFile loads the format file specified.
func (m *MDTable) formatFromFile() error {
	// not really considering this an error that stops things, just one
	// that requirs error level logging. Is this right?
	if m.FormatType != "file" {
		log.Printf("formatFromFile: nothing to do, formatType was %s, expected file", m.FormatType)
		return nil
	}

	// if formatSource isn't set, nothing todo
	if m.FormatSource == "" {
		log.Printf("formatFromFile: nothing to do, formatSource was not set", m.FormatType)
		return nil
	}

	// Read from the format file
	table, err := mog.ReadCSVFile(m.FormatSource)
	if err != nil {
		log.Print(err)
		return err
	}

	//Row 0 is the header information
	m.HeaderRow = table[0]
	m.ColumnAlignment = table[1]
	m.ColumnEmphasis = table[2]

	return nil
}

// Write saves the md table as a source.md; the original extension of source is replaced
// by ,md, markdown, for the markdown output.
func (m *MDTable) Write() (n int, err error) {
	// figure out the output filename using source
	dname, fname := filepath.Split(m.Source.Path)
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

	n, err = f.Write(m.MD)
	if err != nil {
		return n, err
	}
	err = f.Close()
	if err != nil {
		return n, err
	}

	return n, nil
}
