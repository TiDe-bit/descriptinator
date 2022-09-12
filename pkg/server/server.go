package server

import (
	"context"
	"descriptinator/pkg/file_supply"
	"descriptinator/pkg/marshaller"
	"descriptinator/pkg/server/api"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

var _ IServer = &ServeSenator{}

type ServeSenator struct {
	address    string
	marshaller *marshaller.Marshaller
	loader     file_supply.ITextLoader
}

func NewServinator(address string) *ServeSenator {
	ctx := context.Background()
	loader := file_supply.NewMongoTextLoader(ctx)

	return &ServeSenator{
		address:    address,
		marshaller: new(marshaller.Marshaller),
		loader:     loader,
	}
}

func (s *ServeSenator) Serve() {
	engine := gin.Default()

	engine.Group(
		"ebay",
		func(c *gin.Context) {
			s.HandleRoutes(c, engine)
		},
	)

	defer log.Info("Shutting down...")

	if err := engine.Run(s.address); err != nil {
		log.WithError(err).Fatal("Shutting down...")
	}
}

type queryParameterValue struct {
	Used  bool
	Value string
}

func (s *ServeSenator) HandleRoutes(gtx *gin.Context, engine *gin.Engine) {
	queryParams := extractQueryParams(gtx.Params)
	log.Debugf("params %+v", queryParams)

	fullPath := gtx.FullPath()
	fullPathSegments := strings.Split(fullPath, "/")
	artikelNr := fullPathSegments[len(fullPathSegments)-1]

	engine.Group(file_supply.VERSAND_BRIEF.String(), s.Handler(artikelNr, file_supply.VERSAND_BRIEF))
	engine.Group(string(file_supply.VERSAND_PAKET), s.Handler(artikelNr, file_supply.VERSAND_PAKET))
	engine.Group(file_supply.VERSAND_BRIEFTAUBE.String(), s.Handler(artikelNr, file_supply.VERSAND_BRIEFTAUBE))

	api.Run(engine)
}

func extractQueryParams(params gin.Params) map[marshaller.Parameter]queryParameterValue {
	marshalingOptions := make(map[marshaller.Parameter]queryParameterValue)
	queryParams := params

	for key := range marshalingOptions {
		value, ok := queryParams.Get(key.String())
		if ok {
			marshalingOptions[key] = queryParameterValue{
				Used:  ok,
				Value: value,
			}
		}
	}
	return marshalingOptions
}

func (s *ServeSenator) marshalParams(params gin.Params) {
	extraParams := make(map[marshaller.Parameter]string, len(params))
	for key := range extraParams {
		extraParams[key] = params.ByName(string(rune(key)))
	}
}

func (s *ServeSenator) Handler(artikelNr string, method file_supply.Versand) gin.HandlerFunc {
	var data file_supply.FileData
	var ok = false

	data, ok = file_supply.LoadFile(file_supply.FilePathFromArtikelNr(artikelNr))

	if !ok {
		entry := marshaller.DefaultEntry(artikelNr)
		// ToDo: lol
		// s.loader.LoadBriefText()
		// entry.WithTitle()

		s.marshaller.SetEntry(&entry)

		newFile, err := s.marshaller.CreatDescription()
		if err != nil {
			return sendFailure(err)
		}
		data = newFile
	}

	return sendDescription(data)
}

func sendFailure(err error) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.JSON(
			http.StatusInternalServerError,
			err.Error(),
		)
	}
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
