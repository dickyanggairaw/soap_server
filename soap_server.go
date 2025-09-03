package soap_server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Listen(r *mux.Router, url string, operations map[string]func(body string) interface{}, wsdl string) {
	r.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" && req.URL.RawQuery == "wsdl" {
			WsdlHandler(w, req, wsdl)
		} else {
			SoapHandler(w, req, operations)
		}
	})
}
