package transmogrifier

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type JSON struct {
	Source resource
	Table  [][]string
	B      []byte
}

func NewJSON() *JSON {
	return &JSON{Table: [][]string{}}
}

func (j *JSON) TmogCSVTableReader(r io.Reader) ([]byte, error) {
	var err error
	j.Table, err = ReadCSV(r)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	//Now convert the data to md
	return j.toJSON()
}

func (j *JSON) TmogCSVTable() ([]byte, error) {
	csvF, err := os.Open(j.Source.Path)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer csvF.Close()
	return j.TmogCSVTableReader(csvF)
}

// Marshal Table to json, store in B
func (j *JSON) toJSON() ([]byte, error) {
	var rows [][]string
	for _, row := range j.Table {
		rows = append(rows, row)
	}
	var err error
	j.B, err = json.Marshal(rows)
	return j.B, err
}
