/*
 * HTTP redirection service to be used as a fallback handler for
 * unhandled HTTP connections to SID.
 *
 * (c) 2012 Andrew Clausen, Bernd Fix
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
	"text/template"
	"bufio"
	"net"
)

///////////////////////////////////////////////////////////////////////
// Local variables and constants

var redirTempl = template.Must(template.New("redir").Parse(redirStr))

const redirStr = `HTTP/1.0 301 Moved Permanently
Location: {{.URL}}
Content-Type: text/html
Content-Length: 0

`

///////////////////////////////////////////////////////////////////////
// Public types

/*
 * Redirection service for HTTP connections.
 */
type Redir struct {
	URL string
}

///////////////////////////////////////////////////////////////////////
// Public functions and methods

/*
 * Instantiate a new redirection service.
 * @param URL string - where to redirect to
 * @return *Redir - redirection service instance reference
 */
func NewRedir(URL string) *Redir {
	return &Redir{URL}
}

//---------------------------------------------------------------------
/*
 * Handle HTTP connection: send HTTP redirect response on any
 * incoming HTTP request.
 * @param client net.Conn - connection to client
 */
func (self *Redir) Process(client net.Conn) {

	// close client connection on function exit
	defer client.Close()

	// read complete client request
	rdr := bufio.NewReader(client)
	pending := ""
	for {
		// read next line
		b, broken, _ := rdr.ReadLine()
		line := string(b)
		if broken {
			pending = line
			continue
		} else if len(pending) > 0 {
			line = pending + line
		}
		if len(line) == 0 {
			break
		}
	}

	// send redirection response	
	redirTempl.Execute(client, self)
}

//---------------------------------------------------------------------
/*
 * Get service name.
 * @return string - name of control service (for logging purposes)
 */
func (self *Redir) GetName() string {
	return "sid_custom.redir"
}

//---------------------------------------------------------------------
/*
 * Check for TCP protocol.
 * @param protocol string - connection protocol
 * @return bool - protcol handled?
 */
func (self *Redir) CanHandle(protocol string) bool {
	return protocol == "tcp"
}

//---------------------------------------------------------------------
/*
 * Check for appropriate remote address.
 * @param add string - remote address
 * @return bool - local address?
 */
func (self *Redir) IsAllowed(remote string) bool {
	return true
}
