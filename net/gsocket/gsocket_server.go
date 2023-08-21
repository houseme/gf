// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsocket

import (
	"net/http"
)

// WebSocketServer wraps the underlying websocket server connection
// and provides convenient functions.
type WebSocketServer interface {
	// Connect accepts the http request and upgrades it to websocket connection.
	Connect(r *http.Request, w http.ResponseWriter) (err error)
	WebSocket
}
