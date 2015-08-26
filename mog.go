// Copyright © 2014, All rights reserved
// Joel Scoble, https://github.com/mohae/transmogrifier
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
package transmogrifier

import (
	"errors"
	"path"
	"path/filepath"
)

/*
type mogger interface{} {

}
*/
const (
	UnsupportedFormat FormatType = iota
	CSV
	MD
	MDTable
)

const(
	UnsupportedResource ResourceType  = iota
	File
)

// FormatType is the format of a resource
type FormatType int

func (f FormatType) String() string { return formatTypes[f] }

var formatTypes = [...]string{
	"unsupported",
	"csv",
	"md",
	"mdtable"
}

func FormatTypeFromString(s string) FormatType {
	s = strings.ToLower(s)
	switch s {
	case "csv":
		return CSV
	case "md":
		return MD
	case "mdtable":
		return MDTable
	}
	return UnsupportedFormat
}

// ResourceType is the type of a resource
type ResourceType int

func (r ResourceType) String() string { return resourceTypes[r] }

var resourceTypes = [...]string{
	"unsupported",
	"file"
}

// ResourceTypeFromString returns the ResourceType constant
func ResourceTypeFromString(s string) ResourceType {
	s = strings.ToLower(s)
	switch s {
	case "file":
		return File
	}
	return UnsupportedResource
}

// Common errors
var (
	ErrNoSource = errors.New("no source was specified")
)

// Currently only supporting local file.
// TODO enable uri support
type resource struct {
	Name string         // Name of the resource
	Path string         // Path of the resource
	Host   string       // Host of the resource
	Format string       // Format of the resource
	Type   ResourceType // Type of the resource
}

func NewResource(s string, t ResourceType) resource {
	if s == "" {
		return resource{Type: t}
	}
	dir := path.Dir(s)
	// if the path didn't contain a directory, make dir an empty string
	if dir == "." {
		dir = ""
	}
	return resource{Name: path.Base(s), Path: dir, Type: t}
}

func (r resource) String() string {
	if r.Path == "" {
		return r.Name
	}
	return filepath.Join(r.Path, r.Name)
}

func (r *resource) SetName(s string) {
	r.Name = s
}

func (r *resource) SetPath(s string) {
	r.Path = s
}
