package application

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"microtecture/domain/models"
)

// Context is context of web application
type Context struct {
	Request       *http.Request
	Response      http.ResponseWriter
	RemoteAddress string
	User          *models.User
}

// NewContext creates and returns Context
func NewContext() *Context {
	return &Context{}
}

// WithLogger add logger to context instance
//func (self *Context) WithLogger(logger logrus.FieldLogger) *Context {
//	ret := self
//	ret.Logger = logger
//	return ret
//}

// WithUser add user to context instance
func (self *Context) WithUser(user *models.User) *Context {
	ret := self
	ret.User = user
	return ret
}

// WithRequest add http request instance to context instance
func (self *Context) WithRequest(request *http.Request) *Context {
	ret := self
	ret.Request = request
	return ret
}

// WithResponseWriter add http response wirter instance to context instance
func (self *Context) WithResponseWriter(responseWriter http.ResponseWriter) *Context {
	ret := self
	ret.Response = responseWriter
	return ret
}

// DecodeMoel decodes model from context request to an interface domain model
func (self *Context) DecodeModel(v interface{}) error {
	if err := json.NewDecoder(self.Request.Body).Decode(v); err != nil {
		return NewErrValidation(err.Error())
	}
	return nil
}

// Finish writes data with json format and header to http response
func (self *Context) Finish(status int, v interface{}) error {
	if v != nil {
		if err := self.json(v); err != nil {
			return err
		}
	}

	self.Response.WriteHeader(status)
	return nil
}

func (self *Context) json(v interface{}) error {
	err := json.NewEncoder(self.Response).Encode(v)
	return HandleError(err)
}

// ReadCookie reads cookie from context request
func (self *Context) ReadCookie(cookieName string) (cookieValue string, err error) {
	cookie, err := self.Request.Cookie(cookieName)
	if err != nil {
		return "", errors.New(err.Error())
	}

	return cookie.Value, nil
}

// SetCookie sets cookie to context response
// maxAge is second count
func (self *Context) SetCookie(name, value string, maxAge uint) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  value,
		MaxAge: int(maxAge),
	}

	http.SetCookie(self.Response, cookie)
}
