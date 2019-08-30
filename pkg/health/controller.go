package health

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alextanhongpin/url-shortener/domain"
	"github.com/go-chi/chi"
)

var DeployedAt = time.Now()

type Controller struct {
	version string
}

func NewController(version string) *Controller {
	return &Controller{version}
}

func (ctl *Controller) GetHealth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(domain.Health{
		DeployedAt: DeployedAt,
		Uptime:     time.Since(DeployedAt).String(),
		Version:    ctl.version,
	})
}

func (ctl *Controller) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", ctl.GetHealth)
	return r
}
