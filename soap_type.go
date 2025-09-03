package soap_server

import "encoding/xml"

// SOAP Envelope structures - flexible to handle different namespace formats
type SOAPEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Xmlns   string   `xml:"xmlns:soap,attr,omitempty"`
	Body    SOAPBody `xml:"Body"`
}

type SOAPBody struct {
	XMLName xml.Name `xml:"Body"`
	Content string   `xml:",innerxml"`
}

type SOAPFault struct {
	XMLName xml.Name `xml:"soap:Fault"`
	Code    string   `xml:"faultcode"`
	String  string   `xml:"faultstring"`
}
