package p24

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

// Error reports an error and the Method, URL, response/request body that caused it
type Error struct {
	Err         error
	URL, Method string
	Req, Resp   []byte
}

func (e *Error) Error() (err string) {
	// nolint:gocritic // e.Err can be nil
	return fmt.Sprint(e.Err)
}

func (e *Error) Unwrap() error { return e.Err }

func (e *Error) Cause() error { return e.Err }

func newError(err error, url, method string, resp, req []byte) *Error {
	return &Error{err, url, method, resp, req}
}

// Client performs p24 api calls with given Doer, Merchant, Logger.
// Implements p24 information API client.
// see: https://api.privatbank.ua/#p24/main
type Client struct {
	http     Doer
	log      Logger
	merchant Merchant
}

// ClientOpts is a full set of all parameters to initialize Client
type ClientOpts struct {
	HTTP     Doer
	Log      Logger
	Merchant Merchant
}

// NewClient returns Client instance with given opts
func NewClient(opts ClientOpts) *Client {
	var log Logger = LogFunc(func(format string, args ...interface{}) {})
	if opts.Log != nil {
		log = opts.Log
	}
	return &Client{
		http:     opts.HTTP,
		log:      log,
		merchant: opts.Merchant,
	}
}

// DoContext performs a p24 http api call with given url, method, request
// and unmarshal response body to resp if no errors occurred
// nolint:gocyclo // Is a complexity function
func (c *Client) DoContext(ctx context.Context, url, method string, req Request, resp *Response) error {
	// process http req
	httpReqBody, err := xml.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "can`t marshal req")
	}
	httpReqBody = []byte(xml.Header + string(httpReqBody)) // insert xml header above
	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(httpReqBody))
	if err != nil {
		return errors.Wrap(err, "can`t make http request")
	}
	httpReq.Header.Add("Content-Type", "application/xml; charset=utf-8")

	// process http resp
	httpResp, err := c.http.Do(httpReq)
	if err != nil {
		return newError(errors.Wrap(err, "http request failed"), url, method, httpReqBody, nil)
	}
	defer func() {
		if err = httpResp.Body.Close(); err != nil {
			c.log.Logf("[WARN] failed to close http response body: %+v\n", err)
		}
	}()
	httpRespBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return newError(errors.Wrap(err, "can`t read http response body"), url, method, httpReqBody, nil)
	}
	if httpResp.StatusCode >= 300 {
		return newError(errors.Errorf("unexpected http status code %d", httpResp.StatusCode), url, method, httpReqBody, httpRespBody)
	}

	// parse xml resp
	xmlResp := xmlResp(httpRespBody)
	if err = xmlResp.CheckErr(); err != nil {
		return newError(errors.Wrap(err, "xml response with error"), url, method, httpReqBody, httpRespBody)
	}
	if err = xmlResp.CheckContent(); err != nil {
		return newError(errors.Wrap(err, "unexpected xml response content"), url, method, httpReqBody, httpRespBody)
	}
	if err = xmlResp.VerifySign(c.merchant); err != nil {
		return newError(errors.New("xml response with invalid signature"), url, method, httpReqBody, httpRespBody)
	}
	if err = xml.Unmarshal(xmlResp, resp); err != nil {
		return newError(errors.Wrap(err, "can`t unmarshal xml response"), url, method, httpReqBody, httpRespBody)
	}

	return nil
}
