package network

import (
	"log"

	"io/ioutil"
	"net/http"
)


// BatchGet gets data from urls in parallel.
// Each request tried up to "retries" times.
// Responses are returned on channel as []byte.
// Nil is returned if an error occurs: the error is printed.
func BatchGet(c chan []byte, urls []string, retries int) {
	if retries < 0 { retries = 0 }
	for _, url := range urls {
		go func(url string) {
			try := 0
			// 0 retries means try once
			for try < retries + 1 {

				resp, err := http.Get(url)
				if err != nil { 
					// final try
					if try == retries {
						log.Printf("Error in get from url: %s\n", url)
						c <- nil
						return
					} else {
						try++
						log.Printf("Retry get #%d from url: %s\n", try, url)
						continue
					}
				}
				defer resp.Body.Close()

				var body []byte
				body, err = ioutil.ReadAll(resp.Body)
				if err != nil { 
					// final try
					if try == retries {
						log.Printf("Error reading response from url: %s\n", url)
						c <- nil
						return
					} else {
						try++
						log.Printf("Retry reading response #%d from url: %s\n", try, url)
						continue
					}
				}

				c <- body
				return
			}
		}(url)
	}
}