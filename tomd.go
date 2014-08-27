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
func CSVtoTable(source *io.Reader, destination *io.Writer) {

}

// ReadCSV takes a reader, and reads the data connected with it as CSV data.
// A slice of slice of type string, or an error, are returned. This reads the
// entire file, so if the file is very large and you don't have sufficent RAM
// you will not like the results. There may be a row oriented implementation 
// in the future.
func ReadCSV(r io.Reader )  ([][]string,  error) {
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


