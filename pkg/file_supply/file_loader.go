package file_supply

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	defaultArticle Article = Article{
		GeneralInfo: "",
		Description: "",
		Fitting:     "",
		Condition:   "",
		Shipping:    "",
	}
)

// ToDo
var defaultEntry Entry = Entry{
	ArtikelNum: "",
	Title:      "",
	Subtitle:   "",
	Article:    defaultArticle,
	Shipping:   "Bitte holen Sie den gekauften Artikel innerhalb von einer Woche nach Auktionsende bei mir ab.\n\nWenn Sie Versand wünschen, dann betragen die Versandkosten innerhalb von Deutschland:\n\n1 Doppelpack (2 Packungen ohne Umkartons): 2,00 Euro per Post - Warensendung\nab 2 Doppelpacks (ab 4 Packungen ohne Umkartons): 4,00 Euro per DHL- Paket\n\nDie Übergabe Ihrer Bestellung an den Versanddienstleister erfolgt bei mir von Montag bis Freitag, einen Tag nach Zahlungseingang.\n\nLaufzeit von Post-Warensendungen: 7 Tage\nLaufzeit von Paketen: 2 Tage\n\nSollten Sie also schnelleren Versand wünschen, empfehle ich den Paketversand.\n\nDie Paketzustellung erfolgt innerhalb von 2 Tagen (Monat bis Freitag). Sie erhalten eine email mit einem Link zum Versanddienstleister zum Überprüfen des Versandstatus. Bei Postversand per Warensendung erhalten Sie eine email am Tag der Übergabe an die Post ohne Trackinglink.\nSollten Sie keine email von mir erhalten oder Fragen zum Versand Ihrer Sendung haben, setzen Sie sich bitte per email oder telefonisch mit mir in Verbindung.",
	Legal:      "Die Vertragssprache ist deutsch.\nVertragsabschluß bei \"Sofort-Kaufen\"-Angeboten ab dem 12.03.2014:\nDer Vertragsabschluß richtet sich nach § 6 Nr. 4 der Allgemeinen Geschäftsbedingungen von eBay, die ich wie folgt zitiere:\n\"Bei Festpreisartikeln nimmt der Käufer das Angebot an, indem er den Button \"Sofort-Kaufen\" anklickt und anschließend bestätigt. Bei Festpreisartikeln, bei denen der Verkäufer die Option \"sofortige Bezahlung\" ausgewählt hat, nimmt der Käufer das Angebot an, indem er den Button \"Sofort-Kaufen\" anklickt und den unmittelbar nachfolgenden Zahlungsvorgang abschließt. Der Käufer kann Angebote für mehrere Artikel auch dadurch annehmen, dass er die Artikel in den Warenkorb (sofern verfügbar) legt und den unmittelbar nachfolgenden Zahlungsvorgang abschließt.",
	Auction:    "Die Zahlungsabwicklung erfolgt gem. §4 Abs. 2 der eBay-AGB. Gezahlt werden kann mit allen von eBay zur Verfügung gestellten Zahlungsmethoden über eBay, nicht jedoch mehr per Überweisung direkt an den Verkäufer.",
	Seller:     "Dipl.-Ing. (FH) - Dipl.-Wirt.-Ing.(FH)\nThomas Dellmann\nGroße Aue 6\n32361 Preußisch Oldendorf\n\nVollständiges Impressum, sh. unten.\n\nOnline-Streitbeilegung\nOnline-Streitbeilegung gemäß Art. 14 Abs. 1 ODR-VO: Die Europäische Kommission stellt eine Plattform zur Online-Streitbeilegung (OS) bereit, zu der Sie über den Punkt \"Impressum\" unten gelangen.\nZur Teilnahme an einem Streitbeilegungsverfahren vor einer Verbraucherschlichtungsstelle sind wir nicht verpflichtet und nicht bereit.",
	Dsgvo:      "Ich unterrichte Sie hierdurch gemäß Telemediengesetz und DSGVO, dass ich personenbezogene Daten durch elektronische Datenverarbeitung (EDV) in dem zum Zwecke der Begründung, inhaltlichen Ausgestaltung oder Änderung des Kaufvertrages (Kaufabwicklung) erhebe, verarbeite und nutze.\nDarüberhinaus gebe ich Ihre Daten nicht an Dritte weiter. Zur ausführlichen Information klicken Sie bitte unten auf das Feld 'Datenschutz' oder in das Feld zu den rechtlichen Informationen des Verkäufers in die Zeile 'Allgemeine Geschäftsbedingungen' zu diesem Angebot.",
}

func ConnectToMongodb(ctx context.Context) (*mongo.Client, error) {
	//ToDo
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb:27017"))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewMongoTextLoader(ctx context.Context) ITextLoader {
	client, err := ConnectToMongodb(ctx)
	if err != nil {
		logrus.Fatal(err)
	}

	defaultsCollection, articlesCollection, err := setupDB(ctx, client)
	if err != nil {
		logrus.Fatal(err)
	}

	loader := &MongoTextLoader{
		articlesCollection,
		defaultsCollection,
		"",
	}

	// saving default Article
	loader.SaveEntry(ctx, &defaultEntry)

	return loader
}

