package confapi

import "github.com/gin-gonic/gin"

type ConfAPI struct {
}

func NewConfApi(router *gin.IRouter) *ConfAPI {
	return &ConfAPI{}
}

func (c *ConfAPI) GetDeviceList(router *gin.IRouter) {

}
