package marshaller

import (
	"html/template"
	"os"
)

type Entry struct {
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

// ToDo: caching

func Marshal(fileNames []string, entry Entry) error {
	file, err := os.ReadFile(fileNames[0])
	if err != nil {
		return err
	}

	tmpl, err := template.New("html-tmpl").Parse(string(file))
	if err != nil {
		return err
	}

	// ToDo: cache with file
	newFile, err := os.Create(fileNames[0] + ".html")
	if err != nil {
		return err
	}

	err = tmpl.Execute(newFile, entry)

	return err
}
