package mog

import (
	"encoding/csv"
	_ "fmt"
	"io"
	"os"
	_ "path/filepath"
	_ "strconv"
	_ "strings"
)

// CSV is a struct for representing and working with csv data.
type CSV struct {
	// producer information.
	producer resource

	// consumer information
	consumer resource

	// format information
	format resource

	// Variables consistent with stdlib's Reader struct in the csv package,
	// with the exception of csv.Reader.TrailingComma, which is ommitted
	// since it is ignored.
	//
	// Anything set here will override the Reader's default value set by
	// csv.NewReader(). Please check golang.org/pkg/encoding/csv for more
	// info about the variables.
	comma		 rune
	comment		 rune
	fieldsPerRecord  int
	lazyQuotes	 bool
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
		producer: resource{},
		consumer: resource{},
		format: resource{},
		hasHeaderRow: true,
		table: [][]string{},
	}
	return C
}

// ReadCSV takes a reader, and reads the data connected with it as CSV data.
// A slice of slice of type string, or an error, are returned. This reads the
// entire file, so if the file is very large and you don't have sufficent RAM
// you will not like the results. There may be a row or chunk oriented
// implementation in the future.
func ReadCSV(r io.Reader ) ([][]string, error) {
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
