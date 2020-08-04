package users

import (
	"fmt"
	"go-admin/internal/auth"
	"go-admin/internal/utils/db"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("User not found.")

	ErrUserExists = errors.New("User already exist.")

	ErrIncorrectUsernameOrPassword = errors.New("Incorrect username or password.")

	ErrUsernameOrPasswordNil = errors.New("Username or password can not be null.")

	ErrGenerateTokenFailed = errors.New("Generate token failed")
)

func Create(cu CreateUser) error {
	if err := db.Mysql.
		Where("username = ?", cu.Username).
		Find(&User{}).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return ErrUserExists
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cu.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}

	u := User{Username: cu.Username, Password: string(hash), Email: cu.Email, Status: 1}
	ug := UserGroup{UserID: u.ID, GroupID: cu.GroupID}
	ur := UserRole{UserID: u.ID, RoleID: cu.RoleID}

	tx := db.Mysql.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&u).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&ug).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&ur).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	if _, err := auth.Casbin.AddGroupingPolicy(fmt.Sprintf("user::%d", u.ID), fmt.Sprintf("role::%d", cu.RoleID)); err != nil {
		return err
	}

	return nil
}

func Get(userID string) (*UserInfo, error) {
	var u UserInfo

	if err := db.Mysql.
		Model(&User{}).
		Where("id = ?", userID).
		Scan(&u).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, ErrUserNotFound
		}
	}

	return &u, nil
}

func Update(userID string, uu UpdateUser) error {
	if _, err := Get(userID); err != nil {
		return err
	}

	updateUser := make(map[string]interface{})
	if uu.Username != "" {
		updateUser["username"] = uu.Username
	}
	if uu.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(uu.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
		}
		updateUser["password"] = hash
	}
	if uu.Email != "" {
		updateUser["email"] = uu.Email
	}
	if uu.Status != 0 {
		updateUser["status"] = uu.Status
	}

	tx := db.Mysql.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.
		Model(&User{}).
		Where("id = ?", userID).
		Updates(updateUser).Error; err != nil {
		tx.Rollback()
		return err
	}

	if uu.GroupID != 0 {
		ug := UserGroup{GroupID: uu.GroupID}

		if err := tx.
			Model(&UserGroup{}).
			Where("user_id = ?", userID).
			Updates(&ug).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if uu.RoleID != 0 {
		ur := UserRole{RoleID: uu.RoleID}

		if err := tx.
			Model(&UserRole{}).
			Where("user_id = ?", userID).
			Updates(&ur).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	if uu.RoleID != 0 {
		if _, err := auth.Casbin.RemoveGroupingPolicy(fmt.Sprintf("user::%s", userID), fmt.Sprintf("role::%d", uu.RoleID)); err != nil {
			return err
		}

		if _, err := auth.Casbin.AddGroupingPolicy(fmt.Sprintf("user::%s", userID), fmt.Sprintf("role::%d", uu.RoleID)); err != nil {
			return err
		}
	}

	return nil
}

func Delete(userID string) error {
	if _, err := Get(userID); err != nil {
		return err
	}

	u := User{Status: 2}

	tx := db.Mysql.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.
		Where("id = ?", userID).
		Delete(&u).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.
		Where("user_id = ?", userID).
		Delete(&UserGroup{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.
		Where("user_id = ?", userID).
		Delete(&UserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	if _, err := auth.Casbin.RemoveFilteredGroupingPolicy(0, fmt.Sprintf("user::%s", userID)); err != nil {
		return err
	}

	return nil
}

func List() ([]UserInfo, error) {
	var users []UserInfo

	if err := db.Mysql.
		Select("u.id, u.username, u.email, u.status, u.created_at, u.updated_at, g.id AS group_id, g.group_name, r.id AS role_id, r.role_name").
		Table("user_ AS u").
		Joins("LEFT JOIN user_group ug ON ug.user_id = u.id").
		Joins("LEFT JOIN group_ g ON g.id = ug.group_id").
		Joins("LEFT JOIN user_role ur ON ur.user_id = u.id").
		Joins("LEFT JOIN role_ r ON r.id = ur.role_id").
		Where("u.deleted_at IS NULL").
		Order("u.id ASC").
		Scan(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func Token(l Login) (LoginInfo, error) {
	var li LoginInfo
	if len(l.Username) <= 0 || len(l.Password) <= 0 {
		return li, ErrUsernameOrPasswordNil
	}

	var u User
	if err := db.Mysql.
		Where("username = ? ", l.Username).
		Find(&u).Error; err != nil {
		return li, ErrIncorrectUsernameOrPassword
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(l.Password)); err != nil {
		return li, ErrIncorrectUsernameOrPassword
	}

	if u.Status != 1 {
		return li, ErrIncorrectUsernameOrPassword
	}

	token, tokenErr := auth.GenerateToken(u.ID)

	if tokenErr != nil {
		return li, ErrGenerateTokenFailed
	}
	li.Username = l.Username
	li.ID = u.ID
	li.Token = token
	return li, nil
}
