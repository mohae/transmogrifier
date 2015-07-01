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

// Common errors
var (
	ErrNoSource = errors.New("no source was specified")
)

// Currently only supporting local file.
// TODO enable uri support
type resource struct {
	// Name of the resource
	Name string
	Path string
	//	Scheme string
	Host   string
	Format string
	Type   string
}

func NewResource(s string) resource {
	if s == "" {
		return resource{}
	}
	dir := path.Dir(s)
	// if the path didn't contain a directory, make dir an empty string
	if dir == "." {
		dir = ""
	}
	return resource{Name: path.Base(s), Path: dir}
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
