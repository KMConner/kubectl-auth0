package web

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

type LoginResult struct {
	Token        string
	RefreshToken string
}

func WaitCallback(callback func(map[string][]string) (*LoginResult, error)) (*LoginResult, error) {
	m := http.NewServeMux()
	s := http.Server{
		Addr:    ":8088",
		Handler: m,
	}
	var result *LoginResult
	m.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		var err error = nil
		result, err = callback(r.URL.Query())
		if err == nil {
			_, err := w.Write([]byte("<h1>OK</h1>"))
			if err != nil {
				log.Printf("failed to write response: %v\n", err)
			}

			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, time.Second)

			ch := make(chan struct{})
			go func() {
				<-ch
				<-time.After(time.Millisecond * 100)
				err = s.Shutdown(ctx)
				if err != nil {
					log.Printf("error on shutdown the server: %v\n", err)
				}
				cancel()
			}()
			ch <- struct{}{}
		} else {
			_, err := w.Write([]byte("<h1>Error !</h1>"))
			if err != nil {
				log.Printf("failed to write response: %v\n", err)
			}
		}
	})
	err := s.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return result, nil
	}
	return result, err
}
