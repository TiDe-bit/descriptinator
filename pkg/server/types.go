package server

import (
	"github.com/gin-gonic/gin"
)

type Versand string

const (
	VERSAND_PAKET      Versand = "paket"
	VERSAND_BRIEF              = "brief"
	VERSAND_BRIEFTAUBE         = "brieftaube"
)

type Parameter string

const (
	Q_AMOUT Parameter = "amout"
)

type IServer interface {
	HandleShipmentPath(gtx *gin.Context, engine *gin.Engine)
	Handler(method Versand) gin.HandlerFunc
	marshalParams(params gin.Params)
	getShipmentMethod(path string) Versand
}
