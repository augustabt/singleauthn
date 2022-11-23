package helpers

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

func ParseDomain() (string, string) {
	var origin, rpid string

	domain := os.Getenv("DOMAIN")
	if len(domain) > 1 {
		parsedDomain, err := url.Parse(domain)
		if err != nil {
			log.Fatal("Required environment variable DOMAIN invalid")
		}

		splitDomain := strings.Split(parsedDomain.Hostname(), ".")
		rpid = fmt.Sprintf("%s.%s", splitDomain[len(splitDomain)-2], splitDomain[len(splitDomain)-1])
		origin = fmt.Sprintf("%s://%s", parsedDomain.Scheme, parsedDomain.Hostname())
	} else {
		rpid = "localhost"
		origin = "http://localhost:7633"
	}

	return origin, rpid
}
