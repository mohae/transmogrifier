package transmogrifier

import (
	"fmt"
	_ "io"
	"os"
	"path/filepath"
	_ "strconv"
	"strings"
)

// MDTable format representations.
var (
	// Pipe is the MD column separator
	mdPipe = `|`
	// LeftJustify is the MD for left justification of columns.
	mdLeftJustify = []byte(":---")
	// RightJustify is the Md for right justification of columns,
	mdRightJustify = []byte("---:")
	mdCentered     = []byte(":---:")
	mdDontJustify  = []byte("---")
)

// MDTable is a struct for representing and working with markdown tables
type MDTable struct {
	// data source, if applicable.
	source       resource
	sourceFormat FormatType
	// sink, if applicable
	dest resource
	// FormatSource: the location and name of the source file to use. It
	// can either be explicitely set, or TOMD will look for it as
	// `source.fmt`, for `source.csv`.
	formatSource string
	// useFormat: whether or not a format file should be used for the table.
	// 'useFormat' == false implies 'hasHeaderRow' == true.
	useFormat bool
	// hasColumnNames
	hasColumnNames bool
	// columnNames contains the name of each table column
	columnNames []string
	// columnAlignment contains the alignment information, if any, for each
	// column.  This is supplied by the format.
	columnAlignment []string
	// columnEmphasis contains the emphasis information, if any. for each column.
	// This is supplied by the format.
	columnEmphasis []string
	// md is the md table, in bytes
	md []byte
}

// NewMDTable returns an empty MDTable struct.
func NewMDTable() *MDTable {
	return &MDTable{columnNames: []string{}, columnAlignment: []string{}, columnEmphasis: []string{}, md: []byte{}}
}

func (m MDTable) String() string {
	return string(m.md)
}

func (m MDTable) Bytes() []byte {
	return m.md
}

// SetSource set's the Source to the passed value.  The source's extension is
// checked to see if it is a supported format.  An error is returned if it
// isn't.
func (m *MDTable) SetSource(s string) error {
	if s == "" {
		return fmt.Errorf("source string was empty")
	}
	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return fmt.Errorf("unable to determine format of %q", s)
	}
	m.sourceFormat = FormatTypeFromString(parts[len(parts)-1])
	if m.sourceFormat == FmtUnsupported {
		return fmt.Errorf("unsupported format for %q: %q", s, parts[len(parts)-1])
	}
	m.source = NewResource(s, m.sourceFormat, File)
	return nil
}

// SetUseFormat: whether or not a format should be applied to the MD table.
func (m *MDTable) SetUseFormat(b bool) {
	m.useFormat = b
	m.SetHasColumnNames(b)
}

// SetFormatSource set's the source of the format information and sets
// useFormat to 'true'.  If the formatSource != "", it will be used as the
// location of the formatting information for the MD Table. If it isn't set and
// useFormat == true, the format source is expected to be in the same location
// as the source, with the same name + an extension of '.fmt'.
func (m *MDTable) SetFormatSource(s string) error {
	if s == "" {
		return fmt.Errorf("unable to set format source: received empty string")
	}
	m.formatSource = s
	m.useFormat = true
	return nil
}

// SetHasColumnNames
func (m *MDTable) SetHasColumnNames(b bool) {
	m.hasColumnNames = b
}

// SetColumnNames
func (m *MDTable) SetColumnNames(cols []string) {
	m.columnNames = make([]string, len(cols))
	copy(m.columnNames, cols)
}

func (m *MDTable) SetColumnAlignment(cols []string) {
	m.columnAlignment = make([]string, len(cols))
	copy(m.columnAlignment, cols)
}

func (m *MDTable) SetColumnEmphasis(cols []string) {
	m.columnEmphasis = make([]string, len(cols))
	copy(m.columnEmphasis, cols)
}

