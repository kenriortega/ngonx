package interceptor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Can be use this
// c := &http.Client{
// 	Transport: Interceptor{http.DefaultTransport},
//   }

type Interceptor struct {
	core http.RoundTripper
}

func (Interceptor) ModifyRequest(r *http.Request) *http.Request {
	reqBody := MustHumanize(r.Body)

	// this is the request modification
	modReqBody := []byte(fmt.Sprintf(`{"req": %s}`, reqBody))
	ModReqBodyLen := len(modReqBody)

	req := r.Clone(context.Background())
	req.Body = io.NopCloser(bytes.NewReader(modReqBody))
	req.ContentLength = int64(ModReqBodyLen)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", ModReqBodyLen))

	return req
}

func (Interceptor) ModifyResponse(r *http.Response) *http.Response {
	respBody := MustHumanize(r.Body)

	// this is the response modification
	modRespBody := []byte(fmt.Sprintf(`{"resp": %s}`, respBody))
	ModRespBodyLen := len(modRespBody)

	r.Body = io.NopCloser(bytes.NewReader(modRespBody))
	r.ContentLength = int64(ModRespBodyLen)
	r.Header.Set("Content-Length", fmt.Sprintf("%d", ModRespBodyLen))

	return r
}

func (i Interceptor) RoundTrip(r *http.Request) (*http.Response, error) {
	defer func() {
		_ = r.Body.Close()
	}()

	// modify before the request is sent
	newReq := i.ModifyRequest(r)

	// send the request using the DefaultTransport
	resp, _ := i.core.RoundTrip(newReq)
	defer func() {
		_ = resp.Body.Close()
	}()

	// modify after the response is received
	newResp := i.ModifyResponse(resp)

	return newResp, nil
}

func MustHumanize(r io.Reader) string {
	var m map[string]interface{}
	_ = json.NewDecoder(r).Decode(&m)
	b, _ := json.MarshalIndent(m, "", "  ")
	return string(b)
}
