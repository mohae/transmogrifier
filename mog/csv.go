package mog

import (
	"encoding/csv"
	_ "fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	_ "strconv"
	"strings"
)

// CSV is a struct for representing and working with csv data.
type CSV struct {
	// source information.
	source Resource

	// sink information
	sink Resource

	// format information
	format Resource

	// Automatically set the format source-or not.
	formatSourceAutoSet bool

	// useFormat use a format file
	useFormat bool

	// formatType is the type of format being used
	formatType string

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

	// format
	hasFormat bool

	// hasHeaderRows: whether the csv data includes a header row as its
	// first row. If the csv data does not include header data, the header
	// data must be provided via template, e.g. false implies
	// 'useFormat' == true. True does not have any implications on using
	// the format file.
	hasHeaderRow bool

	// The csv file data:
	// headerRow contains the header row information. This is when a format
	// has been supplied, the header row information is set.
	headerRow []string

	// table is the parsed csv data
	table [][]string
}

// NewCSV returns an initialize CSV object. It still needs to be configured
// for use.
func NewCSV() *CSV {
	C := &CSV{
		source:       Resource{},
		sink:         Resource{},
		format:       Resource{},
		hasHeaderRow: true,
		table:        [][]string{},
	}
	return C
}

// NewSourceCSV creates a new *CSV with its source set and initialized.
func NewSourcesCSV(t, s string, b bool) *CSV {
	c := NewCSV()
	c.useFormat = b

	// currently anything that's not file uses the default "", which
	// means set it yourself to use it.
	switch t {
	case "bytes":
	case "file":
		c.formatType = "file"
	}

	c.SetSource(s)
	return c
}

// ReadCSV takes a reader, and reads the data connected with it as CSV data.
// A slice of slice of type string, or an error, are returned. This reads the
// entire file, so if the file is very large and you don't have sufficent RAM
// you will not like the results. There may be a row or chunk oriented
// implementation in the future.
func ReadCSV(r io.Reader) ([][]string, error) {
	cr := csv.NewReader(r)
	rows, err := cr.ReadAll()
	if err != nil {
		//		logger.Error(err)
		return nil, err
	}

	return rows, nil
}

// ReadCSVFile takes a path, reads the contents of the file and returns int.
func ReadCSVFile(f string) ([][]string, error) {
	file, err := os.Open(f)
	if err != nil {
		//		logger.Error(err)
		return nil, err
	}

	// because we don't want to forget or worry about hanldling close prior
	// to every return.
	defer file.Close()

	// Read the file into csv
	csv, err := ReadCSV(file)
	return csv, nil
}

// SetSource sets the source and has the formatFile updated, if applicable.
func (c *CSV) SetSource(s string) {
	c.source = Resource{Name: s}
	c.autoSetFormatFile()
}

// autoSetFormatSource sets the formatSource if it is not already set or if the
// previously set value was set by setFormatSource. The latter allows auto-
// generated default source name to be updated when the source is while
// preserving overrides.
func (c *CSV) autoSetFormatFile() error {
	// if the source isn't set, nothing to do.
	if c.source.Name == "" {
		log.Printf("setFormatSource exit: source not set")
		return nil
	}

	// if formatSource isn't empty and wasn't set by setFormatSource,
	// nothing to do
	if c.source.Format != "" && !c.formatSourceAutoSet {
		log.Printf("setFormatSource exit: formatSource was already set to %s", c.source.Format)
		return nil
	}

	if c.source.Type != "file" {
		log.Printf("setFormatSource exit: not using format file, format type is %s", c.source.Type)
		return nil
	}

	// Figure out the filename
	dir, file := filepath.Split(c.source.Name)

	// break up the filename into its part, the last is extension.
	var fname string
	fParts := strings.Split(file, ".")

	if len(fParts) <= 2 {
		fname = fParts[0]
	} else {
		// Join all but the last part together for the name
		// This handles names with multiple `.`
		fname = strings.Join(fParts[0:len(fParts)-2], ".")
	}

	fname += ".md"
	c.source.Path = filepath.Join(dir, fname)
	c.formatSourceAutoSet = true
	return nil
}
