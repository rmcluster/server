package microservices

import "net/http"

type Microservice interface {
	RegisterHandlers(mux *http.ServeMux)
}
