// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsocket

import (
	"net/http"
)

// WebSocketClient wraps the underlying websocket client connection
// and provides convenient functions.
type WebSocketClient interface {
	// Dial connects to the url with the http.Header.
	// It's the same as function Connect.
	Dial(url string, requestHeader http.Header) error
	WebSocket
}
