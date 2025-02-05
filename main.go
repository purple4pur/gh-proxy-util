package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

var HOSTNAME = os.Getenv("GH_PROXY_HOSTNAME")

func modifyResponse(r *http.Response) error {
	//log.Printf("[response] Content-Encoding: %s\n", r.Header.Get("Content-Encoding"))
	r.Header.Del("Content-Security-Policy")
	if l := r.Header.Get("Location"); l != "" {
		l = strings.Replace(l, "https://raw.githubusercontent.com", "https://raw-githubusercontent." + HOSTNAME, 1)
		l = strings.Replace(l, "https://objects.githubusercontent.com", "https://objects-githubusercontent." + HOSTNAME, 1)
		r.Header.Set("Location", l)
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.Header.Get("Content-Encoding") == "gzip" {
		r.Header.Del("Content-Encoding")
		zr, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return err
		}
		defer zr.Close()
		buf, err := io.ReadAll(zr)
		if err != nil {
			return err
		}
		b = buf
	}

	b = bytes.ReplaceAll(b, []byte("=\"https://github.com"), []byte("=\"https://ghproxy." + HOSTNAME))
	b = bytes.ReplaceAll(b, []byte("=\"https://github.githubassets.com"), []byte("=\"https://github-githubassets." + HOSTNAME))
	b = bytes.ReplaceAll(b, []byte("=\"https://avatars.githubusercontent.com"), []byte("=\"https://avatars-githubusercontent." + HOSTNAME))

	r.Body = io.NopCloser(bytes.NewReader(b))
	r.ContentLength = int64(len(b))
	r.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("error: no listen port or target URL.\nusage: %s LISTEN_PORT->TARGET_URL [LISTEN_PORT->TARGET_URL ...]\n", os.Args[0])
	}

	var wg sync.WaitGroup

	for _, arg := range os.Args[1:] {
		cfg := strings.Split(arg, "->")
		if len(cfg) != 2 {
			log.Fatalf("error: bad format of LISTEN_PORT->TARGET_URL.\nusage: %s LISTEN_PORT->TARGET_URL [LISTEN_PORT->TARGET_URL ...]\n", os.Args[0])
		}
		listenPort := ":" + cfg[0]
		targetUrl := cfg[1]

		wg.Add(1)

		go func(listenPort string, targetUrl string) {
			defer wg.Done()
			u, _ := url.Parse(targetUrl)
			log.Printf("[main] Forwarding %s -> %s\n", listenPort, targetUrl)

			proxy := httputil.NewSingleHostReverseProxy(u)
			proxy.ModifyResponse = modifyResponse

			log.Fatal(http.ListenAndServe(listenPort, proxy))
		}(listenPort, targetUrl)
	}

	wg.Wait()
}
