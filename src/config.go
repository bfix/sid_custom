/*
 * Custom cover server configuration: Example for a simple image
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
	"gospel/logger"
	"gospel/parser"
	"sid"
	"strconv"
)

///////////////////////////////////////////////////////////////////////
// Public types

/*
 * Custom configuration data structure: Settings for using images
 * as cover content (cover server is a picture post) and an
 * additional HTTPS server part for low-security up- and downloads.
 */
type CustomCfg struct {
	ImageDefs string // name of cover image definition file
	HttpsPort int    // port for HTTPS sessions
	HttpsCert string // name of HTTPS certificate file
	HttpsKey  string // name of HTTPS key file
}

///////////////////////////////////////////////////////////////////////
// Global/local variables

/*
 * Custom configuration instance:
 */
var Cfg CustomCfg = CustomCfg{
	ImageDefs: "./images/images.xml", // image definition file
	HttpsPort: 443,                   // port for HTTPS connections
	HttpsCert: "./cert.pem",          // HTTPS certificate
	HttpsKey:  "./key.pem",           // HTTPS key
}

///////////////////////////////////////////////////////////////////////
// Public functions and methods

/*
 * Handle custom configuration options.
 * @param mode int - parameter mode
 * @param param *parser.Parameter - key/value setting
 * @return bool - continue configuration processing?
 */
func CustomConfig(mode int, param *parser.Parameter) bool {
	if param != nil {
		if mode != parser.LIST {
			switch param.Name {
			case "HttpsPort":
				sid.SetIntValue(&Cfg.HttpsPort, param.Value)
			case "HttpsCert":
				Cfg.HttpsCert = param.Value
			case "HttpsKey":
				Cfg.HttpsKey = param.Value
			case "ImageDefs":
				Cfg.ImageDefs = param.Value
			}
		}
	}
	return true
}

//---------------------------------------------------------------------
/*
 * Show custom configration settings.
 */
func ShowCustomConfig() {
	// list current configuration data
	logger.Println(logger.INFO, "[config] !==========< custom configuration >===============")
	logger.Println(logger.INFO, "[config] !Image library definition: "+Cfg.ImageDefs)
	logger.Println(logger.INFO, "[config] !Port for HTTPS sessions: "+strconv.Itoa(Cfg.HttpsPort))
	logger.Println(logger.INFO, "[config] !HTTPS certificate: "+Cfg.HttpsCert)
	logger.Println(logger.INFO, "[config] !HTTPS SSL key: "+Cfg.HttpsKey)
	logger.Println(logger.INFO, "[config] !==========================================")
}
