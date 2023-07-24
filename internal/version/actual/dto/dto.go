package dto

import "encoding/xml"

type Container struct {
	XMLName xml.Name `xml:"container"`
	Text    string   `xml:",chardata"`
	Gar     string   `xml:"gar,attr"`
	Version string   `xml:"version,attr"`
	Header  struct {
		Text        string `xml:",chardata"`
		Uid         string `xml:"uid"`
		Created     string `xml:"created"`
		PreviewFile string `xml:"previewFile"`
	} `xml:"header"`
	ReferenceActual struct {
		Text           string       `xml:",chardata"`
		ExtractionDate string       `xml:"extractionDate"`
		DataVersion    string       `xml:"dataVersion"`
		Organizators   Organizators `xml:"organizators"`
		Operators      Operators    `xml:"operators"`
		Participants   Participants `xml:"participants"`
	} `xml:"referenceActual"`
}

type Organizators struct {
	Text        string        `xml:",chardata"`
	Organizator []Organizator `xml:"organizator"`
}

type Organizator struct {
	Text         string `xml:",chardata"`
	Uid          string `xml:"uid,attr"`
	IedmsId      string `xml:"iedmsId,attr"`
	Title        string `xml:"title"`
	Organization string `xml:"organization"`
	Authority    string `xml:"authority"`
	Phone        string `xml:"phone"`
	Email        string `xml:"email"`
}

type Operators struct {
	Text     string     `xml:",chardata"`
	Operator []Operator `xml:"operator"`
}

type Operator struct {
	Text         string `xml:",chardata"`
	Uid          string `xml:"uid,attr"`
	IedmsId      string `xml:"iedmsId,attr"`
	Title        string `xml:"title"`
	Organization string `xml:"organization"`
	Authority    string `xml:"authority"`
	Phone        string `xml:"phone"`
	Email        string `xml:"email"`
}

type Participants struct {
	Text        string `xml:",chardata"`
	Participant map[string]Participant
}

type Participant struct {
	DbID                 int
	Text                 string               `xml:",chardata"`
	Uid                  string               `xml:"uid,attr"`
	IedmsId              string               `xml:"iedmsId,attr"`
	Title                string               `xml:"title"`
	Organization         string               `xml:"organization"`
	Authority            string               `xml:"authority"`
	Phone                string               `xml:"phone"`
	Email                string               `xml:"email"`
	CommunicationService CommunicationService `xml:"communicationService"`
	OrganizationData     OrganizationData     `xml:"organizationData"`
}

type CommunicationService struct {
	Text        string `xml:",chardata"`
	OperatorUid string `xml:"operatorUid"`
	IsActive    string `xml:"isActive"`
	IsSecure    string `xml:"isSecure"`
}

type OrganizationsData struct {
	Text             string             `xml:",chardata"`
	OrganizationData []OrganizationData `xml:"organizationData"`
}

type OrganizationData struct {
	Text           string       `xml:",chardata"`
	ParticipantUid string       `xml:"participantUid,attr"`
	Organization   Organization `xml:"organization"`
	Attestations   Attestations `xml:"attestations"`
	Departments    Departments  `xml:"departments"`
	Persons        Persons      `xml:"persons"`
}

type Organization struct {
	dbID      int
	Text      string `xml:",chardata"`
	OrgRegNum string `xml:"orgRegNum,attr"`
	Title     string `xml:"title"`
	Address   string `xml:"address"`
	Phone     string `xml:"phone"`
	Email     string `xml:"email"`
	Website   string `xml:"website"`
}

type Attestations struct {
	Text           string           `xml:",chardata"`
	Classification []Classification `xml:"classification"`
}

type Classification struct {
	Text string `xml:",chardata"`
	ID   string `xml:"id,attr"`
}

type Departments struct {
	Text       string       `xml:",chardata"`
	Department []Department `xml:"department"`
}

type Department struct {
	dbID int
	Text string `xml:",chardata"`
	ID   string `xml:"id,attr"`
}

type Persons struct {
	Text   string   `xml:",chardata"`
	Person []Person `xml:"person"`
}

type Person struct {
	Text         string `xml:",chardata"`
	ID           string `xml:"id,attr"`
	DepartmentId string `xml:"departmentId,attr"`
	Post         string `xml:"post"`
	Name         string `xml:"name"`
	Phone        string `xml:"phone"`
	Email        string `xml:"email"`
}

func NewActualDto() *Container {
	var container Container
	return &container
}
