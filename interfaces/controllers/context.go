package controllers

import "net/http"

type ResponseWriter interface {
	Header() http.Header
	WriteHeader(int)
	Write([]byte) (int, error)
}

type Request interface {
}

type Params interface {
	ByName(string) string
}
