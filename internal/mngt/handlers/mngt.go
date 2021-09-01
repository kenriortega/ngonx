package mngt

import (
	"encoding/json"
	"net/http"

	domain "github.com/kenriortega/ngonx/internal/mngt/domain"
	services "github.com/kenriortega/ngonx/internal/mngt/services"
	"github.com/kenriortega/ngonx/pkg/logger"
)

type MngtHandler struct {
	service services.IMngtService
}

func NewMngtHandler(service services.IMngtService) MngtHandler {
	return MngtHandler{
		service: service,
	}
}

func (mh MngtHandler) GetAllEndpoints(w http.ResponseWriter, r *http.Request) {

	endpoints, err := mh.service.ListEnpoints()

	if err != nil {
		logger.LogError("handler: " + err.Error())
		writeResponse(w, http.StatusInternalServerError, err)

	}
	writeResponse(w, http.StatusOK, endpoints)

}

func (mh MngtHandler) RegisterEnpoint(data map[string]interface{}) {
	endpoint := domain.NewEnpoint(
		data["path_url"].(string),
		data["status"].(string),
	)

	endpoint.FromMapToJSON(data)
	err := mh.service.RegisterEnpoint(endpoint)
	if err != nil {
		logger.LogError(err.Error())
	}
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.LogError(err.Error())
	}

}
