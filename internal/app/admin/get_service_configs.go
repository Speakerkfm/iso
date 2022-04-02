package admin

import (
	"encoding/json"
	"net/http"
)

// HandleGetServiceConfigs ...
func (i *Implementation) HandleGetServiceConfigs(rw http.ResponseWriter, r *http.Request) {
	res, err := i.ruleSvc.GetServiceConfigs(r.Context())
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	if err := json.NewEncoder(rw).Encode(res); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.WriteHeader(http.StatusOK)
}
