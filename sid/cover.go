/*
 * Custom cover server implementation: Example for a simple image
 * cover server using the public picture post "imgon.net".
 *
 * (c) 2012 Bernd Fix   >Y<
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or (at
 * your option) any later version.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package sid

///////////////////////////////////////////////////////////////////////
// Public types 

type CustomCover struct {
	posts	map[string]([]byte)		// list of cover POST replacements
}

///////////////////////////////////////////////////////////////////////
// Public functions

func NewCocer() *CustomCover {
	return &CustomCover {
		posts: make (map[string]([]byte)),
	}
}

//=====================================================================
/*
 * Get address of cover server.
 * @return string - server address (name:port)
 */
func (self *CustomCover) GetAddress() string {
	return "imgon.net:80"
}

//---------------------------------------------------------------------
func (self *CustomCover) GetHtmls() map[string]string {
	htmls := make (map[string]string)
	htmls["/"] = "[UPLOAD]"
	return htmls
}

//---------------------------------------------------------------------
/*
 * get cover site POST content for given boundary id.
 * @param id string - boundary id (key used to store POST content)
 * @return []byte - POST content
 */
func (self *CustomCover) GetPostContent (id string) []byte {
	if post,ok := self.posts[id]; ok {
		// delete POST from list
		self.posts[id] = nil,false
		return post
	}
	return nil
}

//---------------------------------------------------------------------
/*
 * Get client-side upload form for next cover content.
 * @param delim striong - boundary identifier
 * @param name string - image name
 * @param mime string - image MIME type
 * @param cmt stribng - image comment
 * @param data []byte - image data (base64-encoded)
 * @return string - action name
 * @return int - size of cover content
 *
 * =====================================
 * POST request format for cover server:
 * =====================================
 *
 *-----------------------------<boundary>
 *Content-Disposition: form-data; name="imgUrl"
 *<nl>
 *<nl>
 *-----------------------------<boundary>
 *Content-Disposition: form-data; name="fileName[]"
 *<nl>
 *<nl>
 *-----------------------------<boundary>
 *Content-Disposition: form-data; name="file[]"; filename="<name>"
 *Content-Type: <mime>
 *<nl>
 *<content>
 *-----------------------------<boundary>
 *Content-Disposition: form-data; name="alt[]"
 *<nl>
 *<description>
 *-----------------------------<boundary>
 *Content-Disposition: form-data; name="new_width[]"
 *<nl>
 *<nl>
 *-----------------------------<boundary>
 *Content-Disposition: form-data; name="new_height[]"
 *<nl>
 *<nl>
 *-----------------------------<boundary>
 *Content-Disposition: form-data; name="submit"
 *<nl>
 *Upload
 *-----------------------------<boundary>--
 *<nl>
 */
func (self *CustomCover) GetUploadForm (delim, name, mime, cmt string, data []byte) (string, int) {

	// build POST content suitable for upload to cover site
	// and save it in the handler structure
	lb := "\r\n"
	lb2 := lb + lb
	lb3 := lb2 + lb
	sep := "-----------------------------" + delim
	post :=
		sep + lb +
		"Content-Disposition: form-data; name=\"imgUrl\"" + lb3 +
		sep + lb +
		"Content-Disposition: form-data; name=\"fileName[]\"" + lb3 +
		sep + lb +
		"Content-Disposition: form-data; name=\"file[]\"; filename=\"" + name + "\"" + lb +
 		"Content-Type: " + mime + lb2 +
 		string(data) + lb +
		sep + lb +
		"Content-Disposition: form-data; name=\"alt[]\"\n\n" +
 		cmt + lb +
		sep + lb +
 		"Content-Disposition: form-data; name=\"new_width[]\"" + lb3 +
		sep + lb +
		"Content-Disposition: form-data; name=\"new_height[]\"" + lb3 +
		sep + lb +
		"Content-Disposition: form-data; name=\"submit\"" + lb2 + "Upload" + lb +
		sep + "--" + lb2
	
	self.posts[delim] = []byte(post)
	return "/upload/" + delim, len(self.posts[delim])+32
} 
