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

## TODO:
Add template support. Templates can define justification for columns, and, optionally, define header row data, column names, for a CSV file.

### Wishlist:
* support templating in a CSV cell.
* template support for _italics_, __bold__, ~~strikethrough~~ 
