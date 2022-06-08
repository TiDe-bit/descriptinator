package marshaller

import (
	"html/template"
	"os"
)

type Entry struct {
	Title    string ``
	SubTitle string
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

// ToDo: caching

func Marhal(fileNames []string, entry Entry) error {
	tmpl := template.Must(
		template.New("html-tmpl").ParseFiles(fileNames...),
	)

	file, err := os.Create(fileNames[0])
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, entry)

	return err
}
