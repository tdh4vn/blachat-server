package repo

import (
	"github.com/gocql/gocql"
)

type ContactsRepo interface {
	GetContactsOfUser(primaryUserID gocql.UUID) ([]gocql.UUID, error)
	CreateContactForUser(primaryID gocql.UUID, secondaryID gocql.UUID) error
	DeleteContact(primaryID gocql.UUID, secondaryID gocql.UUID) error
	CreateContactsForUser(primaryID gocql.UUID, secondaryIDs []gocql.UUID) error
	AddUserToContactOfUsers(primaryIds []gocql.UUID, secondaryId gocql.UUID) error

	GetUsersRelated(userId string) ([]gocql.UUID, error)
}

type ContactsRepoImpl struct {
	DbSession *gocql.Session
}

func NewContactsRepo(db *gocql.Session) ContactsRepo {
	repoImpl := &ContactsRepoImpl{
		DbSession: db,
	}
	return ContactsRepo(repoImpl)
}

const getContactsOfUserQuery = "SELECT secondary_user FROM contacts WHERE primary_user = ?"
const getUserHasContact = "SELECT primary_user FROM contacts WHERE secondary_user = ?"

const createContactsQuery = "INSERT INTO contacts(primary_user, secondary_user) VALUES (?, ?)"
const removeContactViaIDQuery = "DELETE FROM contacts WHERE primary_user = ? AND secondary_user = ?"

func (repo *ContactsRepoImpl) GetUsersRelated(userId string) ([]gocql.UUID, error) {

	var contacts []gocql.UUID

	var iter *gocql.Iter

	iter = repo.DbSession.Query(getUserHasContact, userId).Iter()

	var uId *gocql.UUID

	for iter.Scan(uId) {
		secondaryValueID := *uId
		contacts = append(contacts, secondaryValueID)
	}

	return contacts, nil
}

func (repo *ContactsRepoImpl) GetContactsOfUser(primaryUserID gocql.UUID) ([]gocql.UUID, error) {
	var contacts []gocql.UUID

	var iter *gocql.Iter

	iter = repo.DbSession.Query(getContactsOfUserQuery, primaryUserID).Iter()

	var uId *gocql.UUID

	for iter.Scan(uId) {
		secondaryValueID := *uId
		contacts = append(contacts, secondaryValueID)
	}

	return contacts, nil
}

func (repo *ContactsRepoImpl) CreateContactForUser(primaryID gocql.UUID, secondaryID gocql.UUID) error {
	if err := repo.DbSession.Query(createContactsQuery, primaryID, secondaryID).Exec(); err != nil {
		return err
	} else {
		return nil
	}
}

func (repo *ContactsRepoImpl) CreateContactsForUser(primaryID gocql.UUID, secondaryIDs []gocql.UUID) error {
	batch := repo.DbSession.NewBatch(gocql.LoggedBatch)
	for id := range secondaryIDs {
		batch.Query(createContactsQuery, primaryID, id)
	}

	if err := repo.DbSession.ExecuteBatch(batch); err != nil {
		return err
	}

	return nil
}

func (repo *ContactsRepoImpl) DeleteContact(primaryID gocql.UUID, secondaryID gocql.UUID) error {
	if err := repo.DbSession.Query(removeContactViaIDQuery, primaryID, secondaryID).Exec(); err != nil {
		return err
	} else {
		return nil
	}
}

func (repo *ContactsRepoImpl) AddUserToContactOfUsers(primaryIds []gocql.UUID, secondaryId gocql.UUID) error {
	batch := repo.DbSession.NewBatch(gocql.LoggedBatch)
	for id := range primaryIds {
		batch.Query(createContactsQuery, id, secondaryId)
	}

	if err := repo.DbSession.ExecuteBatch(batch); err != nil {
		return err
	}

	return nil
}

