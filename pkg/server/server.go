package server

import (
	"context"
	"descriptinator/pkg/file_supply"
	"descriptinator/pkg/marshaller"
	"descriptinator/pkg/server/api"
	"net/http"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

	basePath := "./app/dist"

	engine.Use(static.ServeRoot("/", basePath))
	engine.Use(static.ServeRoot("/edit", basePath))

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

	engine.Group(file_supply.VersandBrief.String(), s.Handler(gtx, artikelNr, file_supply.VersandBrief))
	engine.Group(file_supply.VersandPaket.String(), s.Handler(gtx, artikelNr, file_supply.VersandPaket))
	engine.Group(file_supply.VersandBrieftaube.String(), s.Handler(gtx, artikelNr, file_supply.VersandBrieftaube))

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

func (s *ServeSenator) Handler(ctx context.Context, artikelNr string, method file_supply.Versand) gin.HandlerFunc {

	data, ok := file_supply.LoadFile(file_supply.FilePathFromArtikelNr(artikelNr))

	if !ok {

		savedEntry := s.loader.LoadEntry(ctx, method)

		s.marshaller.SetEntry(savedEntry)

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
