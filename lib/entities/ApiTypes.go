package entities

import (
	"context"
	"net/http"
	"reflect"
)

type ApiError struct {
	error
	StatusCode int
	RequestUrl string
}
type ScanFunc func(dest ...any) error

type ApiResponse struct {
	StatusCode int
	RequestUrl string
	Data       interface{}
}

type APIFunction func(context.Context, http.ResponseWriter, *http.Request) (*ApiResponse, *ApiError)

type RequestHtmlFunction func(context.Context, http.ResponseWriter, *http.Request) *ApiError

func InternalServerError(err error, uri string) *ApiError {
	return NewApiError(http.StatusInternalServerError, uri, err)
}

func BadRequestError(err error, uri string) *ApiError {
	return NewApiError(http.StatusBadRequest, uri, err)
}
func NotFoundError(err error, uri string) *ApiError {
	return NewApiError(http.StatusNotFound, uri, err)
}

func NotImplementedError(err error, uri string) *ApiError {
	return NewApiError(http.StatusNotImplemented, uri, err)
}

func NewApiError(status int, uri string, err error) *ApiError {
	return &ApiError{
		error:      err,
		StatusCode: status,
		RequestUrl: uri,
	}
}

func HiddenApiResponse() * ApiResponse{
	return &ApiResponse{
		StatusCode: 418,
		RequestUrl: "",
		Data: nil,
	}
}

func NewApiResponse(statusCode int, uri string, data interface{}) *ApiResponse {
	v := reflect.TypeOf(data).Kind()
	if v != reflect.Array && v != reflect.Slice {
		data = []interface{}{data}
	}
	return &ApiResponse{
		StatusCode: statusCode,
		RequestUrl: uri,
		Data:       data,
	}
}
