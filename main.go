package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	certFile = flag.String("cert", "./cert.pem", "A PEM eoncoded certificate file.")
	keyFile  = flag.String("key", "./private.pem", "A PEM encoded private key file.")
)

func main()  {
	fmt.Println("Starting script...")

	// Load client cert
	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		log.Fatal(err)
	}

	// SSLKEYLOGFILE writer
	w, err := os.OpenFile("tls-secrets.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		KeyLogWriter: w,
	}
	c := colly.NewCollector(
		colly.AllowedDomains("some-site"),
	)
	c.WithTransport(&http.Transport{TLSClientConfig: tlsConfig})
	c.OnHTML("#search-results", func (e *colly.HTMLElement) {
		numChildren := e.DOM.Children().Length()
		fmt.Printf("Found %s child elements... ", strconv.Itoa(numChildren))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting ... ", r.URL.String())
	})

	url := "some-site"
	err = c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}
}
