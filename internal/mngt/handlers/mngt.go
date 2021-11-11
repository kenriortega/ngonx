package mngt

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	domain "github.com/kenriortega/ngonx/internal/mngt/domain"
	services "github.com/kenriortega/ngonx/internal/mngt/services"
	"github.com/kenriortega/ngonx/pkg/errors"
	"github.com/kenriortega/ngonx/pkg/logger"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		EnableCompression: false,
		CheckOrigin:       func(*http.Request) bool { return true },
	}
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

	endpoints, err := mh.service.ListEndpoints()

	if err != nil {
		logger.LogError(errors.Errorf("ngonx mngt: :%v", err).Error())

		writeResponse(w, http.StatusInternalServerError, err)

	}
	writeResponse(w, http.StatusOK, endpoints)

}

func (mh MngtHandler) RegisterEndpoint(data map[string]interface{}) {
	endpoint := domain.NewEnpoint(
		data["path_url"].(string),
		data["status"].(string),
	)

	endpoint.FromMapToJSON(data)
	err := mh.service.RegisterEndpoint(endpoint)
	if err != nil {
		logger.LogError(errors.Errorf("ngonx mngt: :%v", err).Error())
	}
}

func (mh MngtHandler) UpdateEndpoint(endpoint domain.Endpoint) {

	err := mh.service.UpdateEndpoint(endpoint)
	if err != nil {
		logger.LogError(errors.Errorf("ngonx mngt: :%v", err).Error())
	}
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.LogError(errors.Errorf("ngonx mngt: :%v", err).Error())
	}

}

func (mh MngtHandler) WssocketHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	ds := query.Get("ds")
	if ds == "" {
		ds = "10s"
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.LogError(errors.Errorf("ngonx mngt:  upgrade ws :%v", err).Error())
		return
	}
	defer c.Close()

	for {

		mt, message, err := c.ReadMessage()
		if err != nil {
			logger.LogError(errors.Errorf("ngonx mngt:  read ws :%v", err).Error())
			break
		}
		switch string(message) {

		case "endpoints":
			for {
				endpoints, _ := mh.service.ListEndpoints()
				bytes, err := json.Marshal(endpoints)
				if err != nil {
					logger.LogError(errors.Errorf("ngonx mngt:  ws :%v", err).Error())
				}
				err = c.WriteMessage(mt, bytes)
				if err != nil {
					logger.LogError(errors.Errorf("ngonx mngt: write endpoints ws :%v", err).Error())
				}
				durations, err := time.ParseDuration(ds)
				if err != nil {
					err = c.WriteMessage(mt, []byte(err.Error()))
					if err != nil {
						logger.LogError(errors.Errorf("ngonx mngt: write durations ws :%v", err).Error())
					}
				}
				time.Sleep(durations)
			}

		default:
			err = c.WriteMessage(mt, []byte("CMD not found"))
			if err != nil {
				logger.LogError(errors.Errorf("ngonx mngt: write default ws :%v", err).Error())
			}
		}
	}
}
