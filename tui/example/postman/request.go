package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Header struct {
	Key   string
	Value string
}

func (h Header) String() string {
	return fmt.Sprintf("%s: %s", h.Key, h.Value)
}

type Request struct {
	Method  string
	URL     string
	Headers []Header
	Body    string
}

type Response struct {
	StatusCode int
	Status     string
	Headers    map[string][]string
	Body       string
	Error      string
}

func SendRequest(req Request) Response {
	client := &http.Client{Timeout: 10 * time.Second}

	var body io.Reader
	if req.Body != "" {
		body = strings.NewReader(req.Body)
	}

	httpReq, err := http.NewRequest(req.Method, req.URL, body)
	if err != nil {
		return Response{Error: fmt.Sprintf("Request creation error: %v", err)}
	}

	for _, h := range req.Headers {
		httpReq.Header.Set(h.Key, h.Value)
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return Response{Error: fmt.Sprintf("Request error: %v", err)}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{Error: fmt.Sprintf("Body read error: %v", err)}
	}

	return Response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    resp.Header,
		Body:       string(respBody),
	}
}
