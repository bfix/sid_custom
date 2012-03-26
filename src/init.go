/*
 * Custom cover server start-up code.
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
	"sid"
)

///////////////////////////////////////////////////////////////////////
// Main application start-up code

func main() {

	// set custom configuration handler
	sid.CustomConfigHandler = CustomConfig

	// set custom initialization handler
	sid.CustomInitialization = func() *sid.Cover {

		// show custom configuration settings
		ShowCustomConfig()

		// start manager for image-based cover content
		InitImageHandler()

		// fire off HTTPS server
		go httpsServe()

		// return new custom cover instance
		return NewCover()
	}

	// start framework
	sid.Startup()
}
