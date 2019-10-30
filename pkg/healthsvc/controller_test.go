package healthsvc_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alextanhongpin/url-shortener/domain/health"
	"github.com/alextanhongpin/url-shortener/pkg/healthsvc"

	"github.com/stretchr/testify/assert"
)

func TestController(t *testing.T) {
	assert := assert.New(t)

	version := "abc"

	ctl := healthsvc.NewController(version)
	r, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(ctl.GetHealth)
	handler.ServeHTTP(w, r)

	assert.Equal(http.StatusOK, w.Code)
	var req health.Health
	err := json.Unmarshal(w.Body.Bytes(), &req)
	assert.Nil(err)
	assert.Equal(healthsvc.DeployedAt.Format(time.RFC3339), req.DeployedAt.Format(time.RFC3339))
	assert.Equal(version, req.Version)
}
