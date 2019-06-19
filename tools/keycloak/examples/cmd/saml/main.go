package main

import (
	"fmt"
	"net/http"
	"net/url"

	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"github.com/crewjam/saml/samlsp"
	"log"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%+v\n", samlsp.Token(r.Context()))
	fmt.Fprintf(w, "Hello, %s!", samlsp.Token(r.Context()).Attributes.Get("username"))
}

func main() {
	port := flag.Int("port", 8000, "set port")
	x509key := flag.String("key", "key1.key", "set key")
	x509cert := flag.String("cert", "cert1.cert", "set cert")
	flag.Parse()

	// curl http://localhost/auth/realms/Sample/protocol/saml/descriptor すると XML が返ってくる
	confIDPMetadataURL := "http://localhost/auth/realms/Sample/protocol/saml/descriptor"
	confRootURL := fmt.Sprintf("http://localhost:%d", *port)
	confPort := fmt.Sprintf(":%d", *port)

	keyPair, err := tls.LoadX509KeyPair(*x509cert, *x509key)
	if err != nil {
		log.Fatal(err)
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		log.Fatal(err)
	}

	idpMetadataURL, err := url.Parse(confIDPMetadataURL)
	if err != nil {
		log.Fatal(err)
	}

	rootURL, err := url.Parse(confRootURL)
	if err != nil {
		log.Fatal(err)
	}

	samlSP, _ := samlsp.New(samlsp.Options{
		URL: *rootURL,
		// Key はリクエストの署名に使うもの
		Key:            keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate:    keyPair.Leaf,
		IDPMetadataURL: idpMetadataURL,
	})

	// curl http://localhost:8000/saml/metadata で SP のメタデータがとれる

	app := http.HandlerFunc(hello)
	// http://localhost/auth/realms/Sample/protocol/saml?RelayState=b54Dizebx2PgCm5fg3Jqtk9n_63obxDmv-O_LI832wfFZ7OAS34G2QQ1&SAMLRequest=nFLNbtNAEH4Va%2B7JrjdO06xqS6ERIlKBqAkcuI3tCVlpf8zOGOjbo7itVISUA9fd%2Bf5mvjvG4Ae7GeUcH%2BnHSCzF7%2BAj28tHDWOONiE7thEDsZXOHjYfH6yZa4vMlMWlCG8gw3XMkJOkLnkodtsaXD8rS9OXq9asqmVZmVav2q7Cdt3fLvoFdmZZ6eWpXbcExVfK7FKswcw1FDvmkXaRBaPUYHS5numbWWmORlt9Y6vFvDLfoNgSi4soE%2FAsMlilfOrQnxOLwlHOKhP6wOqAYfCkXg2qSxgo3qfc0bSdGk7omaDYvMa%2BT5HHQPlA%2Bafr6Mvjw78S9lZrPXEp7BiK%2FQv9Oxd7F79fX1b7PMT2w%2FG4n%2B0%2FH47QTAezU%2Fp8cRdQrpNcXlw%2FO02jlqI4eYLmis9Agj0K3qk3Us1LUT5hoN12n7zrnv5DXjJGdhQFio336dd9JhSqQfJIoJpnyb%2Fr2PwJAAD%2F%2Fw%3D%3D
	http.Handle("/hello", samlSP.RequireAccount(app))
	// TODO: 見に行くのはさらに /saml/acs ? とかなので、/saml でなく、/saml/ である必要あり
	http.Handle("/saml/", samlSP)
	http.ListenAndServe(confPort, nil)
}
