package index_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Luzilla/dnsbl_exporter/internal/index"
	"github.com/stretchr/testify/assert"
)

func TestIndexHandler(t *testing.T) {
	iHandler := index.IndexHandler{
		Name:    "test",
		Version: "1.2.3",
		Path:    "/something",
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	iHandler.Handler(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Contains(t, string(data), "test")
	assert.Contains(t, string(data), "1.2.3")
	assert.Contains(t, string(data), "/something")
}
