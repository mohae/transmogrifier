// Copyright © 2014, All rights reserved
// Joel Scoble, https://github.com/mohae/tomd
//
// This is licensed under The MIT License. Please refer to the included
// LICENSE file for more information. If the LICENSE file has not been
// included, please refer to the url above.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License
//
// tomd: to markdown, takes input and converts it to markdown
//
// Notes: 
//	* This is not a general markdown processor. It is a package to provide
//      functions that allow things to be converted to their representation
//      in markdown.
//      Currently that means taking a .csv file and converting it to a table.
//	* Uses seelog for 'library logging', to enable logging see:
//        http://github.com/cihub/seelog/wiki/Writing-libraries-with-Seelog
package tomd

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cihub/seelog"
)

// MD constants
var (
	// Pipe is the MD column separator
	pipe []byte = []byte("|")
	
	// LeftJustify is the MD for left justification of columns.
	leftJustify []byte = []byte(":--")

	// RightJustify is the Md for right justification of columns,
	rightJustify []byte = []byte("--:")
	centered []byte = []byte(":--:")
	dontJustify []byte = []byte("--")
)

var logger seelog.LoggerInterface

func init() {	
	// Disable logger by default
	DisableLog()
}

// DisableLog disables all library log output.
func DisableLog() {
	logger = seelog.Disabled
}

// UseLogger uses a specified seelog.LoggerInterface to output library log.
// Use this func if you are using Seelog logging system in your app.
func UseLogger(newLogger seelog.LoggerInterface) {
	logger = newLogger
}

// SetLogWriter uses a specified io.Writer to output library log.
// Use this func if you are not using the Seelog logging system in your app.
func SetLogWriter(writer io.Writer) error {
	if writer == nil {
		return errors.New("Nil writer")
	}

	newLogger, err := seelog.LoggerFromWriterWithMinLevel(writer, seelog.TraceLvl)
	if err != nil {
		return err
	}

	UseLogger(newLogger)	
	return nil
}

// FlushLog must be called before app shutdown.
func FlushLog() {
	logger.Flush()
}

// CSV is a struct for representing and working with csv data.
type CSV struct {
	// HasHeaderRows: whether the csv data includes a header row as its
	// first row. If the csv data does not include header data, the header
	// data must be provided via template
	hasHeaderRow bool

	// Source is the source of the CSV data. It is currently assumed to be
	// a path location
	Source string

	// Destination is where the generated markdown should be put, if it is
	// to be put anywhere. When used, this setting is used in conjunction 
	// with destinationType. Not all destinationTypes need to specify a
	// destinatin, bytes, for example.
	Destination string

	// DestinationType is the type of destination for the md, e.g. file.
	// If the destinationType requires specification of the destination,
	// the Destination variable should be set to that value.
	destinationType string

	// HasFormat: whether there's a format to use with the CSV or not. For
	// files, this is a file with the same name as the CSV file
	HasFormat bool

	// HeaderRow contains the header row information. This is when a format
	// has been supplied, the header row information is set.
	HeaderRow []string
	
	// ColumnAlignment contains the alignment information for each column
	// in the table. This is supplied by the format
	ColumnAlignment []string

	// ColumnEmphasis contains the emphasis information, if any. for each
	// column. This is supplied by the format.
	ColumnEmphasis []string

	// table is the parsed csv data
	table [][]string

	// md holds the md representation of the csv data
	md []byte
}

func NewCSV() *CSV {
	// Only explicitely set the defaults that are not consistent with the
	// variable types initialization state.
	C := &CSV{hasHeaderRow: true, destinationType: "bytes", table: [][]string{}}
	return C
}

// Table takes a reader for csv and converts the read csv to a markdown
// table.
// To get the md, call CSV.MD()
func (c *CSV) ToTable(r io.Reader) error {
	var err error
	c.table, err = ReadCSV(r)
	if err != nil {
		logger.Error(err)
		return err
	}

	//No convert the data to MD
	c.toMD()
	return nil
}

