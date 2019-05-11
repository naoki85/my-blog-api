package controllers

type ResponseWriter interface {
	Header() Header
	WriteHeader(int)
	Write([]byte)
}

type Request interface {
}

type Params interface {
	ByName(string) string
}

type Header map[string][]string

func (h Header) Set(key, value string) {}
