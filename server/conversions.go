package server

import (
	"fmt"
	"time"

	"github.com/wk-y/rama-swap/llama"
)

func convertModelList(ramaModels []llama.Model) ([]Model, error) {
	var models []Model
	for _, ramaModel := range ramaModels {
		t, err := time.Parse(time.RFC3339, ramaModel.Modified)
		if err != nil {
			return nil, fmt.Errorf("failed to parse model timestamp %#v: %v", ramaModel.Modified, err)
		}
		models = append(models, Model{
			Id:      ramaModel.Name,
			Object:  "model",
			Created: int(t.Unix()),
			OwnedBy: "rama-swap",
		})
	}
	return models, nil
}
