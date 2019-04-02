package domain

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

// Ticket es un ticket de autenticación
type Ticket struct {
	XMLName xml.Name `xml:"sso"`
	Version string   `xml:"version,attr"`
	ID      struct {
		Src      string `xml:"src,attr"`
		Dst      string `xml:"dst,attr"`
		UniqueID string `xml:"unique_id,attr"`
		GenTime  string `xml:"gen_time,attr"`
		ExpTime  string `xml:"exp_time,attr"`
	} `xml:"id"`
	Operation struct {
		Type  string `xml:"type,attr"`
		Value string `xml:"value,attr"`
		Login struct {
			Service    string `xml:"service,attr"`
			Entity     string `xml:"entity,attr"`
			UID        string `xml:"uid,attr"`
			Authmethod string `xml:"authmethod,attr"`
			Regmethod  string `xml:"regmethod,attr"`
			Info       []struct {
				Name  string `xml:"name,attr"`
				Value string `xml:"value,attr"`
			} `xml:"info"`
		} `xml:"login"`
	} `xml:"operation"`
}

// GetInfo returns the value associated with a given key
func (tkt Ticket) GetInfo(key string) string {
	for _, i := range tkt.Operation.Login.Info {
		if i.Name == key {
			return i.Value
		}
	}
	return ""
}

// ToUser makes a user object from a ticket
func (tkt Ticket) ToUser() User {
	var cuit, _ = strconv.ParseInt(tkt.GetInfo("cuil"), 10, 64)
	var u = User{
		ID:       cuit,
		Email:    fmt.Sprintf("%s@zaffarano.com.ar", tkt.Operation.Login.UID),
		Name:     tkt.GetInfo("cn"),
		Username: tkt.GetInfo("username"),
	}
	return u
}

// ParseTicket parses an XML string as a authentication ticket
func ParseTicket(str string) (Ticket, error) {
	var tkt Ticket
	var err error

	dec := xml.NewDecoder(strings.NewReader(str))
	dec.CharsetReader = charset.NewReaderLabel
	dec.Strict = false

	err = dec.Decode(&tkt)

	return tkt, err
}
