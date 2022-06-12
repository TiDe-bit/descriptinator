package file_supply

import (
	"errors"
	"fmt"
	"html/template"
	"os"
)

type Entry struct {
	ID       string
	Title    string ``
	Subtitle string
	Article
	Shipping
	Legal
	Auction
	Seller
	Dsgvo
}

type Article struct {
	GeneralInfo string
	Description string
	Fitting     string
	Condition   string
	Shipping
}

type Shipping string

type Auction string
type Legal string
type Seller string
type Dsgvo string

type Marshaller struct {
	url      string
	entry    Entry
	tmplPath string
}

func NewMarshaller() *Marshaller {
	return new(Marshaller)
}

func (m *Marshaller) CreatDescription() (FileData, error) {
	rootPath, err := gotoTmpl()
	if err != nil {
		return nil, err
	}

	tmplPath, err := getTmplPath(rootPath)
	if err != nil {
		return nil, err
	}

	err = marshalOne(tmplPath, m.entry)
	if err != nil {
		return nil, err
	}

	fileData, ok := LoadFile(m.getFileName())
	if !ok {
		return nil, errors.New("no file found to load")
	}

	return fileData, nil
}

func (m *Marshaller) getFileName() string {
	return getFileName(m.entry)
}

func getFileName(entry Entry) string {
	return fmt.Sprintf("%s.html", entry.ID)
}

func marshalOne(fileName string, entry Entry) error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	tmpl, err := template.New("html-tmpl").Parse(string(file))
	if err != nil {
		return err
	}

	newFile, err := os.Create(getFileName(entry))
	if err != nil {
		return err
	}

	err = tmpl.Execute(newFile, entry)

	return err
}
