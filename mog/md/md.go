package mogger

import (
	"encoding/csv"
	"fmt"
	"io"
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
	mdLeftJustify []byte = []byte(":--")

	// RightJustify is the Md for right justification of columns,
	mdRightJustify []byte = []byte("--:")
	mdCentered []byte = []byte(":--:")
	mdDontJustify []byte = []byte("--")
)

// MDTable is a struct for representing and working with markdown tables
type MDTable struct {
	// producer/consumer information.
	mogger

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

	// whether for formatSource was autoset or not.
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

	// table is the parsed csv data
	table [][]string

	// md holds the md representation of the csv data
	md []byte
}

// NewCSV returns an initialize CSV object. It still needs to be configured
// for use.
func NewCSV() *CSV {
	C := &CSV{
		hasHeaderRow: true,
		destinationType: "bytes",
		table: [][]string{},
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

// MarshalTable takes a reader for csv and converts the read csv to a markdown
// table. To get the md, call CSV.md()
func (c *CSV) MarshalTable(r io.Reader) ([]byte, error) {
	var err error
	c.table, err = ReadCSV(r)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	//Now convert the data to md
	c.toMD()
	return c.md, nil
}

// MarshalTable marshals the CSV info to a markdown table.
func (c *CSV) MarshalTable() error{
	logger.Debugf("MarshalTable enter, source: %s", c.source)
	var err error
	// Try to read the source
	c.table, err = ReadCSVFile(c.source)
	if err != nil {
		logger.Error(err)
		return err
	}
		
	var formatName string
	// otherwise see if  HasFormat
	if c.useFormat {
//		c.setFormatFile()
		if c.formatType == "file" {
			//derive the format filename
			filename := filepath.Base(c.source)
			if filename == "." {
				err = fmt.Errorf("unable to determine format filename")
				logger.Error(err)
				return err
			}
	
			dir := filepath.Dir(c.source)
			parts := strings.Split(filename, ".")
			formatName = parts[0] + ".fmt"
			if dir != "." {
				formatName = dir + formatName
			}
		}
	}
	
	if c.useFormat {
		err := c.formatFromFile()
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	// Now convert the data to md
	c.toMD()

	logger.Debug("FileToMDTable exit with error: nil")
	return nil
}

// MD() returns the markdown as []byte
func (c *CSV) MD() []byte {
	return c.md
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

// setFormatSource sets the formatSource if it is not already set or if the 
// previously set value was set by setFormatSource. The latter allows auto-
// generated default source name to be updated when the source is while 
// preserving overrides.
func (c *CSV) autosetFormatFile() error {
	// if the source isn't set, nothing to do.
	if c.source == "" {
		logger.Trace("setFormatSource exit: source not set")
		return nil
	}

	// if formatSource isn't empty and wasn't set by setFormatSource,
	// nothing to do
	if c.formatSource != "" && !c.formatSourceAutoset {
		logger.Infof("setFormatSource exit: formatSource was already set to %s", c.formatSource)
		return nil
	}

	if c.formatType != "file" {
		logger.Trace("setFormatSource exit: not using format file, format type is %s", c.formatType)
		return nil
	}

	// Figure out the filename
	dir, file := filepath.Split(c.source)
	
	// break up the filename into its part, the last is extension.
	var fname string
	fParts := strings.Split(file, ".")

	if len(fParts) <= 2 {
		fname = fParts[0]
	} else {
		// Join all but the last part together for the name
		// This handles names with multiple `.`
		fname = strings.Join(fParts[0:len(fParts) - 2], ".")
	}

	fname += ".md"
	c.formatSource = dir + fname	
	c.formatSourceAutoset = true
	return nil
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
func (c *CSV) formatFromFile() error {
	// not really considering this an error that stops things, just one
	// that requirs error level logging. Is this right?
	if c.formatType != "file" {
		logger.Error("formatFromFile: nothing to do, formatType was %s, expected file", c.formatType)
		return nil
	}

	// if formatSource isn't set, nothing todo
	if c.formatSource == "" {
		logger.Error("formatFromFile: nothing to do, formatSource was not set", c.formatType)
		return nil
	}

	// Read from the format file
	table, err := ReadCSVFile(c.formatSource)
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

// Source returns the source of the CSV
func (c *CSV) Source() string {
	return c.source
}

// SetSource sets the source and has the formatFile updated, if applicable.
func (c *CSV) SetSource(s string) {
	c.source = s
	c.autosetFormatFile() 
}

// Destination is the of destination for the output, if applicable.
func (c *CSV) Destination() string {
	return c.destination
}

// SetDestination sets the destination of the output, if applicable.
func (c *CSV) SetDestination(s string) {
	c.destination = s
}

// DestinationType is the type of destination for the output.
func (c *CSV) DestinationType() string {
	return c.destinationType
}

// SetDestinationType sets the destinationType.
func (c *CSV) SetDestinationType(s string) {
	c.destinationType = s
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

// SetUseFormat sets whether a format should be used. This triggers a setting
// of the FormatFilename, if applicable.
func (c *CSV) SetUseFormat(b bool) {
	c.useFormat = b
	c.autosetFormatFile()
}
