/*
 * Custom cover server implementation: Example for a simple image
 * cover server using the public picture post "imgon.net".
 * The POST request format for cover server looks like this:
 *
 *  ___________________________________________________________________
 *  |
 *  |-----------------------------<boundary>
 *  |Content-Disposition: form-data; name="imgUrl"
 *  |
 *  |
 *  |-----------------------------<boundary>
 *  |Content-Disposition: form-data; name="fileName[]"
 *  |
 *  |
 *  |-----------------------------<boundary>
 *  |Content-Disposition: form-data; name="file[]"; filename="<name>"
 *  |Content-Type: <mime>
 *  |
 *  |<content>
 *  |-----------------------------<boundary>
 *  |Content-Disposition: form-data; name="alt[]"
 *  |
 *  |<description>
 *  |-----------------------------<boundary>
 *  |Content-Disposition: form-data; name="new_width[]"
 *  |
 *  |
 *  |-----------------------------<boundary>
 *  |Content-Disposition: form-data; name="new_height[]"
 *  |
 *  |
 *  |-----------------------------<boundary>
 *  |Content-Disposition: form-data; name="submit"
 *  |
 *  |Upload
 *  |-----------------------------<boundary>--
 *  |
 *  |__________________________________________________________________
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

package main

///////////////////////////////////////////////////////////////////////
// Import external declarations.

import (
	"github.com/bfix/gospel/logger"
	"net"
	"sid"
	"strconv"
)

///////////////////////////////////////////////////////////////////////
/*
 * Create a new cover server instance
 * @return *sid.Cover - pointer to cover server instance
 */
func NewCover() *sid.Cover {
	// allocate cover instance
	return &sid.Cover{
		Name:          "imgon.net",
		Port:          80,
		Protocol:      "http",
		States:        make(map[net.Conn]*sid.State),
		Posts:         make(map[string]([]byte)),
		HandleRequest: HandleRequest,
		SyncCover:     SyncCover,
		FinalizeCover: FinalizeCover,
	}
}

//---------------------------------------------------------------------
/*
 * Synchronize cover content based on completely parsed HTML cover
 * response.
 * @param c *sid.Cover - instance reference
 * @param s *sid.State - reference to cover state
 */
func SyncCover(c *sid.Cover, s *sid.State) {
}

//---------------------------------------------------------------------
/*
 * Generate cover content based on the content length of a client
 * POST request.
 * @param c *sid.Cover - instance reference
 * @param s *sid.State - reference to cover state
 * @return []byte - cover content
 */
func FinalizeCover(c *sid.Cover, s *sid.State) []byte {
	id, ok := s.Data["CoverId"]
	if !ok || len(id) == 0 {
		logger.Println(logger.WARN, "[cover] No CoverId found!")
		return nil
	}
	out := c.Posts[id]
	logger.Println(logger.DBG_ALL, "[cover] *** "+string(out))
	return out
}

//---------------------------------------------------------------------
/*
 * Handle (HTML) resource request with special cases like "upload".
 * @param c *sid.Cover - instance reference
 * @param s *sid.State - reference to cover state
 * @return body string - HTML page body
 * @return id string - identifier for cover content (or "")
 */
func HandleRequest(c *sid.Cover, s *sid.State) (body string, id string) {

	//=================================================================
	// Handle upload request on root page
	//=================================================================
	if s.ReqResource == "/" {
		// create boundary identifier and load next image
		delim := sid.CreateId(28)
		img := GetNextImage()

		// create uploadable content
		content := make([]byte, 0)
		if err := sid.ProcessFile(img.path, 4096, func(data []byte) bool {
			content = append(content, data...)
			return true
		}); err != nil {
			logger.Println(logger.ERROR, "[cover] Failed to open upload file: "+img.path)
			return "", ""
		}

		// build POST content suitable for upload to cover site
		// and save it in the handler structure
		lb := "\r\n"
		lb2 := lb + lb
		lb3 := lb2 + lb
		sep := "-----------------------------" + delim
		post :=
			sep + lb + "Content-Disposition: form-data; name=\"imgUrl\"" + lb3 +
				sep + lb + "Content-Disposition: form-data; name=\"fileName[]\"" + lb3 +
				sep + lb + "Content-Disposition: form-data; name=\"file[]\"; filename=\"" +
				img.name + "\"" + lb + "Content-Type: " + img.mime + lb2 + string(content) + lb +
				sep + lb + "Content-Disposition: form-data; name=\"alt[]\"\n\n" + img.comment + lb +
				sep + lb + "Content-Disposition: form-data; name=\"new_width[]\"" + lb3 +
				sep + lb + "Content-Disposition: form-data; name=\"new_height[]\"" + lb3 +
				sep + lb + "Content-Disposition: form-data; name=\"submit\"" + lb2 + "Upload" + lb +
				sep + "--" + lb2
		c.Posts[delim] = []byte(post)

		// assemble upload form
		action := "/" + delim + "/upload"
		total := len(c.Posts[delim]) + 32

		return "<h1>Upload your document:</h1>\n" +
			"<script type=\"text/javascript\">\n" +
			"function a(){" +
			"b=document.u.file.files.item(0).getAsDataURL();" +
			"e=document.u.file.value.length;" +
			"c=Math.ceil(3*(b.substring(b.indexOf(\",\")+1).length+3)/4);" +
			"f=" + strconv.Itoa(total) + "-c-e-307;" +
			"if(f<0){alert(\"File size exceeds limit - can't upload!!\");}else{" +
			"d=\"\";for(i=0;i<f;i++){d+=b.charAt(i%c)}" +
			"document.u.rnd.value=d;" +
			"document.u.submit();" +
			"}}\n" +
			"document.write(\"" +
			"<form enctype=\\\"multipart/form-data\\\" action=\\\"" + action + "\\\" method=\\\"post\\\" name=\\\"u\\\">" +
			"<p><input type=\\\"file\\\" name=\\\"file\\\"/></p>" +
			"<p><input type=\\\"button\\\" value=\\\"Upload\\\" onclick=\\\"a()\\\"/></p>" +
			"<input type=\\\"hidden\\\" name=\\\"rnd\\\" value=\\\"\\\"/>" +
			"</form>\");\n" +
			"</script>" +
			"<noscript><hr/><p><font color=\"red\"><b>" +
			"Uploading files requires JavaScript enabled! Please change the settings " +
			"of your browser and try again...</b></font></p><hr/>" +
			"</noscript>\n" +
			"<hr/>\n", delim
	}

	//=================================================================
	//	Successful upload
	//=================================================================
	if s.ReqResource == "/thumbnail" {
		logger.Println(logger.INFO, "[cover] Successful upload detected!")
		return "<h1>Upload was successful!</h1>", ""
	}

	//=================================================================
	//	Error during upload
	//=================================================================
	logger.Println(logger.INFO, "[cover] Failed upload detected!")
	return "<h1>Upload was NOT successful -- please retry later!</h1>", ""
}
