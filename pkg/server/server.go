package server

import (
	"descriptinator/pkg/file_supply"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var _ IServer = &ServeSenator{}

type ServeSenator struct {
	address string
}

func NewServinator(address string) *ServeSenator {
	return &ServeSenator{
		address: address,
	}
}

func (s *ServeSenator) Serve() {
	engine := gin.Default()

	engine.Group(
		"ebay",
		func(c *gin.Context) {
			s.HandleShipmentPath(c, engine)
		},
	)

	err := engine.Run(s.address)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *ServeSenator) HandleShipmentPath(gtx *gin.Context, engine *gin.Engine) {
	queryParams := gtx.Params

	engine.Group(VERSAND_BRIEF, s.Handler(VERSAND_BRIEF))
	engine.Group(string(VERSAND_PAKET), s.Handler(VERSAND_PAKET))
	engine.Group(VERSAND_BRIEFTAUBE, s.Handler(VERSAND_BRIEFTAUBE))

}

func (s *ServeSenator) marshalParams(params gin.Params) {
	extraParams := make(map[Parameter]string, len(params))
	for key := range extraParams {
		extraParams[key] = params.ByName(string(key))
	}
}

func (s *ServeSenator) Handler(method Versand) gin.HandlerFunc {
	switch method {
	case VERSAND_BRIEF:

	}
}

func (s *ServeSenator) getShipmentMethod(path string) Versand {

}

func sendDescription(data file_supply.FileData) func(g *gin.Context) {
	return func(g *gin.Context) {
		g.Data(
			http.StatusOK,
			http.DetectContentType(*data),
			*data,
		)
	}

}
