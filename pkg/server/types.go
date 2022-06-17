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

type Parameter int

const (
	Amount Parameter = iota + 1 // EnumIndex = 1
	Legal                       // EnumIndex = 2
)

func (p Parameter) String() string {
	return [...]string{"Amount", "Legal"}[p-1]
}

func (p Parameter) EnumIndex() int {
	return int(p)
}

type IServer interface {
	HandleShipmentPath(gtx *gin.Context, engine *gin.Engine)
	marshalParams(params gin.Params)
}
