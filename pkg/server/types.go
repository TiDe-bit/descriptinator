package server

import "github.com/gin-gonic/gin"

type IServer interface {
	HandleShipmentPath(gtx *gin.Context, engine *gin.Engine)
	marshalParams(params gin.Params)
}
