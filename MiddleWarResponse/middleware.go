package middlewarresponse

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func MiddleWare() {
	srv := http.Server{
		Addr:    ":29091",
		Handler: timeoutMiddleware(http.HandlerFunc(func1), customWriteTimeout),
	}
	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}

func func1(w http.ResponseWriter, req *http.Request) {
	log.Println("starttime", time.Now())
	time.Sleep(12 * time.Second)
	fmt.Println("My func Println")
	log.Println("printtime", time.Now())
	io.WriteString(w, "My func!\n")
	log.Println("stringttime", time.Now())
}

func customWriteTimeout(r *http.Request) time.Duration {
	// Example: Set no timeout for a particular endpoint
	if r.URL.Path == "/noTimeout" {
		return -1
	}
	// Default write timeout
	return 10 * time.Second
}

func timeoutMiddleware(next http.Handler, getWriteTimeout func(*http.Request) time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeTimeout := getWriteTimeout(r)

		if writeTimeout < 0 {
			// No timeout for this endpoint
			next.ServeHTTP(w, r)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), writeTimeout)
		defer cancel()
		r = r.WithContext(ctx)
		done := make(chan struct{})

		go func() {
			next.ServeHTTP(w, r)
			close(done) // Signal that handling is finished
		}()

		select {
		case <-done:
			// request completed
		case <-ctx.Done():
			// request timed out
			w.WriteHeader(http.StatusGatewayTimeout)
			// handle timeout, like logging or sending a specific response
		}
	})
}
