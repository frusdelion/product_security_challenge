package repositories

import (
	"github.com/boltdb/bolt"
	"github.com/frusdelion/zendesk-product_security_challenge/models"
	"github.com/jinzhu/gorm"
	"github.com/klauspost/password"
	"github.com/klauspost/password/drivers/boltpw"
	"github.com/snwfdhmp/errlog"
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	ListUsers() ([]*models.User, error)
	GetUser(userID uint) (*models.User, error)
	UpdateUser(user *models.User, diffUpdate models.User) (*models.User, error)
	DeleteUser(user *models.User) error

	FindUserByUsername(username string) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	FindUserByUserID(userID uint) (*models.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	// Password DB
	pwdb, err := bolt.Open("./common-passwords.db", 0666, nil)
	if errlog.Debug(err) {
		panic(err)
	}

	chk, err := boltpw.New(pwdb, "commonpwd")
	if errlog.Debug(err) {
		panic(err)
	}

	return &userRepository{db: db, pwdb: chk}
}

type userRepository struct {
	db   *gorm.DB
	pwdb *boltpw.BoltDB
}

func (u userRepository) FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	return user, u.db.Find(user, models.User{Email: email}).Error
}

func (u userRepository) FindUserByUserID(userID uint) (*models.User, error) {
	user := &models.User{}

	return user, u.db.First(user, userID).Error
}

func (u userRepository) FindUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	return user, u.db.Where(models.User{Username: username}).First(user).Error
}

func (u userRepository) UpdateUser(user *models.User, diffUpdate models.User) (*models.User, error) {
	user, err := u.GetUser(user.ID)
	if errlog.Debug(err) {
		return nil, err
	}

	// Check if password is set
	if diffUpdate.Password != "" {
		// Check if the password is too simple
		if err := u.checkPasswordWeak(diffUpdate.Password); errlog.Debug(err) {
			return nil, err
		}
	}

	return user, u.db.Model(user).Updates(diffUpdate).Error
}

func (u userRepository) ListUsers() ([]*models.User, error) {
	var users []*models.User
	return users, u.db.Find(users).Error
}

func (u userRepository) CreateUser(user *models.User) (*models.User, error) {

	if user.Password != "" {
		if err := u.checkPasswordWeak(user.Password); errlog.Debug(err) {
			return nil, err
		}
	}

	return user, u.db.Create(user).Error
}

func (u userRepository) checkPasswordWeak(pwd string) error {
	if err := password.Check(pwd, u.pwdb, nil); errlog.Debug(err) {
		return err
	}

	return nil
}

func (u userRepository) GetUser(userID uint) (*models.User, error) {
	user := &models.User{}
	return user, u.db.First(user, userID).Error
}

func (u userRepository) DeleteUser(user *models.User) error {
	return u.db.Model(user).Delete(user).Error
}
