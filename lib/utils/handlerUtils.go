package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pulsone21/threattrack/lib/entities"
)

func CreateApiHandleFunc(apiF entities.APIFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "uri", r.RequestURI)
		fmt.Println(r.URL.Path)
		ctx = context.WithValue(ctx, "entity", strings.Split(r.URL.Path, "/")[1])
		fmt.Println("----------New API Request-------------")
		res, err := apiF(ctx, w, r)
		if err != nil {
			Respond(w, err.StatusCode, map[string]string{"RequestUrl": err.RequestUrl, "Message": err.Error()})
			fmt.Println("----------API Request finished with error-------------")
			return
		}
		if res != nil {
			Respond(w, res.StatusCode, res)
			fmt.Println("----------API Request finished with no result-------------")
			return
		}
		if res.StatusCode == 418 {
			fmt.Println("----------API Request finished with hidden response-------------")
			w.Write(nil)

		}
		Respond(w, http.StatusInternalServerError, map[string]string{"RequestUrl": r.RequestURI, "Message": "error and response are nil"})
		fmt.Println("----------API Request finished successfully-------------")
	}
}

func CreateHtmlHandleFunc(fn entities.RequestHtmlFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("----------New API Request-------------")
		ctx := r.Context()
		ctx = context.WithValue(ctx, "uri", r.RequestURI)
		fmt.Println(r.URL)
		ctx = context.WithValue(ctx, "entity", strings.Split(r.URL.Path, "/")[1])
		if err := fn(ctx, w, r); err != nil {
			Respond(w, err.StatusCode, map[string]string{"RequestUrl": err.RequestUrl, "Message": err.Error()})
			fmt.Println(err.Error())
			fmt.Println("----------API Request finished with error-------------")
			return
		}
		fmt.Println("----------API Request finished successfully-------------")
	}
}

func Respond(w http.ResponseWriter, status int, val any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(val)
}

func TestFunction() string {
	return "This Function Call comes from lib/utils package"
}
