package confapi

import "github.com/gin-gonic/gin"

type ConfAPIServer struct {
	// store
}

func NewConfAPIServer() *ConfAPIServer {
	return &ConfAPIServer{}
}

func (s *ConfAPIServer) GetCameraList(c *gin.Context) {
}

func (s *ConfAPIServer) CreateCamera(c *gin.Context) {
}

func (s *ConfAPIServer) DeleteCamera(c *gin.Context) {
}
