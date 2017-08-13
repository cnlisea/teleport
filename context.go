// Copyright 2015-2017 HenryLee. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package teleport

import (
	"net/url"

	"github.com/henrylee2cn/goutil"

	"github.com/henrylee2cn/teleport/socket"
)

// SvrContext server controller context.
// For example:
//  type Home struct{ SvrContext }
type SvrContext interface {
	Uri() string
	Path() string
	Query() url.Values
	Public() goutil.Map
	PublicLen() int
	SetCodec(string)
	Ip() string
	// RealIp() string
}

// SocketCtx a context interface of Conn for one request/outputonse.
type SocketCtx interface {
	socket.Socket
}

type context struct {
	session    *Session
	ctxPublic  goutil.Map
	handler    *Handler
	packet     *socket.Packet
	inputBody  reflect.Value
	inputUri   *url.URL
	inputQuery url.Values
}

var (
	_ SvrContext = new(context)
	_ SocketCtx  = new(context)
)

// newCtx creates a context of socket.Socket for one request/outputonse.
func newCtx(s *Session) *context {
	ctx := new(context)
	ctx.reset(s)
	return ctx
}

func (c *context) reset(s *Session) {
	c.session = s
	c.ctxPublic = goutil.RwMap()
	if s.PublicLen() > 0 {
		s.socket.Public().Range(func(key, value interface{}) bool {
			c.ctxPublic.Store(key, value)
			return true
		})
	}
	ctx.packet = socket.NewPacket(ctx.loadBody)
}

func (c *context) free() {
	c.packet.Reset(nil)
	c.ctxPublic = nil
}

// Public returns temporary public data of Conn Context.
func (c *context) Public() goutil.Map {
	return c.ctxPublic
}

// PublicLen returns the length of public data of Conn Context.
func (c *context) PublicLen() int {
	return c.ctxPublic.Len()
}

func (c *context) Uri() string {
	return c.inputHeader.GetUri()
}

func (c *context) getUri() *url.URL {
	if c.inputUri == nil {
		c.inputUri, _ = url.Parse(c.Uri())
	}
	return c.inputUri
}

func (c *context) Path() string {
	return c.getUri().Path
}

func (c *context) Query() url.Values {
	if c.inputQuery == nil {
		c.inputQuery = c.getUri().Query()
	}
	return c.inputQuery
}

func (c *context) SetCodec(codec string) {
	c.outputHeader.Codec = codec
}

func (c *context) Ip() string {
	return c.Socket.RemoteAddr().String()
}

func (c *context) loadBody(header *socket.Header) interface{} {
	var ok bool
	c.handler, ok = c.session.router.lookup(header)
	if ok {
		c.inputBody = handler.NewArg()
		return c.inputBody.Interface()
	}
	return nil
}
