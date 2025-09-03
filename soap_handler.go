package soap_server

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func SoapHandler(w http.ResponseWriter, r *http.Request, operations map[string]func(body string) interface{}) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Set response headers
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Header().Set("SOAPAction", "")

	// Read the full request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		sendSOAPFault(w, "Client", "Cannot read request body")
		return
	}

	// Parse the SOAP envelope - try multiple formats
	var envelope SOAPEnvelope

	// First try with soap: prefix
	// soapPrefixedBody := strings.Replace(string(bodyBytes), "<Envelope", "<soap:Envelope", 1)
	// soapPrefixedBody = strings.Replace(soapPrefixedBody, "</Envelope>", "</soap:Envelope>", 1)
	// soapPrefixedBody = strings.Replace(soapPrefixedBody, "<Body", "<soap:Body", 1)
	// soapPrefixedBody = strings.Replace(soapPrefixedBody, "</Body>", "</soap:Body>", 1)

	// Try parsing the original format first
	err = xml.Unmarshal(bodyBytes, &envelope)
	if err != nil {
		sendSOAPFault(w, "Client", "Invalid SOAP request format")
		return
	}
	// Extract operation from body content
	bodyContent := strings.TrimSpace(envelope.Body.Content)
	var response interface{}
	var foundMethod func(args string) interface{}
	var foundMethodName string

	for methodName, operation := range operations {
		if strings.Contains(bodyContent, methodName) {
			foundMethod = operation
			foundMethodName = methodName
			break
		}
	}
	if foundMethod != nil {
		fmt.Println("Method found:", foundMethodName)
		response = foundMethod(bodyContent)
	} else {
		sendSOAPFault(w, "Client", "Unknown operation")
		return
	}
	sendSOAPResponse(w, response)
}

func sendSOAPResponse(w http.ResponseWriter, response interface{}) {
	// Create response envelope
	responseXML, err := xml.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		sendSOAPFault(w, "Server", "Internal server error")
		return
	}

	// Create full SOAP response
	soapResponse := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
    <soap:Body>
        %s
    </soap:Body>
</soap:Envelope>`, string(responseXML))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(soapResponse))

	log.Printf("Sent response: %s", soapResponse)
}

func sendSOAPFault(w http.ResponseWriter, code, message string) {
	soapFault := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
    <soap:Body>
        <soap:Fault>
            <faultcode>%s</faultcode>
            <faultstring>%s</faultstring>
        </soap:Fault>
    </soap:Body>
</soap:Envelope>`, code, message)

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(soapFault))

	log.Printf("Sent fault: %s - %s", code, message)
}
