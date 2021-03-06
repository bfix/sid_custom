/*
 * Handle HTTPS session for non-TOR document uploads and custom
 * downloads.
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
	"net/http"
	"sid"
	"strconv"
	"strings"
)

///////////////////////////////////////////////////////////////////////
/*
 * Handle HTTPS requests.
 * @param resp http.ResponseWriter - response buffer
 * @param req *http.Request - request data
 */
func handler(resp http.ResponseWriter, req *http.Request) {

	// get requested resource reference
	ref := req.URL.String()
	mth := req.Method
	logger.Println(logger.DBG_ALL, "[https] "+mth+" "+ref)

	//-----------------------------------------------------------------
	// check for POST request (document upload)
	//-----------------------------------------------------------------
	if mth == "POST" {
		// get upload data
		rdr, _, err := req.FormFile("file")
		if err != nil {
			logger.Println(logger.INFO, "[https] Error accessing uploaded file: "+err.Error())
			// show error page
			ref = "/upload_err.html"
			mth = "GET"
		} else {
			content := make([]byte, 0)
			if err = sid.ProcessStream(rdr, 4096, func(data []byte) bool {
				content = append(content, data...)
				return true
			}); err != nil {
				logger.Println(logger.INFO, "[https] Error accessing uploaded file: "+err.Error())
				// show error page
				ref = "/upload_err.html"
				mth = "GET"
			} else {
				// post-process uploaded document
				sid.PostprocessUploadData(content)
				// set resource ref to response page
				ref = "/upload_resp.html"
				mth = "GET"
			}
		}
	}

	//-----------------------------------------------------------------
	// handle GET requests
	//-----------------------------------------------------------------
	if mth == "GET" {
		// set default page
		if ref == "/" {
			ref = "/index.html"
		}
		// handle resource file
		switch {
		case strings.HasSuffix(ref, ".html"):
			resp.Header().Set("Content-Type", "text/html")
		}
		if err := sid.ProcessFile("./www"+ref, 4096, func(data []byte) bool {
			// append data to response buffer
			resp.Write(data)
			return true
		}); err != nil {
			logger.Println(logger.ERROR, "[https] Resource failure: "+err.Error())
		}
	}
}

///////////////////////////////////////////////////////////////////////
/*
 * Start-up the HTTPS server instance.
 */
func httpsServe() {

	// check for disabled HTTPS server
	if Cfg.HttpsPort < 0 {
		logger.Println(logger.INFO, "[https] HTTPS server disabled.")
		return
	}

	// define handlers
	http.HandleFunc("/", handler)

	// start server
	addr := ":" + strconv.Itoa(Cfg.HttpsPort)
	logger.Println(logger.INFO, "[https] Starting server on "+addr)
	if err := http.ListenAndServeTLS(addr, Cfg.HttpsCert, Cfg.HttpsKey, nil); err != nil {
		logger.Println(logger.ERROR, "[https] "+err.Error())
	}
}
