package server

import (
	"descriptinator/pkg/file_supply"
	"descriptinator/pkg/marshaller"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var _ IServer = &ServeSenator{}

type ServeSenator struct {
	address    string
	marshaller *marshaller.Marshaller
}

func NewServinator(address string) *ServeSenator {
	return &ServeSenator{
		address:    address,
		marshaller: new(marshaller.Marshaller),
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

type QueryParameterValue struct {
	Used  bool
	Value string
}

func (s *ServeSenator) HandleShipmentPath(gtx *gin.Context, engine *gin.Engine) {
	queryParams := extractQueryParams(gtx.Params)
	log.Debugf("params %+v", queryParams)

	engine.Group(VERSAND_BRIEF, s.Handler(VERSAND_BRIEF))
	engine.Group(string(VERSAND_PAKET), s.Handler(VERSAND_PAKET))
	engine.Group(VERSAND_BRIEFTAUBE, s.Handler(VERSAND_BRIEFTAUBE))

}

func extractQueryParams(params gin.Params) map[Parameter]QueryParameterValue {
	marshalingOptions := make(map[Parameter]QueryParameterValue)
	queryParams := params

	for key := range marshalingOptions {
		value, ok := queryParams.Get(key.String())
		if ok {
			marshalingOptions[key] = QueryParameterValue{
				Used:  ok,
				Value: value,
			}
		}
	}
	return marshalingOptions
}

func (s *ServeSenator) marshalParams(params gin.Params) {
	extraParams := make(map[Parameter]string, len(params))
	for key := range extraParams {
		extraParams[key] = params.ByName(string(rune(key)))
	}
}

func (s *ServeSenator) Handler(artikelNr string, method Versand) gin.HandlerFunc {
	var data file_supply.FileData
	switch method {
	case VERSAND_BRIEF:

	}

	entry := marshaller.DefaultEntry(artikelNr)
	entry.WithTitle()

	s.marshaller.CreatDescription()

	return sendDescription(data)
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