// FileToTable takes a filename, opens it, reads the csv, and then converts
// it to a markdown table.
// CSV.md is used to build the table markdown during processing. It also stores
// the md and is available through the CSV.MD() method, which returns []byte
// MD() method.
//
// CSV filename = sources[0]
// FMT filename = sources[1]
func (c *CSV) FileToTable(sources ...string) error {
	var err error

	//Get the CSV from the source
	c.table, err = ReadCSVFile(sources[0])
	if err != nil {
		logger.Error(err)
		return err
	}

	// If there's a format, load it,
	var formatName string
	if len(sources) >= 2 {
		// If the slice passed was long enough, assume 2nd is format file
		formatName =  sources[1]

		// and set HasFormat to true
		c.HasFormat = true
	} else {
		// otherwise see if  HasFormat
		if c.HasFormat {
			//derive the format filename
			filename := filepath.Base(sources[0])
			if filename == "." {
				err := errors.New("unable to determine format filename")
				logger.Error(err)
				return err
			}

			dir := filepath.Dir(sources[0])
			parts := strings.Split(filename, ".")
			formatName = parts[0] + ".fmt"

			if dir != "." {
				formatName = dir + formatName
			}
		}
	}
	
	if c.HasFormat {
		err := c.FormatFromFile(formatName)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	// Now convert the data to MD
	c.toMD()

	return nil
}

// MD() returns the markdown as []byte
func (c *CSV) MD() []byte {
	return c.md
}

// ReadCSV takes a reader, and reads the data connected with it as CSV data.
// A slice of slice of type string, or an error, are returned. This reads the
// entire file, so if the file is very large and you don't have sufficent RAM
// you will not like the results. There may be a row oriented implementation 
// in the future.
func ReadCSV(r io.Reader ) ([][]string, error) {
	cr := csv.NewReader(r)
	rows, err := cr.ReadAll()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return rows, nil
}

// ReadCSVFile takes a path, reads the contents of the file and returns int.
func ReadCSVFile(f string) ([][]string, error) {
	file, err := os.Open(f)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	
	// because we don't want to forget or worry about hanldling close prior
	// to every return.
	defer file.Close()
	
	//
	data, err := ReadCSV(file)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return data, nil
}

// ToMD does table header processing then converts its table data to MD,
func (c *CSV) toMD() ()  {
	// Process the header first
	c.addHeader()

	// for each row of table data, process it.
	for _, row := range c.table {
		c.rowToMD(row)
	}
	
	return
}

// RowToMD takes a csv table row and returns the md version of it consistent
// with its configuration.
func (c *CSV) rowToMD(cols []string) {
	c.appendColumnSeparator()

	for _, col := range cols {
		// TODO this is where column data decoration would occur
		// with templates
		bcol := []byte(col)
		c.md = append(c.md, bcol...)
		c.appendColumnSeparator()
	}

}

// addHeader adds the table header row and the separator row that goes between
// the header row and the data.
func (c *CSV) addHeader() () {
	if c.hasHeaderRow {
		c.rowToMD(c.table[0])
		//remove the first row
		c.table = append(c.table[1:])
	} else {
		if c.HasFormat {
			c.rowToMD(c.HeaderRow)
		}
	}

	c.appendHeaderSeparatorRow(len(c.table[0]))
	return
}

// appendHeaderSeparator adds the configured column  separator
func (c *CSV) appendHeaderSeparatorRow(cols int) {
	c.appendColumnSeparator()

	for i := 0; i < cols; i++ {
		var separator []byte	

		if c.HasFormat {
			switch c.ColumnAlignment[i] {
			case "left", "l":
				separator = leftJustify
			case "center", "c":
				separator = centered
			case "right", "r":
				separator = rightJustify
			default:
				separator = dontJustify
			}
		} else {
			separator = dontJustify
		}

		separator = append(separator, pipe...)
	
		c.md = append(c.md, separator...)
	}

	return
			
}

// appendColumnSeparator appends a pip to the md array
func (c *CSV) appendColumnSeparator() {
	c.md = append(c.md, pipe...)
}

// FormatFromFile loads the format file specified. 
func (c *CSV) FormatFromFile(s string) error {
	table, err := ReadCSVFile(s)
	if err != nil {
		logger.Error(err)
		return err
	}
	
	//Row 0 is the header information
	c.HeaderRow = table[0]
	c.ColumnAlignment = table[1]
	c.ColumnEmphasis = table[2]

	return nil
}
