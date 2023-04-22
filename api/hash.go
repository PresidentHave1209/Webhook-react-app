package handler

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.URL.String())

	spliturl := strings.Split(r.URL.String(), "?")
	webhook := strings.Split(spliturl[1], "~")

	fullurl := "https://discordapp.com/api/webhooks/" + webhook[0]
	fmt.Println("URL:>", fullurl)

	message, failure := url.QueryUnescape(webhook[1])
	if failure != nil {
		return
	}
	
	message = strings.Replace(message, "@", "@ ", -1)

	var jsonStr = []byte(`{"content":"[` + hash(r.RemoteAddr) + "] " + message + `"}`)
	req, err := http.NewRequest("POST", fullurl, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		panic(err)
	}

	if r.Header.Get("If-Unmodified-Since") == "" {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	}
}

func hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	rs := fmt.Sprintf("%X", bs)
	return rs[0:6]
}
