package server

import (
	"log"
	"net/http"

	"github.com/wk-y/rama-swap/microservices/scheduling"
)

type proxyTask struct {
	model string
	w     http.ResponseWriter
	r     *http.Request
}

func newProxyTask(model string, w http.ResponseWriter, r *http.Request) *proxyTask {
	return &proxyTask{
		model: model,
		w:     w,
		r:     r,
	}
}

// Model implements [scheduling.Task].
func (p *proxyTask) Model() string {
	return p.model
}

// PerformInference implements [scheduling.Task].
func (p *proxyTask) PerformInference(instance scheduling.Instance) error {
	log.Printf("Proxying request for model %v", p.model)
	// forward the request to the openai client
	instance.ReverseProxy().ServeHTTP(p.w, p.r)
	return nil
}

var _ scheduling.Task = (*proxyTask)(nil)
