package qishutaLib

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"strings"
)

type ResponseInterface interface {
	Resp() *resty.Response
	Json() gjson.Result
	String() string
	Document() *goquery.Document
	OK() bool
}

type Response struct {
	*resty.Response
}

func (r *Response) Resp() *resty.Response {
	return r.Response
}
func (r *Response) Json() gjson.Result {
	return gjson.Parse(r.String())
}
func (r *Response) Bytes() []byte {
	return r.Response.Body()
}
func (r *Response) OK() bool {
	return r.Response.StatusCode() == 200
}
func (r *Response) Document() *goquery.Document {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(r.String()))
	return doc
}
