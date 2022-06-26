package marshaller

import (
	"context"
	"descriptinator/pkg/file_supply"
	"errors"
	"fmt"
	"html/template"
	"os"
)

type Entry struct {
	ArtikelNum string   `bson:"artikelNum"`
	Title      *string  `bson:"title"`
	Subtitle   *string  `bson:"subtitle"`
	Article    Article  `bson:"article"`
	shipping   Shipping `bson:"shipping"`
	legal      Legal    `bson:"legal"`
	auction    Auction  `bson:"auction"`
	seller     Seller   `bson:"seller"`
	dsgvo      Dsgvo    `bson:"dsgvo"`
}

type Article struct {
	GeneralInfo *string  `bson:"generalInfo"`
	Description *string  `bson:"description"`
	Fitting     *string  `bson:"fitting"`
	Condition   *string  `bson:"condition"`
	shipping    Shipping `bson:"shipping"`
}

type Shipping *string

type Auction *string
type Legal *string
type Seller *string
type Dsgvo *string

func (e *Entry) WithTitle(title *string) {
	e.Title = title
}
func (e *Entry) WithSubtitle(subtitle *string) {
	e.Subtitle = subtitle
}
func (e *Entry) WithShipping(ctx context.Context, shipping Versand, l file_supply.ITextLoader) {
	switch shipping {
	case VERSAND_BRIEF:
		e.shipping = l.LoadBriefText(ctx)
		break
	case VERSAND_PAKET:
		e.shipping = l.LoadPaketText(ctx)
		break
	case VERSAND_BRIEFTAUBE:
		e.shipping = l.LoadPaketBrieftaube(ctx)
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
		ArtikelNum: &id,
		Title:      nil,
		Subtitle:   nil,
		Article: Article{
			GeneralInfo: nil,
			Description: nil,
			Fitting:     nil,
			Condition:   nil,
			shipping:    nil,
		},
		shipping: nil,
		legal:    nil,
		auction:  nil,
		seller:   nil,
		dsgvo:    nil,
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
	return fmt.Sprintf("%s/html/%s.html", wd, *entry.ArtikelNum)
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
