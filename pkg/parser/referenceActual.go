package parser

import (
	"bufio"
	"cmd/internal/version/actual/dto"
	"fmt"
	xmlparser "github.com/tamerh/xml-stream-parser"
	"os"
	"time"
)

func ParseSax(xmlFilePath *string) (*dto.Container, error) {
	// открываем файл для чтения содержимого
	var container dto.Container
	startTime := time.Now()
	xmlFile, err := os.Open(*xmlFilePath)
	defer func(xmlFile *os.File) {
		err = xmlFile.Close()
		if err != nil {

		}
	}(xmlFile)
	if err != nil {
		return nil, err
	}

	buffer := bufio.NewReaderSize(xmlFile, 65536)
	parser := xmlparser.NewXMLParser(
		buffer,
		`gar:organizator`,
		`gar:operator`,
		`gar:participant`,
		`gar:organizationData`,
	)

	container.ReferenceActual.Participants.Participant = make(map[string]dto.Participant)
	for xmlElement := range parser.Stream() {
		//parse organizators
		if xmlElement.Name == `gar:organizator` {
			container.ReferenceActual.Organizators.Organizator =
				append(
					container.ReferenceActual.Organizators.Organizator,
					dto.Organizator{
						Uid:          xmlElement.Attrs["gar:uid"],
						IedmsId:      xmlElement.Attrs["gar:iedmsId"],
						Title:        xmlElement.Childs[`gar:title`][0].InnerText,
						Organization: xmlElement.Childs[`gar:organization`][0].InnerText,
						Authority:    xmlElement.Childs[`gar:authority`][0].InnerText,
						Phone:        xmlElement.Childs[`gar:phone`][0].InnerText,
						Email:        xmlElement.Childs[`gar:email`][0].InnerText,
					},
				)
		}
		//parse operators
		if xmlElement.Name == `gar:operator` {
			container.ReferenceActual.Operators.Operator =
				append(
					container.ReferenceActual.Operators.Operator,
					dto.Operator{
						Uid:          xmlElement.Attrs["gar:uid"],
						IedmsId:      xmlElement.Attrs["gar:iedmsId"],
						Title:        xmlElement.Childs[`gar:title`][0].InnerText,
						Organization: xmlElement.Childs[`gar:organization`][0].InnerText,
						Authority:    xmlElement.Childs[`gar:authority`][0].InnerText,
						Phone:        xmlElement.Childs[`gar:phone`][0].InnerText,
						Email:        xmlElement.Childs[`gar:email`][0].InnerText,
					},
				)
		}
		//parse participants
		if xmlElement.Name == `gar:participant` {
			participantUuid := xmlElement.Attrs["gar:uid"]
			container.ReferenceActual.Participants.Participant[participantUuid] = dto.Participant{
				IedmsId:      xmlElement.Attrs["gar:iedmsId"],
				Title:        xmlElement.Childs[`gar:title`][0].InnerText,
				Organization: xmlElement.Childs[`gar:organization`][0].InnerText,
				Authority:    xmlElement.Childs[`gar:authority`][0].InnerText,
				Phone:        xmlElement.Childs[`gar:phone`][0].InnerText,
				Email:        xmlElement.Childs[`gar:email`][0].InnerText,
				CommunicationService: dto.CommunicationService{
					OperatorUid: xmlElement.Childs[`gar:communicationService`][0].Childs[`gar:operatorUid`][0].InnerText,
					IsActive:    xmlElement.Childs[`gar:communicationService`][0].Childs[`gar:isActive`][0].InnerText,
					IsSecure:    xmlElement.Childs[`gar:communicationService`][0].Childs[`gar:isSecure`][0].InnerText,
				},
			}
		}
	}

	err = ParseOrganizationData(xmlFilePath, &container)
	if err != nil {
		return nil, err
	}

	endTime := time.Now()
	fmt.Printf("Чтение и парсинг файла заняло %v\n", endTime.Sub(startTime))

	return &container, nil
}

