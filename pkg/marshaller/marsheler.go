package marshaller

import (
	"context"
	"descriptinator/pkg/file_supply"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"html/template"
	"os"
)

var _ file_supply.Valid = &Entry{}

type Entry struct {
	ArtikelNum string               `bson:"artikelNum" json:"artikelNum"`
	Title      *string              `bson:"title" json:"title"`
	Subtitle   *string              `bson:"subtitle" json:"subtitle"`
	Article    file_supply.Article  `bson:"article" json:"article"`
	Shipping   file_supply.Shipping `bson:"shipping" json:"shipping"`
	Legal      file_supply.Legal    `bson:"legal" json:"legal"`
	Auction    file_supply.Auction  `bson:"auction" json:"auction"`
	Seller     file_supply.Seller   `bson:"seller" json:"seller"`
	Dsgvo      file_supply.Dsgvo    `bson:"dsgvo" json:"dsgvo"`
}

func (e *Entry) Byte() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		log.Error(err)
		return nil
	}
	return bytes
}

func (e *Entry) WithTitle(title *string) {
	e.Title = title
}
func (e *Entry) WithSubtitle(subtitle *string) {
	e.Subtitle = subtitle
}
func (e *Entry) WithShipping(ctx context.Context, shipping Versand, l file_supply.ITextLoader) {
	switch shipping {
	case VERSAND_BRIEF:
		e.Shipping = l.LoadBriefText(ctx)
		break
	case VERSAND_PAKET:
		e.Shipping = l.LoadPaketText(ctx)
		break
	case VERSAND_BRIEFTAUBE:
		e.Shipping = l.LoadPaketBrieftaube(ctx)
		break
	}
}

type Marshaller struct {
	url      string
	entry    *Entry
	tmplPath string
}

func DefaultEntry(id string) Entry {
	return Entry{
		ArtikelNum: id,
		Title:      nil,
		Subtitle:   nil,
		Article: file_supply.Article{
			GeneralInfo: nil,
			Description: nil,
			Fitting:     nil,
			Condition:   nil,
		},
		Shipping: nil,
		Legal:    nil,
		Auction:  nil,
		Seller:   nil,
		Dsgvo:    nil,
	}
}

func NewMarshaller() *Marshaller {
	return new(Marshaller)
}

func (m *Marshaller) CreatDescription() (file_supply.FileData, error) {
	rootPath, err := file_supply.GotoTmpl()
	if err != nil {
		return nil, err
	}

	tmplPath, err := file_supply.GetTmplFile(rootPath)
	if err != nil {
		return nil, err
	}

	err = marshalOne(tmplPath, m.entry)
	if err != nil {
		return nil, err
	}

	fileData, ok := file_supply.LoadFile(m.getFileName())
	if !ok {
		return nil, errors.New("no file found to load")
	}

	return fileData, nil
}

func (m *Marshaller) SetEntry(entry *Entry) {
	m.entry = entry
}

func (m *Marshaller) getFileName() string {
	return getFileDestination(m.entry)
}

func getFileDestination(entry *Entry) string {
	wd, _ := os.Getwd()
	return fmt.Sprintf("%s/html/%s.html", wd, entry.ArtikelNum)
}

func marshalOne(fileName string, entry *Entry) error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	tmpl, err := template.New("html-tmpl").Parse(string(file))
	if err != nil {
		return err
	}

	newFile, err := os.Create(getFileDestination(entry))
	if err != nil {
		return err
	}

	err = tmpl.Execute(newFile, entry)

	return err
}
