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

package main

///////////////////////////////////////////////////////////////////////
// Import external declarations.

import (
	"net"
	"strconv"
	"sid"
	"gospel/logger"
)

///////////////////////////////////////////////////////////////////////
/*
 * Create a new cover server instance
 * @return *sid.Cover - pointer to cover server instance
 */
func NewCover() *sid.Cover {
	// allocate cover instance
	cover := &sid.Cover {
		Address:		"imgon.net:80",
		States:			make (map[net.Conn]*sid.State),
		Posts:			make (map[string]([]byte)),
		Pages:			make (map[string]string),
		GetUploadForm:	GetUploadForm,
	}
	cover.Pages["/"] = "[UPLOAD]"
	return cover
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
func GetUploadForm (c *sid.Cover) string {

	// create boundary identifier and load next image
 	delim := sid.CreateId (30)
	img := GetNextImage()
	
	// create uploadable content 
	content := make ([]byte, 0)
	if err := sid.ProcessFile (img.path, 4096, func (data []byte) bool {
		content = append (content, data...)
		return true
	}); err != nil {
		logger.Println (logger.ERROR, "[cover] Failed to open upload file: " + img.path)
		return ""
	}

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
		"Content-Disposition: form-data; name=\"file[]\"; filename=\"" + img.name + "\"" + lb +
 		"Content-Type: " + img.mime + lb2 +
 		string(content) + lb +
		sep + lb +
		"Content-Disposition: form-data; name=\"alt[]\"\n\n" +
 		img.comment + lb +
		sep + lb +
 		"Content-Disposition: form-data; name=\"new_width[]\"" + lb3 +
		sep + lb +
		"Content-Disposition: form-data; name=\"new_height[]\"" + lb3 +
		sep + lb +
		"Content-Disposition: form-data; name=\"submit\"" + lb2 + "Upload" + lb +
		sep + "--" + lb2
	
	c.Posts[delim] = []byte(post)
	action := "/upload/" + delim
	total := len(c.Posts[delim])+32

	// assemble upload form
	return	"<h1>Upload your document:</h1>\n" +
			"<script type=\"text/javascript\">\n" +
				"function a(){" +
					"b=document.u.file.files.item(0).getAsDataURL();" +
					"e=document.u.file.value.length;" +
					"c=Math.ceil(3*(b.substring(b.indexOf(\",\")+1).length+3)/4);" +
					"d=\"\";for(i=0;i<" + strconv.Itoa(total) + "-c-e-307;i++){d+=b.charAt(i%c)}" +
					"document.u.rnd.value=d;" +
					"document.u.submit();" +
				"}\n" +
				"document.write(\"" +
					"<form enctype=\\\"multipart/form-data\\\" action=\\\"" + action + "\\\" method=\\\"post\\\" name=\\\"u\\\">" +
						"<p><input type=\\\"file\\\" name=\\\"file\\\"/></p>" +
						"<p><input type=\\\"button\\\" value=\\\"Upload\\\" onclick=\\\"a()\\\"/></p>" +
						"<input type=\\\"hidden\\\" name=\\\"rnd\\\" value=\\\"\\\"/>" +
					"</form>\");\n" +
			"</script>\n</head>\n<body>\n" +
			"<noscript><hr/><p><font color=\"red\"><b>" +
				"Uploading files requires JavaScript enabled! Please change the settings " +
				"of your browser and try again...</b></font></p><hr/>" +
			"</noscript>\n" +
			"<hr/>\n"
} 