func ParseOrganizationData(xmlFilePath *string, container *dto.Container) (err error) {
	xmlFile, err := os.Open(*xmlFilePath)
	defer func(xmlFile *os.File) {
		err = xmlFile.Close()
		if err != nil {

		}
	}(xmlFile)
	if err != nil {
		return err
	}
	buffer := bufio.NewReaderSize(xmlFile, 65536)
	parser := xmlparser.NewXMLParser(
		buffer,
		`gar:organizationData`,
	)
	//parse organizations
	for xmlElement := range parser.Stream() {
		if xmlElement.Name == `gar:organizationData` {
			var title, address, phone, email, website string
			var classifications []dto.Classification
			var departments []dto.Department
			var persons []dto.Person
			for organizationDataChild := range xmlElement.Childs {
				switch organizationDataChild {
				case "gar:organization":
					for organizationChild := range xmlElement.Childs[organizationDataChild][0].Childs {
						value := xmlElement.Childs[organizationDataChild][0].Childs[organizationChild][0].InnerText
						switch organizationChild {
						case "gar:title":
							title = value
						case "gar:address":
							address = value
						case "gar:phone":
							phone = value
						case "gar:email":
							email = value
						case "gar:website":
							website = value
						}
					}
				case "gar:attestations":
					for i := 0; i < len(xmlElement.Childs[organizationDataChild][0].Childs["gar:classification"]); i++ {
						classifications = append(classifications, dto.Classification{
							Text: xmlElement.Childs[organizationDataChild][0].Childs["gar:classification"][i].InnerText,
							ID:   xmlElement.Childs[organizationDataChild][0].Childs["gar:classification"][i].Attrs["gar:id"],
						})
					}
				case "gar:departments":
					for i := 0; i < len(xmlElement.Childs[organizationDataChild][0].Childs["gar:department"]); i++ {
						departments = append(departments, dto.Department{
							Text: xmlElement.Childs[organizationDataChild][0].Childs["gar:department"][i].InnerText,
							ID:   xmlElement.Childs[organizationDataChild][0].Childs["gar:department"][i].Attrs["gar:id"],
						})
					}
				case "gar:persons":
					for i := 0; i < len(xmlElement.Childs[organizationDataChild][0].Childs["gar:person"]); i++ {
						persons = append(persons, dto.Person{
							ID:           xmlElement.Childs[organizationDataChild][0].Childs["gar:person"][i].Attrs["gar:id"],
							DepartmentId: xmlElement.Childs[organizationDataChild][0].Childs["gar:person"][i].Attrs["gar:departmentId"],
							Post:         xmlElement.Childs[organizationDataChild][0].Childs["gar:person"][i].Childs["gar:post"][0].InnerText,
							Name:         xmlElement.Childs[organizationDataChild][0].Childs["gar:person"][i].Childs["gar:name"][0].InnerText,
							Phone:        xmlElement.Childs[organizationDataChild][0].Childs["gar:person"][i].Childs["gar:phone"][0].InnerText,
							Email:        xmlElement.Childs[organizationDataChild][0].Childs["gar:person"][i].Childs["gar:email"][0].InnerText,
						})
					}
				}
			}

			participantUuid := xmlElement.Attrs["gar:participantUid"]
			entry, ok := container.ReferenceActual.Participants.Participant[participantUuid]
			if ok {
				entry.OrganizationData = dto.OrganizationData{
					ParticipantUid: participantUuid,
					Organization: dto.Organization{
						Title:   title,
						Address: address,
						Phone:   phone,
						Email:   email,
						Website: website,
					},
					Attestations: dto.Attestations{
						Classification: classifications,
					},
					Departments: dto.Departments{
						Department: departments,
					},
					Persons: dto.Persons{
						Person: persons,
					},
				}
				container.ReferenceActual.Participants.Participant[participantUuid] = entry
			}
		}
	}

	return nil
}
