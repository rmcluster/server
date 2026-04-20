package openapi

import (
	"github.com/gin-gonic/gin"
	internal "github.com/wk-y/rama-swap/server/openapi/internal/generated/go"
	"github.com/wk-y/rama-swap/tracker"
)

func NewRouter() *gin.Engine {
	return internal.NewRouter(internal.ApiHandleFunctions{
		DefaultAPI: OpenAPIRoutes{},
	})
}

type OpenAPIRoutes struct{}

// AnnounceGet implements [openapi.DefaultAPI].
func (o OpenAPIRoutes) AnnounceGet(c *gin.Context) {
	tracker.DefaultTracker.Announce(c)
}

var _ internal.DefaultAPI = OpenAPIRoutes{}
