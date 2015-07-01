package transmogrifier

import (
	"encoding/csv"
	"io"
	"os"
	_ "path/filepath"
	_ "strconv"
	_ "strings"
)

// CSV is a struct for representing and working with csv data.
type CSV struct {
	// source information.
	source resource
	// sink information
	sink resource
	// Variables consistent with stdlib's Reader struct in the csv package,
	// with the exception of csv.Reader.TrailingComma, which is ommitted
	// since it is ignored.
	//
	// Anything set here will override the Reader's default value set by
	// csv.NewReader(). Please check golang.org/pkg/encoding/csv for more
	// info about the variables.
	comma            rune
	comment          rune
	fieldsPerRecord  int
	lazyQuotes       bool
	trimLeadingSpace bool
	// hasHeader: whether the csv data includes a header row as its
	// first row. If the csv data does not include header data, the header
	// data must be provided via template, e.g. false implies
	// 'useFormat' == true. True does not have any implications on using
	// the format file.
	hasHeader bool
	// The csv file data:
	// headerRow contains the header row information. This is when a format
	// has been supplied, the header row information is set.
	headerRow []string
	// table is the parsed csv data
	rows [][]string
}

// NewCSV returns an initialize CSV object. It still needs to be configured
// for use.
func NewCSV() *CSV {
	C := &CSV{
		source:    resource{},
		sink:      resource{},
		hasHeader: true,
		headerRow: []string{},
		rows:      [][]string{},
	}
	return C
}

// NewCSVSource creates a new *CSV with its source set,
func NewCSVSource(s string) *CSV {
	c := NewCSV()
	c.SetSource(s)
	return c
}

// ReadAll takes a reader, and reads the data connected with it as CSV data.
// If there is a header row, CSV.hasHeader == true, the headerRow field is
// populated with the first row of the source.  This reads the entire file at
// once.  If an error occurs, it is returned
func (c *CSV) Read(r io.Reader) error {
	var err error
	cr := csv.NewReader(r)
	c.rows, err = cr.ReadAll()
	if err != nil {
		return err
	}
	if c.hasHeader {
		c.headerRow = c.rows[0]
		c.rows = c.rows[1:]
	}
	return nil
}

// ReadFile takes a path, reads the contents of the file and returns any error
// encountered. The entire file will be read at once.
func (c *CSV) ReadFile(f string) error {
	if f == "" {
		return ErrNoSource
	}
	file, err := os.Open(f)
	if err != nil {
		return err
	}
	// because we don't want to forget or worry about hanldling close prior
	// to every return.
	defer file.Close()
	// Read the file into csv
	return c.Read(file)
}

func (c *CSV) ReadSource() error {
	return c.ReadFile(c.source.String())
}

// SetSource sets the source and has the formatFile updated, if applicable.
func (c *CSV) SetSource(s string) {
	c.source = NewResource(s)
}

// Source returns the source string
func (c *CSV) Source() string {
	return c.source.String()
}

func (c *CSV) SetHasHeader(b bool) {
	c.hasHeader = b
}

// HasHeader returns the hasHeader bool
func (c *CSV) HasHeader() bool {
	return c.hasHeader
}

// HeaderRow returns the header row, if it exists.
func (c *CSV) HeaderRow() []string {
	return c.headerRow
}

// Rows returns the csv rows.
func (c *CSV) Rows() [][]string {
	return c.rows
}
