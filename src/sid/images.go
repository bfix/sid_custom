/*
 * Custom cover strategy: Use images as cover content and public
 * picture posts as cover server.
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
	"encoding/xml"
	"gospel/crypto"
	"gospel/logger"
	"os"
)

///////////////////////////////////////////////////////////////////////
// Image handler: Provide access to (annotated) image content
// to be used as cover content in upload procedures. Annotations
// include mime type, comments, etc. pp.
///////////////////////////////////////////////////////////////////////
/*
 * The XML definition file for image references looks like this:
 *
 *<?xml version="1.0" encoding="UTF-8"?>
 *<images>
 *    <image>
 *        <name>Test</name>
 *        <comment>blubb</comment>
 *        <path>./image/img001.gif</path>
 *        <mime>image/gif</mime>
 *    </image>
 *</images>
 */

//=====================================================================
/*
 * List of images (for XML parsing)
 */
type ImageList struct {
	Image []ImageDef `xml:"image"`
}

//=====================================================================
/*
 * Image definition (XML).
 */
type ImageDef struct {
	//-----------------------------------------------------------------
	// XML mapped fields
	//-----------------------------------------------------------------
	Name    string `xml:"name"`
	Comment string `xml:"comment"`
	Path    string `xml:"path"`
	Mime    string `xml:"mime"`
}

/*
 * Image definition (List).
 */
type ImageRef struct {
	name    string
	comment string
	path    string
	mime    string
	size    int
}

//=====================================================================
/*
 * List of known image references.
 */
var imgList []*ImageRef

//---------------------------------------------------------------------
/*
 * Initialize image handler: read image definitions from the file
 * specified by the "defs" argument.
 * @param defs string - name of XML-based image definitions 
 */
func InitImageHandler() {

	// prepare parsing of image references
	imgList = make([]*ImageRef, 0)
	rdr, err := os.Open(Cfg.ImageDefs)
	if err != nil {
		// terminate application in case of failure
		logger.Println(logger.ERROR, "[images] Can't read image definitions -- terminating!")
		logger.Println(logger.ERROR, "[images] "+err.Error())
		os.Exit(1)
	}
	defer rdr.Close()

	// parse XML file and build image reference list
	decoder := xml.NewDecoder(rdr)
	var list ImageList
	if err = decoder.Decode(&list); err != nil {
		// terminate application in case of failure
		logger.Println(logger.ERROR, "[images] Can't decode image definitions -- terminating!")
		logger.Println(logger.ERROR, "[images] "+err.Error())
		os.Exit(1)
	}
	for _, img := range list.Image {
		logger.Println(logger.DBG, "[images]: image="+img.Name)
		// get size of image file
		fi, err := os.Stat(img.Path)
		if err != nil {
			logger.Println(logger.ERROR, "[images] image '"+img.Path+"' missing!")
			continue
		}
		// clone to reference instance
		ir := &ImageRef{
			name:    img.Name,
			comment: img.Comment,
			path:    img.Path,
			mime:    img.Mime,
			size:    int(fi.Size()),
		}
		// add to image list
		imgList = append(imgList, ir)
	}
	logger.Printf(logger.INFO, "[images] %d images available\n", len(imgList))
}

//---------------------------------------------------------------------
/*
 * Get next (random) image from repository
 * @return *ImageRef - reference to (random) image
 */
func GetNextImage() *ImageRef {
	return imgList[crypto.RandInt(0, len(imgList)-1)]
}
