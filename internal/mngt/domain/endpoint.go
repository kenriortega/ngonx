package mngt

import (
	"encoding/json"

	"github.com/kenriortega/ngonx/pkg/logger"
)

type Endpoint struct {
	pathUrl string
	status  string
}
type IEndpoint interface {
	ListEnpoints() ([]Endpoint, error)
	RegisterEnpoint(Endpoint) error
}

func NewEnpoint(pathUrl, status string) *Endpoint {
	return &Endpoint{
		pathUrl: pathUrl,
		status:  status,
	}
}

func (ed *Endpoint) PathUrl() string           { return ed.pathUrl }
func (ed *Endpoint) SetPathUrl(pathUrl string) { ed.pathUrl = pathUrl }
func (ed *Endpoint) Status() string            { return ed.status }
func (ed *Endpoint) SetStatus(status string)   { ed.status = status }

func (ed *Endpoint) FromMapToJSON(data map[string]interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		logger.LogError(err.Error())
	}

	json.Unmarshal(b, &ed)
}
