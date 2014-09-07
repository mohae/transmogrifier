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
	_ "strconv"
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
	// Supported:
	//	[]byte	no destination needed
	//	file	destination optional, if not set the output will be
	//		`sourceFilename.md` instead of `sourceFilename.csv`.
	destinationType string

	// hasHeaderRows: whether the csv data includes a header row as its
	// first row. If the csv data does not include header data, the header
	// data must be provided via template, e.g. false implies 
	// 'useFormat' == true. True does not have any implications on using
	// the format file.
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

	// formatSource: the location and name of the source file to use. It
	// can either be explicitely set, or TOMD will look for it as
	// `source.fmt`, for `source.csv`.
	formatSource string

	// useFormat: whether there's a format to use with the CSV or not. For
	// files, this is usually a file, with the same name and path as the
	// source, using the 'fmt' extension. This can also be set explicitely.
	// 'useFormat' == false implies 'hasHeaderRow' == true.
	useFormat bool

	// useFormatFile: whether the Format to use is a format file or not. 
	// When true, the format will be loaded from the file. When false, the
	// format information must be set using their setters: headerRow, 
	// columAlignment, and columnEmphasis.
	useFormatFile bool

	// table is the parsed csv data
	table [][]string

	// md holds the md representation of the csv data
	md []byte
}

// NewCSV returns an initialize CSV object. It still needs to be configured
// for use.
func NewCSV() *CSV {
	C := &CSV{
		hasHeader: true,
		destinationType: "bytes",
		table: [][]string{}
	}
	return C
}

// ToMDTable takes a reader for csv and converts the read csv to a markdown
// table. To get the md, call CSV.md()
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
	if c.useFormat {
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
	
	if c.useFormat {
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
		if c.useFormat {
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

		if c.useFormat {
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

// Destination is the of destination for the output, if applicable.
func (c *CSV) DestinationType() string {
	return c.destinationType
}

// SetDestination sets the destination of the output, if applicable.
func (c *CSV) SetDestinationType(string) {
	return c.destinationType
}

// DestinationType is the type of destination for the output.
func (c *CSV) DestinationType() string {
	return c.destinationType
}

// SetDestinationType sets the destinationType.
func (c *CSV) SetDestinationType(string) {
	return c.destinationType
}

// HasHeaderRow returns whether, or not, this csv file has a format file to
// use.
func (c *CSV) HasHeaderRow() bool {
	return c.hasHeaderRow
}

// SetHasHeaderRow sets whether, or not, the source has a header row.
func (c *CSV) SetHasHeaderRow(b bool) {
	c.hasHeaderRow = b
}

// HeaderRow returns the column headers; i.e., the header row.
func (c *CSV) HeaderRow() []string {
	return c.headerRow
}

// SetHeaderRow sets the headerRow information.
func (c *CSV) SetHeaderRow(s []string) {
	c.headerRow = s
}

// ColumnAlignment returns the columnAlignment information. This can be set
// either explicitely or using a format file.
func (c *CSV) ColumnAlignment() []string {
	return c.columnAlignment
}

// SetColumnAlignment sets the columnAlignment informatin.
func (c *CSV) SetColumnAlignment(s []string) {
	c.columnAlignment = s
}

// ColumnEmphasis returns the columnEmphasis information. This can be set
// either explicitly or with a format file.
func (c *CSV) ColumnEmphasis() []string {
	return c.columnEmphasis
}

// SetColumnEmphasis sets columnEmphasis information.
func (c *CSV) SetColumnEmphasis(s []string) {
	c.columnEmphasis = s
}

// FormatSource returns the formatSource information.
func (c *CSV) FormatSource() string {
	return c.formatSource
}

// SetFormatSource sets formatSource information. A side-affect of this is that
// setting the format file will automatically set `useFormat` and
// `useFormatFile`.
func (c *CSV) SetFormatSource(s string) {
	c.formatSource = s
}

// UseFormat returns whether this csv file has a format file to use.
func (c *CSV) UseFormat() bool {
	return c.useFormat
}

// SetUseFormat sets whether a format should be used.
func (c *CSV) SetUseFormat(b bool) {
	c.useFormat = b
}

// UseFormat returns whether this csv file has a format file to use.
func (c *CSV) UseFormatFile() bool {
	return c.useFormat
}

// SetUseFormatFile sets whether a format file should be used.
func (c *CSV) SetUseFormatFile(b bool) {
	c.useFormatFile = b
}

