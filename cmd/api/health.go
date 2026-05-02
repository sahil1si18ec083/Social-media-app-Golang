package main

import (
	"net/http"
	"time"
)

// HealthResponse represents the health check response payload.
type HealthResponse struct {
	Status  string `json:"status" example:"ok"`
	Version string `json:"version" example:"1.0.0"`
}

// healthCheckHandler godoc
//
//	@Summary		Health check
//	@Description	Returns server health status and version information
//	@Tags			system
//	@ID				health-check
//	@Produce		json
//	@Success		200	{object}	HealthResponse	"Server is healthy"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/health [get]
func (a *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	time.Sleep(time.Second * 3)
	data := HealthResponse{
		Status:  "ok",
		Version: "version",
	}

	err := writeJSON(w, http.StatusOK, data)
	if err != nil {
		a.internalServerError(w, r, err)
	}

}
