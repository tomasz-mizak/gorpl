package model

import (
	"encoding/xml"
)

// Common types based on commonTypes.xsd
type BooleanAsString string   // "TAK", "NIE", ""
type BigIntAsString string    // max length 19
type IntegerAsString string   // max length 10
type LimitedString string     // max length 255
type DateAsString string      // max length 10
type DeletedAsString string   // "Skasowane", ""
type ChangeTypeString string  // "Nowy", "Zmodyfikowany", "Usuniety"

// Main structures based on schematXmlRejestr_Produktow_Leczniczych.xsd
type ProduktyLecznicze struct {
	XMLName         xml.Name          `xml:"http://rejestry.ezdrowie.gov.pl/rpl/eksport-danych-v6.0.0 produktyLecznicze"`
	StanNaDzien     DateAsString      `xml:"stanNaDzien,attr"`
	ProduktyLecznicze []ProduktLeczniczy `xml:"produktLeczniczy"`
}

type ProduktLeczniczy struct {
	KodyATC                  *KodyATC                 `xml:"kodyATC"`
	DrogiPodania             *DrogiPodania            `xml:"drogiPodania"`
	SubstancjeCzynne         *SubstancjeCzynne        `xml:"substancjeCzynne"`
	Opakowania               *Opakowania              `xml:"opakowania"`
	DaneOWytworcy            *DaneOWytworcy           `xml:"daneOWytworcy"`
	MaterialyEdukacyjne      *MaterialyEdukacyjne     `xml:"materialyEdukacyjne"`
	
	NazwaProduktu            LimitedString           `xml:"nazwaProduktu,attr"`
	RodzajPreparatu          LimitedString           `xml:"rodzajPreparatu,attr"`
	NazwaPowszechnieStosowana LimitedString          `xml:"nazwaPowszechnieStosowana,attr"`
	NazwaPoprzedniaProduktu   LimitedString          `xml:"nazwaPoprzedniaProduktu,attr"`
	Moc                      string                  `xml:"moc,attr"`
	NazwaPostaciFarmaceutycznej LimitedString        `xml:"nazwaPostaciFarmaceutycznej,attr"`
	PodmiotOdpowiedzialny    string                  `xml:"podmiotOdpowiedzialny,attr"`
	TypProcedury             LimitedString           `xml:"typProcedury,attr"`
	NumerPozwolenia          LimitedString           `xml:"numerPozwolenia,attr"`
	WaznoscPozwolenia        string                  `xml:"waznoscPozwolenia,attr"`
	PodstawaPrawna           string                  `xml:"podstawaPrawna,attr"`
	ZakazStosowaniaUZwierzat BooleanAsString         `xml:"zakazStosowaniaUZwierzat,attr"`
	Ulotka                   string                  `xml:"ulotka,attr"`
	Charakterystyka          string                  `xml:"charakterystyka,attr"`
	EtykietoUlotka           string                  `xml:"etykietoUlotka,attr"`
	UlotkaImportRownolegly   string                  `xml:"ulotkaImportRownolegly,attr"`
	EtykietoUlotkaImportRownolegly string            `xml:"etykietoUlotkaImportRownolegly,attr"`
	OznaczenieOpakowanImportRownolegly string        `xml:"oznaczenieOpakowanImportRownolegly,attr"`
	ID                       BigIntAsString          `xml:"id,attr"`
	Status                   ChangeTypeString        `xml:"status,attr"`
}

type KodyATC struct {
	KodATC []LimitedString `xml:"kodATC"`
}

type DrogiPodania struct {
	DrogaPodania []DrogaPodania `xml:"drogaPodania"`
}

type DrogaPodania struct {
	Gatunki            *Gatunki `xml:"gatunki"`
	DrogaPodaniaNazwa  string   `xml:"drogaPodaniaNazwa,attr"`
}

type Gatunki struct {
	Gatunek []Gatunek `xml:"gatunek"`
}

type Gatunek struct {
	OkresyKarencji  *OkresyKarencji `xml:"okresyKarencji"`
	NazwaGatunku    string          `xml:"nazwaGatunku,attr"`
}

type OkresyKarencji struct {
	OkresKarencji []OkresKarencji `xml:"okresKarencji"`
}

type OkresKarencji struct {
	NazwaTkanki     string `xml:"nazwaTkanki,attr"`
	WartoscMiary    string `xml:"wartoscMiary,attr"`
	JednostkaMiary  string `xml:"jednostkaMiary,attr"`
}

