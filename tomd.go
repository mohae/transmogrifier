// Copyright © 2014, All rights reserved
// Joel Scoble, https://github.com/mohae/car
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

	"github.com/cihub/seelog"
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

// CSVtoTable takes an stream of bytes that is CSV and outputs it as a stream
// of bytes representing that table in Markdown.
//
// Customized CSV to MD table translations can be specified using template
// files that specify how the CSV should be translated. Template files enable
// support of all the implemented features listed in the support list. If the
// feature has an _ in front of it. it is not supported yet.
//
// Supports:
//	_ Headers
//	_ No Headers
//	_ Right Justified
//	_ Left Justified
//	_ Centered
//	_ Justified Headers
//	_ Justified Fields
// 
// TODO: implement unsopported features
//

// CSV is a struct for representing and working with csv data.
type CSV struct {
	// WHeaders signifies whether the datasource's first row is a header
	// row or not. If true, the sources first row is a header row and will
	// be treated as such, otherwise the header information comes from the
	// template file for the source CSV.
	HasHeaderRow bool

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

	// Template is the name of the template to use. This is for justifying
	// the MD table. If no justificattion is wanted, leave template empty.
	// TODO: currently not supported
	Template string
	
	// table is the parsed csv data
	table [][]string

	md []byte
}

func NewCSV() *CSV {
	// Only explicitely set the defaults that are not consistent with the
	// variable types initialization state.
	C := &CSV{HasHeaderRow: true, destinationType: "bytes", table: [][]string{}}
	return C
}

// Table converts the incoming csv to a markdown table. It is expected that any
// other settings that are needed to get the desired MD table output to the
// correct destination will be set. 
// The crated markdown table is captured by CSV and is available through its 
// MD() method. If no destination is set, this is how the generated markdown
// can be retrieved.
func (c *CSV) FileToMDTable(source string) error {
	var err error
	//Get the CSV from the source
	c.table, err = ReadCSVFile(source)
	if err != nil {
		logger.Error(err)
		return err
	}
	
	// Now convert the data to MD
	
	return nil
}

// ReadCSV takes a reader, and reads the data connected with it as CSV data.
// A slice of slice of type string, or an error, are returned. This reads the
// entire file, so if the file is very large and you don't have sufficent RAM
// you will not like the results. There may be a row oriented implementation 
// in the future.
func ReadCSV(r io.Reader ) ([][]string,  error) {
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

// CSVtoMD takes a [][]string and returns a markdown representation of it as
// []byte
func (c *CSV) ToMD() ()  {
	// Process the header first
	c.addHeader()

	// for each row of table data, process it.
	for _, row := range c.table {
		c.RowToMD(row)
	}
	
	return
}

// RowToMD takes a csv table row and returns the md version of it consistent
// with its configuration.
func (c *CSV) RowToMD(cols []string) {
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
	if c.HasHeaderRow {
		c.RowToMD(c.table[0])
		//remove the first row
		c.table = append(c.table[1:])
	} else {
		// not implemented--get from template TODO
	}

	c.appendHeaderSeparatorRow(len(c.table[0]))
	return
}

// appendSeparatorRow adds a sepa
// appendSeparator adds the configured column  separator
func (c *CSV) appendHeaderSeparatorRow(cols int) {
	c.appendColumnSeparator()
	val := []byte("-|")
	for i := 0; i < cols; i++ {
		c.md = append(c.md, val...)
	}

	return
			
}

// appendColSeparator appends a pip to the md array
func (c *CSV) appendColumnSeparator() () {
	val := []byte("|")
	c.md = append(c.md, val...)
}
