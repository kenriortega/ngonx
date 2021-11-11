package mngt

import (
	"encoding/json"

	"github.com/kenriortega/ngonx/pkg/errors"
	"github.com/kenriortega/ngonx/pkg/logger"
	uuid "github.com/satori/go.uuid"
)

type Endpoint struct {
	ID      string `json:"id,omitempty"`
	PathUrl string `json:"path_url,omitempty"`
	Status  string `json:"status,omitempty"`
}
type IEndpoint interface {
	ListEndpoints() ([]Endpoint, error)
	RegisterEndpoint(Endpoint) error
	UpdateEndpoint(Endpoint) error
}

func NewEnpoint(pathUrl, status string) Endpoint {
	myuuid := uuid.NewV4().String()

	return Endpoint{
		ID:      myuuid,
		PathUrl: pathUrl,
		Status:  status,
	}
}

func (ed *Endpoint) FromMapToJSON(data map[string]interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		logger.LogError(errors.Errorf("ngonx mngt: :%v", err).Error())

	}

	err = json.Unmarshal(b, &ed)
	if err != nil {
		logger.LogError(errors.Errorf("ngonx mngt: :%v", err).Error())

	}
}
func (ed *Endpoint) ToMAP() (toHashMap map[string]interface{}, err error) {

	fromStruct, err := json.Marshal(ed)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(fromStruct, &toHashMap); err != nil {
		return toHashMap, err
	}

	return toHashMap, nil
}
