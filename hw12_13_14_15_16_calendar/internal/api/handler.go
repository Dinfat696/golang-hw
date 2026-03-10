package api

import "net/http"

func Handler(h interface{}) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}
