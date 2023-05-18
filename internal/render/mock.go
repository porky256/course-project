package render

import "net/http"

type mockWriter struct {
}

func (mw *mockWriter) Header() http.Header {
	return http.Header{}
}

func (mw *mockWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func (mw *mockWriter) WriteHeader(statusCode int) {
}
