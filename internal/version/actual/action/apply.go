package action

import (
	operatorModel "cmd/internal/operator"
	participantModel "cmd/internal/participant"
	"cmd/internal/participant/organization/classification"
	"cmd/internal/version/actual/dto"
	"cmd/pkg/utils/formatter"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Uuids struct {
	Uuid string
}

func ApplyVersion(
	conn *pgxpool.Pool,
	container *dto.Container,
	// удаляем старый ГАС
) (err error) {
	var tx pgx.Tx
	//operators := make([]operatorModel.Operator, 0)
	//var o operatorModel.Operators
	//begin transaction
	tx, err = conn.Begin(context.TODO())
	startTime := time.Now()
	err = DeleteGas(&tx)
	if err != nil {
		_ = tx.Rollback(context.TODO())
		return err
	}
	endTime := time.Now()
	fmt.Printf("Удаление всех записей справочника ГАС выполнен за %v\n", endTime.Sub(startTime))
	//создаем абонентов
	startTime = time.Now()
	err = CreateAbonents(&tx, container)
	if err != nil {
		_ = tx.Rollback(context.TODO())
	}
	endTime = time.Now()
	fmt.Printf("Инсерт записей абонентов выполнен за %v\n", endTime.Sub(startTime))

	//создаем организаторов
	startTime = time.Now()
	err = CreateOrganizators(&tx, &container.ReferenceActual.Organizators)
	if err != nil {
		_ = tx.Rollback(context.TODO())
	}
	endTime = time.Now()
	fmt.Printf("Инсерт записей организаторов выполнен за %v\n", endTime.Sub(startTime))

	//создаем операторов
	startTime = time.Now()
	var operatorCollection *[]operatorModel.Operator
	operatorCollection, err = CreateOperators(&tx, &container.ReferenceActual.Operators)
	if err != nil {
		_ = tx.Rollback(context.TODO())
	}
	endTime = time.Now()
	fmt.Printf("Инсерт записей операторов выполнен за %v\n", endTime.Sub(startTime))

	//создаем участников
	err = CreateParticipants(
		&tx,
		&container.ReferenceActual.Participants,
		operatorCollection)
	if err != nil {
		_ = tx.Rollback(context.TODO())
	}
	err = tx.Commit(context.TODO())
	if err != nil {
		return err
	}

	return nil
}

func CreateAbonents(
	tx *pgx.Tx,
	container *dto.Container,
) (err error) {
	db := *tx
	batch := &pgx.Batch{}
	q := formatter.Query(
		`insert into abonents (
	                        uuid,
	                        abonent_name,
	                        medo_address,
	                        created_at,
	                        updated_at
							)
				values ($1, $2, $3, current_timestamp, current_timestamp)
				`)
	for _, organizator := range container.ReferenceActual.Organizators.Organizator {
		batch.Queue(
			q,
			organizator.Uid,
			organizator.Title,
			organizator.IedmsId,
		)
	}
	for _, operator := range container.ReferenceActual.Operators.Operator {
		batch.Queue(
			q,
			operator.Uid,
			operator.Title,
			operator.IedmsId,
		)
	}
	for uuid, participant := range container.ReferenceActual.Participants.Participant {
		batch.Queue(
			q,
			uuid,
			participant.Title,
			participant.IedmsId,
		)
	}
	br := db.SendBatch(context.TODO(), batch)
	_, err = br.Exec()
	if err != nil {
		return err
	}
	err = br.Close()
	if err != nil {
		return err
	}
	return nil
}

func CreateOrganizators(
	tx *pgx.Tx,
	organizators *dto.Organizators,
) (err error) {
	db := *tx
	batch := &pgx.Batch{}
	q := formatter.Query(
		`INSERT INTO organizators.organizators (
                            uuid,
							organizator_short_name,
							organizator_organization_legal_name,
							organizator_authority_fio,
							organizator_authority_phone,
							organizator_authority_email,
                            created_at,
                            updated_at
							)
				VALUES ($1, $2, $3, $4, $5, $6, current_timestamp, current_timestamp)
				`)
	for _, organizator := range organizators.Organizator {
		batch.Queue(
			q,
			organizator.Uid,
			organizator.Title,
			organizator.Organization,
			organizator.Authority,
			organizator.Phone,
			organizator.Email,
		)
	}
	br := db.SendBatch(context.TODO(), batch)
	_, err = br.Exec()
	if err != nil {
		return err
	}
	err = br.Close()
	if err != nil {
		return err
	}
	return nil
}

func CreateOperators(
	tx *pgx.Tx,
	operators *dto.Operators,
) (collection *[]operatorModel.Operator, err error) {
	var rows pgx.Rows
	var o operatorModel.Operator
	db := *tx
	batch := &pgx.Batch{}
	q := formatter.Query(
		`INSERT INTO operators.operators (
                            uuid,
							operator_short_name,
							operator_organization_legal_name,
							operator_authority_fio,
							operator_authority_phone,
							operator_authority_email,
                            created_at,
                            updated_at
							)
				VALUES ($1, $2, $3, $4, $5, $6, current_timestamp, current_timestamp)
				RETURNING id, uuid
				`)
	for _, operator := range operators.Operator {
		batch.Queue(
			q,
			operator.Uid,
			operator.Title,
			operator.Organization,
			operator.Authority,
			operator.Phone,
			operator.Email,
		)
	}
	br := db.SendBatch(context.TODO(), batch)
	_, err = br.Exec()
	if err != nil {
		return nil, err
	}
	err = br.Close()
	if err != nil {
		return nil, err
	}
	batch = &pgx.Batch{}
	batch.Queue(`SELECT id, uuid from operators.operators`)
	br = db.SendBatch(context.TODO(), batch)
	rows, err = br.Query()
	if err != nil {
		return nil, err
	}
	oCollect := make([]operatorModel.Operator, 0)
	for rows.Next() {
		err = rows.Scan(&o.ID, &o.Uuid)
		if err != nil {
			return nil, err
		}
		oCollect = append(oCollect, o)
	}
	err = br.Close()
	if err != nil {
		return nil, err
	}
	return &oCollect, nil
}

func CreateParticipants(
	tx *pgx.Tx,
	participants *dto.Participants,
	operatorCollection *[]operatorModel.Operator,
) (err error) {
	startTime := time.Now()
	q := formatter.Query(
		`INSERT INTO participants.participants (
                           uuid,
							operator_id,
							participant_short_name,
							participant_authority_fio,
							participant_authority_phone,
							participant_authority_email,
                           participant_can_dsp,
                           participant_is_active,
                           created_at,
                           updated_at
							)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, current_timestamp, current_timestamp)
				RETURNING id
				`)
	var rows pgx.Rows
	db := *tx
	batch := &pgx.Batch{}
	for uuid, participant := range participants.Participant {
		var operatorId *int
		operatorId, err = findOperatorByUuid(&participant.CommunicationService.OperatorUid, operatorCollection)
		if err != nil {
			return err
		}
		participantIsSecure := true
		if participant.CommunicationService.IsSecure == `false` {
			participantIsSecure = false
		}
		participantIsActive := true
		if participant.CommunicationService.IsActive == `false` {
			participantIsActive = false
		}
		batch.Queue(
			q,
			uuid,
			*operatorId,
			participant.Title,
			participant.Authority,
			participant.Phone,
			participant.Email,
			participantIsSecure,
			participantIsActive,
		)
	}
	br := db.SendBatch(context.TODO(), batch)
	_, err = br.Exec()
	if err != nil {
		return err
	}
	err = br.Close()
	if err != nil {
		return err
	}
	endTime := time.Now()
	fmt.Printf("Инсерт записей участников выполнен за %v\n", endTime.Sub(startTime))
	batch = &pgx.Batch{}
	batch.Queue(`SELECT id, uuid from participants.participants`)
	br = db.SendBatch(context.TODO(), batch)
	rows, err = br.Query()
	if err != nil {
		return err
	}
	var pId int
	var pUuid string
	for rows.Next() {
		err = rows.Scan(&pId, &pUuid)
		entry, ok := participants.Participant[pUuid]
		if ok {
			entry.DbID = pId
			participants.Participant[pUuid] = entry
		}
		if err != nil {
			return err
		}
	}
	err = br.Close()
	if err != nil {
		return err
	}

	err = CreateOrganizations(tx, participants)
	if err != nil {
		return err
	}
	return nil
}

func CreateOrganizations(
	tx *pgx.Tx,
	participants *dto.Participants,
) (err error) {
	startTime := time.Now()
	q := formatter.Query(
		`INSERT INTO participants.organizations (
	                            organization_ogrn,
								organization_email,
								organization_legal_address,
								organization_legal_name,
								organization_phone,
								organization_web_address,
	                            participant_id,
	                            created_at,
	                            updated_at
								)
					VALUES ($1, $2, $3, $4, $5, $6, $7, current_timestamp, current_timestamp)
					`)
	batch := &pgx.Batch{}
	db := *tx
	//var organizationLegalName *string
	//var participantId *int
	//var classifications *[]classification.Classification
	//classifications, err = classificationsRepository.NewRepository(db).GetAll(context.TODO())
	if err != nil {
		return err
	}
	for _, participant := range participants.Participant {
		batch.Queue(
			q,
			participant.OrganizationData.Organization.OrgRegNum,
			participant.OrganizationData.Organization.Email,
			participant.OrganizationData.Organization.Address,
			participant.Organization,
			participant.OrganizationData.Organization.Phone,
			participant.OrganizationData.Organization.Website,
			participant.DbID,
		)
	}
	br := db.SendBatch(context.TODO(), batch)
	_, err = br.Exec()
	if err != nil {
		return err
	}
	endTime := time.Now()
	err = br.Close()
	if err != nil {
		return err
	}
	fmt.Printf("Инсерт записей организаций выполнен за %v\n", endTime.Sub(startTime))
	//batch = &pgx.Batch{}
	//batch.Queue(`SELECT id from participants.organizations`)
	//br = db.SendBatch(context.TODO(), batch)
	//var o organizationModel.Organization
	//var rows pgx.Rows
	//rows, err = br.Query()
	//if err != nil {
	//	return err
	//}
	//oCollect := make([]organizationModel.Organization, 0)
	//i := 0
	//batch = &pgx.Batch{}
	//q = formatter.Query(
	//	`INSERT INTO participants.organization_classification (
	//                            organization_id,
	//							classification_id,
	//                            created_at,
	//                            updated_at
	//							)
	//				VALUES ($1, $2, current_timestamp, current_timestamp)
	//				`)
	//var cID *int
	//for rows.Next() {
	//	err = rows.Scan(&o.ID)
	//	if err != nil {
	//		return err
	//	}
	//	if organizations.OrganizationData[i].Attestations.Classification != nil {
	//		for _, c := range organizations.OrganizationData[i].Attestations.Classification {
	//			cID, err = findClassificationId(&c.ID, classifications)
	//			batch.Queue(
	//				q,
	//				o.ID,
	//				*cID,
	//			)
	//		}
	//	}
	//	oCollect = append(oCollect, o)
	//	i++
	//}
	//err = br.Close()
	//if err != nil {
	//	return err
	//}
	//br = db.SendBatch(context.TODO(), batch)
	//rows, err = br.Query()
	//if err != nil {
	//	return err
	//}
	//err = br.Close()
	//if err != nil {
	//	return err
	//}
	//сохраняем грифы доступа
	/*if organization.Attestations.Classification != nil {
		var classificationData *classification.Classification
		qClassifications := formatter.Query(
			`INSERT INTO participants.organization_classification (
	                            organization_id,
								classification_id,
	                            created_at,
	                            updated_at
								)
					VALUES ($1, $2, current_timestamp, current_timestamp)
					`)
		for _, classificationNode := range organization.Attestations.Classification {
			classificationData, err = classificationRepository.NewRepository(*tx).
				FindByClassificationId(context.TODO(), classificationNode.ID)
			if err != nil {
				return err
			}
			if classificationData != nil {
				_, err = db.Exec(
					context.TODO(),
					qClassifications,
					OrganizationId,
					classificationData.ID,
				)
				if err != nil {
					return err
				}
			}
		}
	}*/
	//если у организации имеются подразделения, запускаем их сохранение
	/*if organization.Departments.Department != nil {
		err = CreateDepartments(tx, &OrganizationId, organization)
		if err != nil {
			return err
		}
	}*/
	return nil
}

/*func CreateDepartments(
	tx *pgx.Tx,
	organizationId *int,
	organization *dto.OrganizationData,
) (err error) {
	db := *tx
	q := formatter.Query(
		`INSERT INTO participants.departments (
	                            organization_id,
								department_legal_name,
	                            created_at,
	                            updated_at
								)
					VALUES ($1, $2, current_timestamp, current_timestamp)
					RETURNING id
					`)
	for _, department := range organization.Departments.Department {
		var DepartmentId int
		err = db.QueryRow(
			context.TODO(),
			q,
			organizationId,
			department.Text,
		).Scan(&DepartmentId)
		if err != nil {
			return err
		}
		persons := findPersonsByDepartmentId(&department.ID, organization)
		if persons != nil {
			err = CreatePersons(tx, &DepartmentId, &organization.Persons)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CreatePersons(
	tx *pgx.Tx,
	departmentId *int,
	persons *dto.Persons,
) (err error) {
	db := *tx
	q := formatter.Query(
		`INSERT INTO participants.persons (
	                            department_id,
								person_post,
                                person_name,
                                person_phone,
                                person_email,
	                            created_at,
	                            updated_at
								)
					VALUES ($1, $2, $3, $4, $5, current_timestamp, current_timestamp)
					`)
	for _, person := range persons.Person {
		_, err = db.Exec(
			context.TODO(),
			q,
			departmentId,
			person.Post,
			person.Name,
			person.Phone,
			person.Email,
		)
		if err != nil {
			return err
		}
	}
	return nil
}*/

func DeleteGas(tx *pgx.Tx) error {
	q := formatter.Query(`DELETE FROM public.abonents
			WHERE abonents.uuid is not null`)
	db := *tx
	_, err := db.Exec(context.TODO(), q)
	if err != nil {
		return err
	}
	return nil
}

//найти организацию для участника по аттрибуту uuid
func findOrganizationByParticipantUid(
	uuid string,
	organizations *dto.OrganizationsData,
) (*dto.OrganizationData, error) {
	var organization dto.OrganizationData
	for _, organization = range organizations.OrganizationData {
		if organization.ParticipantUid == uuid {
			return &organization, nil
		}
	}
	return &organization, errors.New(
		fmt.Sprintf(
			`Организация с uuid <%u> не была найдена.`,
			uuid,
		),
	)
}

//найти сотрудников по аттрибуту departmentId
func findPersonsByDepartmentId(
	departmentId *string,
	organization *dto.OrganizationData,
) (result []dto.Person) {
	for _, person := range organization.Persons.Person {
		if *departmentId == person.DepartmentId {
			result = append(result, person)
		}
	}
	return result
}

//найти id оператора по его uuid
func findOperatorByUuid(
	uuid *string,
	operatorsCollection *[]operatorModel.Operator,
) (
	id *int,
	err error,
) {
	for _, operator := range *operatorsCollection {
		if operator.Uuid == *uuid {
			return &operator.ID, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("оператор с uuid <%u> не найден", *uuid))
}

//найти участника по его uuid
func findParticipantByUuid(
	uuid *string,
	participantsCollection *[]participantModel.Participant,
) (
	id *int,
	organizationName *string,
	err error,
) {
	for _, participant := range *participantsCollection {
		if participant.Uuid == *uuid {
			return &participant.ID, &participant.Organization, nil
		}
	}
	return nil, nil, errors.New(fmt.Sprintf("участник с uuid <%u> не найден", *uuid))
}

//найти id грифа по classification_id
func findClassificationId(
	ClassificationId *string,
	collection *[]classification.Classification,
) (
	id *int,
	err error,
) {
	for _, c := range *collection {
		if c.ClassificationId == *ClassificationId {
			return &c.ID, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("гриф доступа с classification_id <%c> не найден", *ClassificationId))
}
