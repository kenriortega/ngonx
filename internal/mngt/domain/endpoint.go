package mngt

import (
	"encoding/json"

	"github.com/kenriortega/ngonx/pkg/logger"
)

type Endpoint struct {
	PathUrl string `json:"path_url,omitempty"`
	Status  string `json:"status,omitempty"`
}
type IEndpoint interface {
	ListEnpoints() ([]Endpoint, error)
	RegisterEnpoint(Endpoint) error
}

func NewEnpoint(pathUrl, status string) Endpoint {
	return Endpoint{
		PathUrl: pathUrl,
		Status:  status,
	}
}

func (ed *Endpoint) FromMapToJSON(data map[string]interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		logger.LogError(err.Error())
	}

	err = json.Unmarshal(b, &ed)
	if err != nil {
		logger.LogError(err.Error())
	}
}
