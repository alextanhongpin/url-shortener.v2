package shorturlsvc_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alextanhongpin/url-shortener/domain/shorturl"
	"github.com/alextanhongpin/url-shortener/pkg/shorturlsvc"
	"github.com/alextanhongpin/url-shortener/pkg/shorturlsvc/internal"

	"github.com/stretchr/testify/assert"
)

func setupRouter(ctl *shorturlsvc.Controller, method, url string, body interface{}) *httptest.ResponseRecorder {
	var js []byte
	if body != nil {
		js, _ = json.Marshal(body)
	}
	r, _ := http.NewRequest(method, url, bytes.NewBuffer(js))
	w := httptest.NewRecorder()

	handler := ctl.Router()
	handler.ServeHTTP(w, r)
	return w
}

func TestGet(t *testing.T) {
	ucase := internal.NewMockUseCase()
	ctl := shorturlsvc.NewController(ucase)

	t.Run("returns error when get failed", func(t *testing.T) {
		resp := setupRouter(ctl, "GET", "/short-code", nil)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, "not implemented\n", resp.Body.String())
	})

	t.Run("redirects when short url exists", func(t *testing.T) {
		var (
			redirectURL = "http://abc.com"
		)
		ucase.Err = nil
		ucase.GetResponse = &shorturl.GetResponse{
			LongURL: redirectURL,
		}
		resp := setupRouter(ctl, "GET", "/short-code", nil)
		assert.Equal(t, http.StatusFound, resp.Code)
		assert.Equal(t, redirectURL, resp.Header().Get("Location"))
	})
}

func TestPost(t *testing.T) {
	ucase := internal.NewMockUseCase()
	ctl := shorturlsvc.NewController(ucase)
	data := shorturl.PutRequest{
		Code: "",
		URL:  "http://abc.com",
	}

	t.Run("returns error when create failed", func(t *testing.T) {
		resp := setupRouter(ctl, "POST", "/", data)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("returns short code when create succeed", func(t *testing.T) {
		code := "abc"
		ucase.Err = nil
		ucase.PutResponse = &shorturl.PutResponse{
			Code: code,
		}
		resp := setupRouter(ctl, "POST", "/", data)

		assert.Equal(t, http.StatusOK, resp.Code)

		b, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)

		var res shorturl.PutResponse
		assert.Nil(t, json.Unmarshal(b, &res))
		assert.Equal(t, code, res.Code)
	})
}

func TestCheckExists(t *testing.T) {
	// Prepare mock usecase.
	ucase := internal.NewMockUseCase()

	// Prepare controller.
	ctl := shorturlsvc.NewController(ucase)

	t.Run("returns error when failed", func(t *testing.T) {
		resp := setupRouter(ctl, "GET", "/short-code/search", nil)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("returns success when exists", func(t *testing.T) {
		// Override.
		ucase.Err = nil
		ucase.CheckExistsResponse = &shorturl.CheckExistsResponse{
			Exist: true,
		}

		resp := setupRouter(ctl, "GET", "/short-code/search", nil)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, `{"exist":true}`, strings.TrimSpace(resp.Body.String()))
	})
}
