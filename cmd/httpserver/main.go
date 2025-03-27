package main

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"http/internal/headers"
	"http/internal/request"
	"http/internal/response"
	"http/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, HandlerFunc)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func HandlerFunc(w *response.Writer, r *request.Request) {
	if r.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, r)
		return
	}
	if r.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, r)
		return
	}
	if strings.HasPrefix(r.RequestLine.RequestTarget, "/httpbin/") {
		proxyHandler(w, r)
		return
	}
	if r.RequestLine.RequestTarget == "/video" {
		handlerVideo(w, r)
		return
	}
	handler200(w, r)
	return
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(400)
	body := `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody([]byte(body))
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(500)
	body := `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody([]byte(body))
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(200)
	body := `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody([]byte(body))
}

func handlerVideo(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(200)
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting pwd: %s", err)
	}
	fileData, err := os.ReadFile(filepath.Join(wd, "assets/vim.mp4"))
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}
	h := response.GetDefaultHeaders(len(fileData))
	h.Override("Content-Type", "video/mp4")
	w.WriteHeaders(h)
	w.WriteBody(fileData)
}

func proxyHandler(w *response.Writer, r *request.Request) {
	lastPart := strings.TrimPrefix(r.RequestLine.RequestTarget, "/httpbin/")
	proxyURL := fmt.Sprintf("https://httpbin.org/%s", lastPart)
	w.WriteStatusLine(200)
	h := response.GetDefaultHeaders(0)
	h.Remove("Content-Length")
	h.Override("Content-Type", "application/json")
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Trailer", "X-Content-SHA256")
	h.Set("Trailer", "X-Content-Length")
	w.WriteHeaders(h)
	resp, err := http.Get(proxyURL)
	if err != nil {
		log.Fatalf("Error requesting: %s", err)
	}
	bodyBuf := new(bytes.Buffer)
	for {
		buf := make([]byte, 1024)
		n, err := resp.Body.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				bodyHash := sha256.Sum256(bodyBuf.Bytes())
				h = headers.NewHeaders()
				h.Set("X-Content-SHA256", fmt.Sprintf("%x", bodyHash))
				h.Set("X-Content-Length", fmt.Sprintf("%d", len(bodyBuf.Bytes())))
				w.WriteChunkedBodyDone(h)

				log.Println("Written the Whole body")
				return
			}
			log.Fatalf("Error requesting: %s", err)
		}
		bodyBuf.Write(buf[:n])
		log.Printf("Read Data: %d", n)
		_, err = w.WriteChunkedBody(buf)
		if err != nil {
			log.Fatalf("Error requesting: %s", err)
		}
	}
}
