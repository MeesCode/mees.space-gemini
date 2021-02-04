package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pitr/gig"
)

func main() {
	// Gig instance
	g := gig.Default()

	g.Static("/", "index.gmi")

	g.Handle("/book", func(c gig.Context) error {
		return c.File("book.gmi")
	})

	g.Handle("/submit", func(c gig.Context) error {

		// get query
		query, err := c.QueryString()
		if err != nil || query == "" {
			return c.NoContent(gig.StatusInput, "Post text")
		}

		// set name default as ip
		name := "anonymous"

		// use name from cert if available
		cert := c.Certificate()
		if cert != nil {
			name = cert.Subject.CommonName
		}

		// If the file doesn't exist, create it, or append to the file
		f, err := os.OpenFile("book.gmi", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		// build the submission
		line := fmt.Sprintf(
			"[%s] <%s> %s\n",
			time.Now().Format("02 Jan 06 15:04 MST"),
			name,
			query,
		)

		// write to file and close again
		if _, err := f.Write([]byte(line)); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}

		return c.NoContent(gig.StatusRedirectTemporary, "/book")
	})

	// Start server on PORT or default port
	g.Run("cert.pem", "key.pem")
}
