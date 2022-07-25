package api

import (
	"descriptinator/pkg/file_supply"
	"descriptinator/pkg/marshaller"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func Run(engine *gin.Engine) {
	// ToDo api path
	engine.GET(getApiPath(marshaller.VERSAND_BRIEF.String()), loadVersand(marshaller.VERSAND_BRIEF.String()))
	engine.GET(getApiPath(marshaller.VERSAND_PAKET.String()), loadVersand(marshaller.VERSAND_PAKET.String()))
	engine.GET(getApiPath(marshaller.VERSAND_BRIEFTAUBE.String()), loadVersand(marshaller.VERSAND_BRIEFTAUBE.String()))
	engine.GET(getApiPath("entry"), func(gtx *gin.Context) {
		artikelNum, exists := gtx.Params.Get("artikelNum")
		if !exists {
			gtx.Error(errors.New("no artikelNum in path"))
			return
		}
		loadEntryFromDB(artikelNum)
	})

	// ToDo do other imports
	engine.PUT(getApiPath(marshaller.VERSAND_BRIEF.String()), saveVersand(marshaller.VERSAND_BRIEF.String()))
	engine.PUT(getApiPath(marshaller.VERSAND_PAKET.String()), saveVersand(marshaller.VERSAND_PAKET.String()))
	engine.PUT(getApiPath(marshaller.VERSAND_BRIEFTAUBE.String()), saveVersand(marshaller.VERSAND_BRIEFTAUBE.String()))
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
	return load[file_supply.Ttext]("versand", versandArt, false)
}

func saveVersand(versandArt string) gin.HandlerFunc {
	return save[file_supply.Ttext]("versand", versandArt, false)
}

func loadEntryFromDB(entryID string) gin.HandlerFunc {
	return load[*marshaller.Entry]("entry", entryID, true)
}

func saveEntry(entryID string) gin.HandlerFunc {
	return save[*marshaller.Entry]("entry", entryID, true)
}

func load[T file_supply.Valid](field, id string, special bool) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var data *T
		var err error
		if special {
			data, err = file_supply.LoadAny[T](
				gtx,
				bson.M{field: id},
				struct{}{},
			)
		} else {
			data, err = file_supply.LoadAny[T](
				gtx,
				bson.M{field: id},
			)
		}
		if err != nil {
			gtx.Error(err)
		}
		tmp := *data
		gtx.Data(http.StatusOK, "data", tmp.Byte())
	}
}

func save[T file_supply.Valid](field, id string, special bool) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var object *T

		rawData, err := gtx.GetRawData()
		if err != nil {
			gtx.Error(err)
		}

		err = json.Unmarshal(rawData, object)
		if err != nil {
			gtx.Error(err)
		}

		if special {
			err = file_supply.SaveAny[T](
				gtx,
				bson.M{field: id},
				object,
				struct{}{},
			)
		} else {
			err = file_supply.SaveAny[T](
				gtx,
				bson.M{field: id},
				object,
			)
		}
		if err != nil {
			gtx.Error(err)
		}

		gtx.Status(http.StatusOK)
	}
}
