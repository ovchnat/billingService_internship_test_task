package apiserver

import (
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIServer_DepositMoney(t *testing.T) {
	srv := NewServer()
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sayHello", nil)
	srv.httpServer.Handler.ServeHTTP(rec, req)
	assert.Equal(t, rec.Body.String(), "Hello World!")
}
