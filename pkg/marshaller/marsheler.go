package marshaller

import (
	"descriptinator/pkg/file_supply"
	"descriptinator/pkg/server"
	"errors"
	"fmt"
	"html/template"
	"os"
)

type Entry struct {
	KundenNr *string
	Title    *string ``
	Subtitle *string
	Article  Article
	shipping Shipping
	legal    Legal
	auction  Auction
	seller   Seller
	dsgvo    Dsgvo
}

type Article struct {
	GeneralInfo *string
	Description *string
	Fitting     *string
	Condition   *string
	shipping    Shipping
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
func (e *Entry) WithShipping(shipping server.Versand) {
	switch shipping {
	case server.VERSAND_BRIEF:
		e.shipping = file_supply.LoadBriefText()
		break
	case server.VERSAND_PAKET:
		e.shipping = file_supply.LoadPaketText()
		break
	case server.VERSAND_BRIEFTAUBE:
		e.shipping = file_supply.LoadPaketBrieftaube()
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
		KundenNr: &id,
		Title:    nil,
		Subtitle: nil,
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
	return fmt.Sprintf("%s/html/%s.html", wd, entry.KundenNr)
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
