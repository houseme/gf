// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gsocket provides a high performance, easy-to-use and feature-rich websocket client/server.
package gsocket

// WebSocket is the interface for websocket.
type WebSocket interface {
	ReadMessage() (messageType int, p []byte, err error)
	SendMessage(messageType int, data []byte) error
	Close() error
}
