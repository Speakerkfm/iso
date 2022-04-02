package admin

import (
	"encoding/json"
	"net/http"

	"github.com/Speakerkfm/iso/internal/pkg/models"
)

// HandleSaveServiceConfigs ...
func (i *Implementation) HandleSaveServiceConfigs(rw http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		return
	}

	var res []models.ServiceConfigDesc

	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	if err := i.ruleSvc.SaveServiceConfigs(r.Context(), res); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.WriteHeader(http.StatusOK)
}