func setupDB(ctx context.Context, client *mongo.Client) (*mongo.Collection, *mongo.Collection, error) {
	db := client.Database("texts")

	err := db.CreateCollection(ctx, "defaults", options.CreateCollection())
	if err != nil {
		return nil, nil, err
	}
	defaultsCollection := db.Collection("defaults")

	err = db.CreateCollection(ctx, "articles", options.CreateCollection())
	if err != nil {
		return nil, nil, err
	}
	articlesCollection := db.Collection("articles")

	return defaultsCollection, articlesCollection, nil
}

var _ ITextLoader = &MongoTextLoader{}

type MongoTextLoader struct {
	articleCollection  *mongo.Collection
	defaultsCollection *mongo.Collection
	currentArtikelNum  string
}

func (l *MongoTextLoader) InitiateAction(articleNumber string) ITextLoader {
	l.currentArtikelNum = articleNumber
	return l
}

func (l *MongoTextLoader) loadAny(ctx context.Context) (*Entry, error) {
	var target *Entry

	logrus.Debugf("loading entry '%s'", l.currentArtikelNum)

	filter := bson.M{"artikelNum": l.currentArtikelNum}

	err := l.articleCollection.FindOne(ctx, filter).Decode(target)
	if err != nil {
		return nil, err
	}

	return target, nil
}

func (l *MongoTextLoader) LoadEntry(ctx context.Context, shippingMethod string) *Entry {
	entry, err := l.loadAny(ctx)
	if err != nil {
		logrus.WithError(err).Warn("error loading entry")
		entry = new(Entry)
	}

	entry.Shipping = *l.LoadShippingText(ctx, shippingMethod)
	// ToDo
	switch true {
	case entry.ArtikelNum == "":
		break
	case entry.Auction == "":
		entry.Auction = defaultEntry.Auction
		break
	case entry.Legal == "":
		entry.Legal = defaultEntry.Legal
	case entry.Dsgvo == "":
		entry.Dsgvo = defaultEntry.Dsgvo
	case entry.Seller == "":
		entry.Seller = defaultEntry.Seller
	}

	return entry
}

func (l *MongoTextLoader) LoadTitleText(ctx context.Context) *string {
	entry, err := l.loadAny(ctx)
	if err != nil {
		logrus.WithError(err).Warn("error loading titel text")
		return nil
	}

	return &entry.Title
}

func (l *MongoTextLoader) LoadLegalText(ctx context.Context) *string {
	entry, err := l.loadAny(ctx)
	if err != nil {
		logrus.WithError(err).Warn("error loading legal text")
		return nil
	}

	return &entry.Legal
}

func (l *MongoTextLoader) LoadAuctionText(ctx context.Context) *string {
	entry, err := l.loadAny(ctx)
	if err != nil {
		logrus.WithError(err).Warn("error loading auction text")
		return nil
	}

	return &entry.Auction
}

type shippinmethodDoc struct {
	Data string `bson:"data"`
}

func mapShippingToPreset(shippingMethod string) string {
	resonse := ""
	switch shippingMethod {
	case VersandBrief.String():
		// TODO
		resonse = "Versandkosten (sh. auch unten):\n"
		break
	case VersandPaket.String():
		resonse = "Versandkosten (sh. auch unten):\n1 Doppelpack: 0 Euro per Post – Brief (oder 2 Euro Aufpreis für DHL-Paketversand)\nab 2 Doppelpacks ( 4 Packungen ohne Umkartons): 0 Euro per DHL-Paket\n\nBriefversand erfolgt geschützt im Briefumschlag, ca. 5 cm dick, bitte 2-3 Tage Zustellzeit seitens der Post einkalkulieren."
		break
	case VersandBrieftaube.String():
		resonse = "Viel Glück die Brieftaube zu fangen."
		break
	default:
		resonse = "not shipping"
		break
	}
	return resonse
}

func (l *MongoTextLoader) LoadShippingText(ctx context.Context, shippingMethod string) *string {
	var result *shippinmethodDoc
	err := l.defaultsCollection.FindOne(ctx, bson.M{"shipping": shippingMethod}).Decode(result)
	if err != nil {
		logrus.WithError(err).Error("loading shipping text")
		return nil
	}

	if result.Data == "" {
		preset := mapShippingToPreset(shippingMethod)
		l.SaveDefaultVersandText(ctx, Versand(shippingMethod), preset)
		if err != nil {
			logrus.Warn(err)
		}
		return &preset
	}

	return &result.Data
}

func (l *MongoTextLoader) LoadSellerText(ctx context.Context) *string {
	entry, err := l.loadAny(ctx)
	if err != nil {
		logrus.WithError(err).Warn("error loading seller text")
		return nil
	}

	return &entry.Seller
}

func (l *MongoTextLoader) SaveEntry(ctx context.Context, data *Entry) error {
	opts := options.Update().SetUpsert(true)

	update := bson.M{"$set": data}

	_, err := l.articleCollection.UpdateOne(ctx, bson.M{"artikelNum": data.ArtikelNum}, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func (l *MongoTextLoader) SaveDefaultVersandText(ctx context.Context, method Versand, data string) error {
	_, err := l.defaultsCollection.UpdateOne(
		ctx,
		bson.M{"shipping": method.String()},
		bson.M{"$set": bson.M{"data": data}},
	)
	if err != nil {
		return err
	}

	return nil
}
