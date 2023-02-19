package repository

import (
	"database/sql"
	"my-rest-api/internal/model"
	"my-rest-api/pkg/store"
)

type UserRepository interface {
	Create(user *model.User) (*model.User, error)
	FindAll() ([]*model.User, error)
	FindById(id int) (*model.User, error)
	Update(user *model.User) error
	Delete(id int) error
	FindByEmail(email string) (*model.User, error)
}

type UserRepositoryImpl struct {
	store *store.Store
}

func NewUserRepository(store *store.Store) UserRepository {
	return &UserRepositoryImpl{
		store: store,
	}
}

func (u *UserRepositoryImpl) Create(user *model.User) (*model.User, error) {

	if err := user.Validate(); err != nil {
		return nil, err
	}

	if us, err := u.FindByEmail(user.Email); err == store.ErrEmailNotUnique {
		return us, err
	}

	if err := user.BeforeCreate(); err != nil {
		return nil, err
	}

	if err := u.store.Db.QueryRow(
		"INSERT INTO \"User\" (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id",
		user.UserName, user.Email, user.Password, user.Role,
	).Scan(&user.Id); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepositoryImpl) FindAll() ([]*model.User, error) {

	rows, err := u.store.Db.Query("SELECT * FROM \"User\"")
	if err != nil {
		return nil, err
	}

	users := make([]*model.User, 0)
	for rows.Next() {
		var user model.User

		//var role string
		err = rows.Scan(&user.Id, &user.UserName, &user.Email, &user.Password, &user.Role)
		//user.Role = enum.ConvertFromString(role)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil

}

func (u *UserRepositoryImpl) FindById(id int) (*model.User, error) {
	query := "SELECT * FROM \"User\" WHERE id= $1;"

	row := u.store.Db.QueryRow(query, id)
	var user model.User

	if err := row.Scan(&user.Id, &user.UserName, &user.Email, &user.Password, &user.Role); err != nil {
		if err != sql.ErrNoRows {
			return nil, err // db error
		} else {
			return nil, nil //not found user
		}
	}

	return &user, nil
}

func (u *UserRepositoryImpl) Update(user *model.User) error {
	query := "UPDATE \"User\" SET username =$1, email = $2, password = $3, role = $4 WHERE id = $5;"

	err := user.BeforeCreate()
	if err != nil {
		return err
	}

	res, err := u.store.Db.Exec(query, user.UserName, user.Email, user.Password, user.Role, user.Id)
	if err != nil {
		return err
	}

	if amount, err := res.RowsAffected(); err != nil {
		return err
	} else if amount == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (u *UserRepositoryImpl) Delete(id int) error {
	query := "DELETE FROM \"User\" WHERE id = $1;"

	res, err := u.store.Db.Exec(query, id)

	if err != nil {
		return err
	}

	if amount, err := res.RowsAffected(); err != nil {
		return err
	} else if amount == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (u *UserRepositoryImpl) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}
	if err := u.store.Db.QueryRow(
		"SELECT * FROM \"User\" WHERE email = $1",
		email,
	).Scan(
		&user.Id,
		&user.UserName,
		&user.Email,
		&user.Password,
		&user.Role,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return user, store.ErrEmailNotUnique

}