type SubstancjeCzynne struct {
	SubstancjaCzynna []SubstancjaCzynna `xml:"substancjaCzynna"`
}

type SubstancjaCzynna struct {
	Value                       string `xml:",chardata"`
	NazwaSubstancji             string `xml:"nazwaSubstancji,attr"`
	IloscSubstancji             string `xml:"iloscSubstancji,attr"`
	JednostkaMiaryIlosciSubstancji string `xml:"jednostkaMiaryIlosciSubstancji,attr"`
	IloscPreparatu              string `xml:"iloscPreparatu,attr"`
	JednostkaMiaryIlosciPreparatu string `xml:"jednostkaMiaryIlosciPreparatu,attr"`
	InnyOpisIlosci              string `xml:"innyOpisIlosci,attr"`
}

type Opakowania struct {
	Opakowanie []Opakowanie `xml:"opakowanie"`
}

type Opakowanie struct {
	JednostkiOpakowania  *JednostkiOpakowania `xml:"jednostkiOpakowania"`
	ZgodyPrezesa         *ZgodyPrezesa        `xml:"zgodyPrezesa"`
	
	KodGTIN              LimitedString       `xml:"kodGTIN,attr"`
	KategoriaDostepnosci LimitedString       `xml:"kategoriaDostepnosci,attr"`
	Skasowane            BooleanAsString     `xml:"skasowane,attr"`
	NumerEu              LimitedString       `xml:"numerEu,attr"`
	DystrybutorRownolegly LimitedString      `xml:"dystrybutorRownolegly,attr"`
	ID                   BigIntAsString      `xml:"id,attr"`
}

type JednostkiOpakowania struct {
	JednostkaOpakowania []JednostkaOpakowania `xml:"jednostkaOpakowania"`
}

type JednostkaOpakowania struct {
	Value               string         `xml:",chardata"`
	LiczbaOpakowan      IntegerAsString `xml:"liczbaOpakowan,attr"`
	RodzajOpakowania    LimitedString  `xml:"rodzajOpakowania,attr"`
	Pojemnosc           string         `xml:"pojemnosc,attr"`
	JednostkaPojemnosci LimitedString  `xml:"jednostkaPojemnosci,attr"`
	InformacjeDodatkowe string         `xml:"informacjeDodatkowe,attr"`
}

type ZgodyPrezesa struct {
	ZgodaPrezesa []ZgodaPrezesa `xml:"zgodaPrezesa"`
}

type ZgodaPrezesa struct {
	NrZgodyPrezesa    LimitedString     `xml:"nrZgodyPrezesa"`
	GTINZagraniczne   *GTINZagraniczne  `xml:"GTINZagraniczne"`
}

type GTINZagraniczne struct {
	GTINZagraniczny []GTINZagraniczny `xml:"GTINZagraniczny"`
}

type GTINZagraniczny struct {
	Numer string `xml:"numer,attr"`
}

type DaneOWytworcy struct {
	Wytworcy []Wytworcy `xml:"wytworcy"`
}

type Wytworcy struct {
	Value                        string `xml:",chardata"`
	NazwaWytworcyImportera       string `xml:"nazwaWytworcyImportera,attr"`
	KrajWytworcyImportera        string `xml:"krajWytworcyImportera,attr"`
	PodmiotOdpowiedzialnywKrajuEksportu string `xml:"podmiotOdpowiedzialnywKrajuEksportu,attr"`
	KrajEksportu                 string `xml:"krajEksportu,attr"`
}

type MaterialyEdukacyjne struct {
	DlaPacjenta *DlaPacjenta `xml:"dlaPacjenta"`
	DlaMedyka   *DlaMedyka   `xml:"dlaMedyka"`
}

type DlaPacjenta struct {
	MaterialEdukacyjny []MaterialEdukacyjny `xml:"materialEdukacyjny"`
}

type DlaMedyka struct {
	MaterialEdukacyjny []MaterialEdukacyjny `xml:"materialEdukacyjny"`
}

type MaterialEdukacyjny struct {
	Value         string `xml:",chardata"`
	NazwaMaterialu string `xml:"nazwaMaterialu,attr"`
	Material      string `xml:"material,attr"`
}

// ProductInfo holds information about a product and associated package
type ProductInfo struct {
	Product  *ProduktLeczniczy
	Package  *Opakowanie
}
