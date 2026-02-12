package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)


const (
	Day   GetEventsParamsPeriod = "day"
	Month GetEventsParamsPeriod = "month"
	Week  GetEventsParamsPeriod = "week"
)


type CreateEventRequest struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
}


type Event struct {
	DateTime time.Time `json:"dateTime"`
	Id       int64     `json:"id"`
	Title    string    `json:"title"`
}


type UpdateEventRequest struct {
	DateTime time.Time `json:"dateTime"`
	Title    string    `json:"title"`
}


type GetEventsParams struct {

	Date openapi_types.Date `form:"date" json:"date"`

	Period *GetEventsParamsPeriod `form:"period,omitempty" json:"period,omitempty"`
}


type GetEventsParamsPeriod string


type PostEventsJSONRequestBody = CreateEventRequest


type PutEventsIdJSONRequestBody = UpdateEventRequest


type ServerInterface interface {

	GetEvents(w http.ResponseWriter, r *http.Request, params GetEventsParams)


	PostEvents(w http.ResponseWriter, r *http.Request)

	DeleteEventsID(w http.ResponseWriter, r *http.Request, id int64)

	PutEventsID(w http.ResponseWriter, r *http.Request, id int64)
}



type Unimplemented struct{}


func (_ Unimplemented) GetEvents(w http.ResponseWriter, r *http.Request, params GetEventsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}


func (_ Unimplemented) PostEvents(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}


func (_ Unimplemented) DeleteEventsID(w http.ResponseWriter, r *http.Request, id int64) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (_ Unimplemented) PutEventsID(w http.ResponseWriter, r *http.Request, id int64) {
	w.WriteHeader(http.StatusNotImplemented)
}

type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler


func (siw *ServerInterfaceWrapper) GetEvents(w http.ResponseWriter, r *http.Request) {

	var err error


	var params GetEventsParams


	if paramValue := r.URL.Query().Get("date"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "date"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "date", r.URL.Query(), &params.Date)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "date", Err: err})
		return
	}


	err = runtime.BindQueryParameter("form", true, false, "period", r.URL.Query(), &params.Period)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "period", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetEvents(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) PostEvents(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostEvents(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}


func (siw *ServerInterfaceWrapper) DeleteEventsId(w http.ResponseWriter, r *http.Request) {

	var err error


	var id int64

	err = runtime.BindStyledParameterWithOptions("simple", "id", chi.URLParam(r, "id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteEventsID(w, r, id)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func (siw *ServerInterfaceWrapper) PutEventsId(w http.ResponseWriter, r *http.Request) {

	var err error


	var id int64

	err = runtime.BindStyledParameterWithOptions("simple", "id", chi.URLParam(r, "id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PutEventsID(w, r, id)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}


func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}


func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}


func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/events", wrapper.GetEvents)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/events", wrapper.PostEvents)
	})
	r.Group(func(r chi.Router) {
		r.Delete(options.BaseURL+"/events/{id}", wrapper.DeleteEventsId)
	})
	r.Group(func(r chi.Router) {
		r.Put(options.BaseURL+"/events/{id}", wrapper.PutEventsId)
	})

	return r
}