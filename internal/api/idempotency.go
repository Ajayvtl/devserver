package api

import (
	"bytes"
	"net/http"
	"sync"
)

type IdempotencyRecord struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

var idempotencyStore sync.Map

type idempotencyResponseWriter struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (w *idempotencyResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *idempotencyResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func IdempotencyMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("Idempotency-Key")
			if key == "" || r.Method == http.MethodGet || r.Method == http.MethodHead {
				next.ServeHTTP(w, r)
				return
			}

			if val, ok := idempotencyStore.Load(key); ok {
				rec := val.(*IdempotencyRecord)
				for k, vv := range rec.Headers {
					for _, v := range vv {
						w.Header().Add(k, v)
					}
				}
				w.WriteHeader(rec.StatusCode)
				w.Write(rec.Body)
				return
			}

			irw := &idempotencyResponseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(irw, r)

			if irw.status >= 200 && irw.status < 500 {
				idempotencyStore.Store(key, &IdempotencyRecord{
					StatusCode: irw.status,
					Headers:    irw.Header(),
					Body:       irw.body.Bytes(),
				})
			}
		})
	}
}
