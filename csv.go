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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// CSV is a struct for representing and working with csv data.
type CSV struct {
	// Source is the source of the CSV data. It is currently assumed to be
	// a path location
	source string

	// destination is where the generated markdown should be put, if it is
	// to be put anywhere. When used, this setting is used in conjunction 
	// with destinationType. Not all destinationTypes need to specify a
	// destinatin, bytes, for example.
	destination string

	// destinationType is the type of destination for the md, e.g. file.
	// If the destinationType requires specification of the destination,
	// the Destination variable should be set to that value.
	destinationType string

	// hasFormat: whether there's a format to use with the CSV or not. For
	// files, this is a file with the same name as the CSV file
	hasFormat bool

	// HasHeaderRows: whether the csv data includes a header row as its
	// first row. If the csv data does not include header data, the header
	// data must be provided via template
	hasHeaderRow bool

	// headerRow contains the header row information. This is when a format
	// has been supplied, the header row information is set.
	headerRow []string
	
	// columnAlignment contains the alignment information for each column
	// in the table. This is supplied by the format
	columnAlignment []string

	// columnEmphasis contains the emphasis information, if any. for each
	// column. This is supplied by the format.
	columnEmphasis []string

	// table is the parsed csv data
	table [][]string

	// md holds the md representation of the csv data
	md []byte
}

func NewCSV() *CSV {
	// Set the defaults based on the environment variables and defaults
	tmp := os.Getenv("hasheader")
	header, err := strconv.ParseBool(tmp)
	if err != nil {
		header = false
	}

	C := &CSV{hasHeaderRow: header, destinationType: "bytes", table: [][]string{}}
	return C
}

// ToMDTable takes a reader for csv and converts the read csv to a markdown
// table.
// To get the md, call CSV.md()
func (c *CSV) ToMDTable(r io.Reader) error {
	var err error
	c.table, err = ReadCSV(r)
	if err != nil {
		logger.Error(err)
		return err
	}

	//Now convert the data to md
	c.toMD()
	return nil
}

// FileToMDTable takes a file and marshals it to a md table.
func (c *CSV) FileToMDTable(source string) error{
	var err error
	// Try to read the source
	c.table, err = ReadCSVFile(source)
	if err != nil {
		logger.Error(err)
		return err
	}
		
	var formatName string
	// otherwise see if  HasFormat
	if c.hasFormat {
		//derive the format filename
		filename := filepath.Base(source)
		if filename == "." {
			err = fmt.Errorf("unable to determine format filename")
			logger.Error(err)
			return err
		}

		dir := filepath.Dir(source)
		parts := strings.Split(filename, ".")
		formatName = parts[0] + ".fmt"
		if dir != "." {
			formatName = dir + formatName
		}
	}
	
	if c.hasFormat {
		err := c.formatFromFile(formatName)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	// Now convert the data to md
	c.toMD()
	return nil
}

// md() returns the markdown as []byte
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

// tomd does table header processing then converts its table data to md,
func (c *CSV) toMD() ()  {
	// Process the header first
	c.addHeader()

	// for each row of table data, process it.
	for _, row := range c.table {
		c.rowToMD(row)
	}
	
	return
}

// rowTomd takes a csv table row and returns the md version of it consistent
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
		if c.hasFormat {
			c.rowToMD(c.headerRow)
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

		if c.hasFormat {
			switch c.columnAlignment[i] {
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
	
		c.md = append(c.md, separator...)
	}

	return
			
}

// appendColumnSeparator appends a pip to the md array
func (c *CSV) appendColumnSeparator() {
	c.md = append(c.md, mdPipe...)
}

// FormatFromFile loads the format file specified. 
func (c *CSV) formatFromFile(s string) error {
	table, err := ReadCSVFile(s)
	if err != nil {
		logger.Error(err)
		return err
	}
	
	//Row 0 is the header information
	c.headerRow = table[0]
	c.columnAlignment = table[1]
	c.columnEmphasis = table[2]

	return nil
}
