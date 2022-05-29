package admin

import (
	"encoding/json"
	"net/http"
)

// HandleGetReport ...
func (i *Implementation) HandleGetReport(rw http.ResponseWriter, r *http.Request) {
	res, err := i.reportSvc.GetReport(r.Context())
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