// Transmogrify transomgrifies the source into a MD table. The result is held
// in md and can be obtained by m.MD().  Any error encountered is returned.
// SetHasHeader needs to be called prior to calling this method.
func (m *MDTable) TransmogrifyStringTable(t [][]string) error {
	if t == nil {
		return fmt.Errorf("unable to tranmogrify string table: received nil")
	}

	// Process the header first
	if m.hasColumnNames {
		m.SetColumnNames(t[0])
		//remove the first row
		t = t[1:]
	}
	m.tableHeader()
	// for each row of table data, process it.
	for _, row := range t {
		m.rowToMD(row)
	}
	return nil
}

func (m *MDTable) tableHeader() {
	useFormat := m.useFormat
	m.useFormat = false
	m.rowToMD(m.columnNames)
	m.appendHeaderSeparatorRow()
	m.useFormat = useFormat
}

// rowTomd takes a table row and returns the md version of it consistent
// with its configuration.
func (m *MDTable) rowToMD(cols []string) {
	m.appendColumnSeparator()
	for i, col := range cols {
		var bcol []byte
		// TODO this is where column data decoration would occur
		// with templates
		if m.useFormat {
			switch m.columnEmphasis[i] {
			case "bold", "b":
				bcol = append(bcol, []byte{'_', '_'}...)
			case "italic", "italics", "i":
				bcol = append(bcol, []byte{'_'}...)
			case "strikethrough", "s":
				bcol = append(bcol, []byte{'~', '~'}...)
			}
		}
		bcol = append(append(bcol, []byte(col)...), bcol...)
		m.md = append(m.md, bcol...)
		m.appendColumnSeparator()
	}
	// add a new line at the end of a row
	m.md = append(m.md, []byte("  \n")...)
}

// appendHeaderSeparator adds the configured column  separator
func (m *MDTable) appendHeaderSeparatorRow() {
	m.appendColumnSeparator()
	for i := 0; i < len(m.columnNames); i++ {
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
	m.md = append(m.md, []byte("  \n")...)
	return
}

// appendColumnSeparator appends a pip to the md array
func (m *MDTable) appendColumnSeparator() {
	m.md = append(m.md, mdPipe...)
}

// FormatFromFile loads the format file specified.
func (m *MDTable) formatFromFile() error {
	// if formatSource isn't set, nothing todo
	if m.formatSource == "" {
		return nil
	}
	// Read from the format file
	fsource := NewCSV()
	fsource.SetHasHeader(false)
	err := fsource.ReadFile(m.formatSource)
	if err != nil {
		return err
	}
	if len(fsource.rows) < 3 {
		return fmt.Errorf("insufficient format rows: expected at least 3, got %d", len(fsource.rows))
	}
	//Row 0 is the header information
	m.columnNames = append(m.columnNames, fsource.rows[0]...)
	//Row 1 is the column alignment information
	m.columnAlignment = append(m.columnAlignment, fsource.rows[1]...)
	//Row 2 is the column emphasis information
	m.columnEmphasis = append(m.columnEmphasis, fsource.rows[2]...)
	return nil
}

// SetDest sets the destination of the Write operation. If the destination is
// an empty string, "", the source name will be concatinated with '.md'. Any
// non-empty dest string will be used as the destination.
//
// Currently, only write to file is supported.
func (m *MDTable) SetDest(s string) {
	m.dest = NewResource(s, FmtMDTable, File)
	if s != "" {
		return
	}
	m.dest.SetPath(m.source.Path)
	m.dest.SetName(mdFilenameFrom(m.source.Name))
}

// TODO add support for writing to the received writer
//
//

//
// Currentky
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

func mdFilenameFrom(source string) string {
	if source == "" {
		return ""
	}
	parts := strings.Split(source, ".")
	if len(parts) < 2 {
		return fmt.Sprintf("%s.md", parts[0])
	}
	var dest string
	for i, part := range parts {
		if i == len(parts)-1 {
			dest += "md"
			return dest
		}
		dest += part + "."
	}
	return dest

}
