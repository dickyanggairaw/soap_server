package soap_server

import "net/http"

func WsdlHandler(w http.ResponseWriter, r *http.Request, wsdl string) {
	wsdlContent := wsdl
	w.Header().Set("Content-Type", "text/xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(wsdlContent))
}
