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
package mog

/*
type mogger interface{} {

}
*/

// Currently only supporting local file.
// TODO enable uri support
type Resource struct {
	// Name of the resource
	Name string
	Path string
	//	Scheme string
	Host   string
	Format string
	Type   string
}

func NewResource(path string) Resource {
	if path == "" {
		return Resource{}
	}
	return Resource{Name: path, Path: path}
}
