tomd
====

tomd: a package that converts stuff **to ma**rkdown

Written to convert .csv files to markdown tables, but csvtomd has too many consonants and not enough vowels for most English speakers, or American speaker in this case. Anyways that kind of consonant/verb ratio frightens me.

So tomd it is!

## What tomd converts.
Not much at the moment:

* csv to markdown tables
  * csv files with a header row are supported.
  * csv file without a header row is supported, partially (complete support requires adding template support)

## How to use it:
### install
    $ go get github.com/mohae/tomd

### convert a CSV file:
Get a CSV object:
    c := tomd.NewCSV()

Call `CSV.FileToTable(filename)` with the CSV filename:
	err := c.FileToTable

Retrieve the markdown []bytes:
    md = c.MD()

### Use a io.Reader
Get a CSV object:
    c := tomd.NewCSV()

Call `CSV.ToTable(io.Reader)` with a io.Reader:
    err := c.ToTable(r)

Retrieve the markdown []bytes:
    md = c.MD()


### io.Reader Example:
```
package main

import (
        "fmt"
        "strings"

        "github.com/mohae/tomd"
)

var csvFile = strings.NewReader(`Author,Title,Year,ISBN
"Douglas Adams","Restaurant at the End of the Universe","1980","ISBN 0-345-39181-0"
"William Gibson","Burning Chrome","1986","ISBN 978-0-06-053982-5"`)

func main() {
        // get a new CSV object, with its defaults set
        c := tomd.NewCSV()

        // convert the csv to a MD table
        err := c.ToTable(csvFile)
        if err != nil {
                fmt.Println(err.Error())
        }
        md := c.MD()
        fmt.Printf("%s" , md)
}
```
## Format support
Basic support for formatting tables has been implemented. Formatting can either be done with a format file, `.fmt`, or by setting the CSV objects formatting information yourself.

The format information mimics a 3 row CSV table, with the number of columns matching the number of columns of the source CSV:

* Row 1 is the header row. The header information of the file for which this is a format must be set in the format file, regardless of its presence in the source data.
* Row 2 is the alignment information for the columns in the table.
* Row 3 is the emphasis information for the columns in the table.

### From a file:
Formatting from a file can be done two ways, either through the CSV.HasFormat flag or by passing the format file name, along with the csv file name, to the CSV.FileToTable() method. 

To pass the format file name the CSV.FileToTable() call must be in the form of:

    CSV.FileToTable(sourceFilename, formatFilename)

Passing the filename will also set the CSV.HasFormat to `true`.

If your format file uses the same path and name as your source file, except that it ends with `.fmt` instead of `.csv`, you can just set the CSV.HasFormat to `true`. When CSV.HasFormat is true and no format filename was passed, the format filename is derived from the source filename.

### Setting it yourself:
__Not fully implemented__
To set the formatting information without a file, you can set the information yourself. `CSV.HasFormat` must be set to true. The following fields need to be populated with their values:

    CSV.HeaderRow           []string
    CSV.ColumnAlignment     []string
    CSV.ColumnEmnphasis     []string //not implemented

### ColumnAlignment
Column alignment sets the column's alignment, if any. Valid values are:

* left, l
* right, r
* center, c

An empty value means the alignment is not set for that column

### CaolumnEmphasis
__Not Implemented__
Column emphasis sets the MD emphasis, if any, for that column. Valid values are:

* bold, b
* emphasis, e
* italics, i
* strikethrough, s

### Wishlist:
* format support for _italics_, __bold__, ~~strikethrough~~ 
