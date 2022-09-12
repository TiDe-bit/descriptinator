package marshaller

import (
	"descriptinator/pkg/file_supply"
	"errors"
	"fmt"
	"html/template"
	"os"
)

type Marshaller struct {
	url      string
	entry    *file_supply.Entry
	tmplPath string
}

func DefaultEntry(id string) file_supply.Entry {
	return file_supply.Entry{
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

func (m *Marshaller) SetEntry(entry *file_supply.Entry) {
	m.entry = entry
}

func (m *Marshaller) getFileName() string {
	return getFileDestination(m.entry)
}

func getFileDestination(entry *file_supply.Entry) string {
	wd, _ := os.Getwd()
	return fmt.Sprintf("%s/html/%s.html", wd, entry.ArtikelNum)
}

func marshalOne(fileName string, entry *file_supply.Entry) error {
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
