# SOAP Server

A lightweight Go package for creating SOAP web services with minimal configuration. This package provides a simple way to handle SOAP requests and responses, as well as serving WSDL files.

## Features

- Simple API for creating SOAP endpoints
- Automatic handling of SOAP envelope parsing and generation
- Support for serving WSDL files
- Integration with Gorilla Mux router
- Error handling with SOAP fault responses

## Installation

```bash
go get github.com/dickyanggairaw/soap_server
```

## Usage

### Basic Example

```go
package main

import (
	"log"
	"net/http"
	"strings"

	"encoding/xml"

	"github.com/dickyanggairaw/soap_server"
	"github.com/gorilla/mux"
)

// Request/Response structures
type AddRequest struct {
	XMLName xml.Name `xml:"add"`
	A       int      `xml:"a"`
	B       int      `xml:"b"`
}

type AddResponse struct {
	XMLName xml.Name `xml:"addResponse"`
	Result  int      `xml:"result"`
}

func main() {
	r := mux.NewRouter()

	wsdlContent := `<?xml version="1.0" encoding="UTF-8"?>
<definitions name="Calculator"
             targetNamespace="http://calculator.example.com/"
             xmlns="http://schemas.xmlsoap.org/wsdl/"
             xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/"
             xmlns:tns="http://calculator.example.com/"
             xmlns:xsd="http://www.w3.org/2001/XMLSchema"
             xmlns:soap="http://schemas.xmlsoap.org/wsdl/soap/">

  <types>
    <xsd:schema targetNamespace="http://calculator.example.com/"
                elementFormDefault="qualified">

      <xsd:complexType name="addRequest">
        <xsd:sequence>
          <xsd:element name="a" type="xsd:int"/>
          <xsd:element name="b" type="xsd:int"/>
        </xsd:sequence>
      </xsd:complexType>

      <xsd:complexType name="addResponse">
        <xsd:sequence>
          <xsd:element name="result" type="xsd:int"/>
        </xsd:sequence>
      </xsd:complexType>

      <xsd:element name="add" type="tns:addRequest"/>
      <xsd:element name="addResponse" type="tns:addResponse"/>

    </xsd:schema>
  </types>

  <message name="addRequest">
    <part name="parameters" element="tns:add"/>
  </message>
  <message name="addResponse">
    <part name="parameters" element="tns:addResponse"/>
  </message>

  <!-- PortType -->
  <portType name="CalculatorPortType">
    <operation name="add">
      <input message="tns:addRequest"/>
      <output message="tns:addResponse"/>
    </operation>
  </portType>

  <!-- Binding -->
  <binding name="CalculatorBinding" type="tns:CalculatorPortType">
    <soap:binding style="document" transport="http://schemas.xmlsoap.org/soap/http"/>

    <operation name="add">
      <soap:operation soapAction="http://calculator.example.com/add"/>
      <input><soap:body use="literal"/></input>
      <output><soap:body use="literal"/></output>
    </operation>

  </binding>

  <!-- Service -->
  <service name="CalculatorService">
    <port name="CalculatorPort" binding="tns:CalculatorBinding">
      <soap:address location="http://localhost:8080/calculator"/>
    </port>
  </service>

</definitions>`

	operations := map[string]func(body string) interface{}{
		"add": func(body string) interface{} {
			var req AddRequest
			decoder := xml.NewDecoder(strings.NewReader(body))
			err := decoder.Decode(&req)
			if err != nil {
				log.Printf("Error parsing add request: %v", err)
				return AddResponse{}
			}

			result := req.A + req.B
			log.Printf("Add: %d + %d = %d", req.A, req.B, result)

			return AddResponse{Result: result}
		},
	}
	soap_server.Listen(r, "/calculator", operations, wsdlContent)
	log.Println("SOAP server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
```

### Accessing the Service

- SOAP Endpoint: `http://localhost:8080/calculator`
- WSDL: `http://localhost:8080/calculator?wsdl`

## API Reference

### `Listen(r *mux.Router, url string, operations map[string]func(body string) interface{}, wsdl string)`

Registers a SOAP endpoint with the given router.

- `r`: A Gorilla Mux router
- `url`: The URL path for the SOAP endpoint
- `operations`: A map of operation names to handler functions
- `wsdl`: The WSDL content to serve when requested

## Dependencies

- [github.com/gorilla/mux](https://github.com/gorilla/mux): HTTP router and dispatcher

## License

MIT