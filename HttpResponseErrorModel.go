package uhttp

import (
	"encoding/json"
	"errors"
	"io"
)

func NewHttpErrorResponse(err error) HttpResponseErrorModel {
	return HttpResponseErrorModel{
		Error: err.Error(),
	}
}

type HttpResponseErrorModel struct {
	Error string `json:"error"`
}

func ErrorFromHttpResponseBody(r io.ReadCloser) (error, error) {
	m := HttpResponseErrorModel{}
	err := json.NewDecoder(r).Decode(&m)
	if err != nil {
		return nil, err
	}
	return errors.New(m.Error), nil
}
