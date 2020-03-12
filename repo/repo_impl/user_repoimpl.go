package repo_impl

import (
	"blachat-server/entities"
	"github.com/gocql/gocql"
	"time"
)

var findByIdQuery = `SELECT id, name, avatar, created_at, updated_at FROM users WHERE id = ?`
var findByIdsQuery = `SELECT id, name, avatar, created_at, updated_at FROM users WHERE id IN ?`
var insertUserQuery = `INSERT INTO users(id, name, avatar, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
var updateUserQuery = `UPDATE users SET name=?, avatar=?, updated_at=? WHERE id = ?`

type UserRepoImpl struct {
	DbSession *gocql.Session
}

func NewUserRepo(db *gocql.Session) *UserRepoImpl {
	return &UserRepoImpl {
		DbSession: db,
	}
}

func (repo *UserRepoImpl) FindByID(id gocql.UUID) (*entities.User, error) {
	user := entities.User {}

	if err := repo.DbSession.Query(findByIdQuery, id).Consistency(gocql.One).Scan(
		&user.ID,
		&user.Name,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepoImpl) FindByIDs(ids []gocql.UUID) ([]*entities.User, error) {
	var users []*entities.User

	iter := repo.DbSession.Query(findByIdsQuery, ids).Iter()

	var id gocql.UUID
	var name string
	var avatar string
	var createdAt time.Time
	var updatedAt time.Time

	for iter.Scan(&id, &name, &avatar, &createdAt, &updatedAt) {
		users = append(users, &entities.User{
			ID:        id,
			Name:      name,
			Avatar:    avatar,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	return users, nil
}

func (repo *UserRepoImpl) Insert(user *entities.User) (*entities.User, error) {
	if newID, err := gocql.RandomUUID(); err != nil {
		return nil, err
	} else {
		user.ID = newID
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

 		if err := repo.DbSession.Query(insertUserQuery,
			user.ID,
			user.Name,
			user.Avatar,
			user.CreatedAt,
			user.UpdatedAt).Exec(); err != nil {
				return nil, err
		}
			return user, nil
	}
}

func (repo *UserRepoImpl) Update(user *entities.User) (*entities.User, error) {
	if err := repo.DbSession.Query(updateUserQuery, user.Name, user.Avatar, time.Now(), user.ID).Exec(); err != nil {
		return nil, err
	}
	user.UpdatedAt = time.Now()
	user.CreatedAt = time.Now()
	return user, nil
}
