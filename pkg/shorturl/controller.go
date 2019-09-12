package shorturl

import (
	"encoding/json"
	"net/http"

	"github.com/alextanhongpin/url-shortener/app/logger"
	"github.com/alextanhongpin/url-shortener/domain"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type Controller struct {
	service domain.Service
}

func NewController(service domain.Service) *Controller {
	return &Controller{service}
}

func (ctl *Controller) GetRedirect(w http.ResponseWriter, r *http.Request) {
	var req domain.GetRequest
	req.Code = chi.URLParam(r, "code")

	ctx := r.Context()
	log := logger.WithContext(ctx, zap.Object("req", &req))

	res, err := ctl.service.Get(ctx, req)
	if err != nil {
		log.Error("getRedirectError", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, res.LongURL, http.StatusFound)
}

func (ctl *Controller) PostCreate(w http.ResponseWriter, r *http.Request) {
	var req domain.PutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	log := logger.WithContext(ctx, zap.Object("req", &req))

	res, err := ctl.service.Put(ctx, req)
	if err != nil {
		log.Error("postCreateError", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (ctl *Controller) GetSearch(w http.ResponseWriter, r *http.Request) {
	var req domain.CheckExistsRequest
	req.Code = chi.URLParam(r, "code")

	ctx := r.Context()
	log := logger.WithContext(ctx, zap.Object("req", &req))

	res, err := ctl.service.CheckExists(ctx, req)
	if err != nil {
		log.Error("checkExistsError", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (ctl *Controller) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/{code}/search", ctl.GetSearch)
	r.Get("/{code}", ctl.GetRedirect)
	r.Post("/", ctl.PostCreate)
	return r
}
