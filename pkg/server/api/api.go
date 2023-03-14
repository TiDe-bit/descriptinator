package api

import (
	"descriptinator/pkg/file_supply"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Run(engine *gin.Engine, loader *file_supply.ITextLoader) {
	engine.GET(getApiPath(file_supply.VersandBrief.String()), loadVersand(file_supply.VersandBrief.String()))
	engine.GET(getApiPath(file_supply.VersandPaket.String()), loadVersand(file_supply.VersandPaket.String()))
	engine.GET(getApiPath(file_supply.VersandBrieftaube.String()), loadVersand(file_supply.VersandBrieftaube.String()))
	engine.GET(getApiPath("entry"), func(gtx *gin.Context) {
		artikelNum, exists := gtx.Params.Get("artikelNum")
		if !exists {
			gtx.Error(errors.New("no artikelNum in path"))
			return
		}
		loadEntryFromDB(artikelNum)
	})

	// ToDo do other imports
	engine.PUT(getApiPath(file_supply.VersandBrief.String()), saveVersand(file_supply.VersandBrief.String()))
	engine.PUT(getApiPath(file_supply.VersandPaket.String()), saveVersand(file_supply.VersandPaket.String()))
	engine.PUT(getApiPath(file_supply.VersandBrieftaube.String()), saveVersand(file_supply.VersandBrieftaube.String()))
	engine.PUT(getApiPath("entry"), func(gtx *gin.Context) {
		artikelNum, exists := gtx.Params.Get("artikelNum")
		if !exists {
			gtx.Error(errors.New("no artikelNum in path"))
			return
		}
		saveEntry(artikelNum)
	})
}

func getApiPath(endPoint string) string {
	return fmt.Sprintf("api/%s", endPoint)
}

func loadVersand(versandArt string) gin.HandlerFunc {
	return load("versand", versandArt, false)
}

func saveVersand(versandArt string) gin.HandlerFunc {
	return save("versand", versandArt, false)
}

func loadEntryFromDB(entryID string, loader file_supply.ITextLoader) gin.HandlerFunc {
	loader.LoadEntry(ctx, entryID)
}

func saveEntry(entryID string) gin.HandlerFunc {
	return save("entry", entryID, true)
}
