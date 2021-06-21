package users

import (
	"admin-go/internal/auth"
	"admin-go/internal/utils/db"
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("User not found.")

	ErrUserExists = errors.New("User already exist.")

	ErrIncorrectUsernameOrPassword = errors.New("Incorrect username or password.")

	ErrUsernameOrPasswordNil = errors.New("Username or password can not be null.")

	ErrGenerateTokenFailed = errors.New("Generate token failed")
)

func Create(cu CreateUser) error {
	u := User{}

	if err := db.Mysql.
		Where("username = ?", cu.Username).
		Find(&u).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	if u.ID != 0 {
		return ErrUserExists
	}

	hash, pwErr := bcrypt.GenerateFromPassword([]byte(cu.Password), bcrypt.DefaultCost)
	if pwErr != nil {
		return pwErr
	}

	nu := User{Username: cu.Username, Password: string(hash), Email: cu.Email, Status: 1}

	err := db.Mysql.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&nu).Error; err != nil {
			return err
		}

		ug := UserGroup{UserID: nu.ID, GroupID: cu.GroupID}
		ur := UserRole{UserID: nu.ID, RoleID: cu.RoleID}
		if err := tx.Create(&ug).Error; err != nil {
			return err
		}
		if err := tx.Create(&ur).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	if _, casbinErr := auth.Casbin.AddGroupingPolicy(fmt.Sprintf("user::%d", u.ID), fmt.Sprintf("role::%d", cu.RoleID)); casbinErr != nil {
		return casbinErr
	}

	return nil
}

func Get(userID uint64) (user User, err error) {
	if err = db.Mysql.Take(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, ErrUserNotFound
		} else {
			return user, err
		}
	}
	return user, nil
}

func GetInfo(userID uint64) (userInfo UserInfo, err error) {
	if err = db.Mysql.
		Model(&User{}).
		Where("id = ?", userID).
		Take(&userInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userInfo, ErrUserNotFound
		} else {
			return userInfo, err
		}
	}
	return userInfo, err
}

func Update(userID uint64, uu UpdateUser) error {

	or := UserRole{}
	if err := db.Mysql.
		Where("user_id = ?", userID).
		Take(&or).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		} else {
			return err
		}
	}

	updateUser := make(map[string]interface{})
	if uu.Username != "" {
		updateUser["username"] = uu.Username
	}
	if uu.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		updateUser["password"] = string(hash)
	}
	if uu.Email != "" {
		updateUser["email"] = uu.Email
	}
	if uu.Status != 0 {
		updateUser["status"] = uu.Status
	}

	err := db.Mysql.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Model(&User{}).
			Where("id = ?", userID).
			Updates(&updateUser).Error; err != nil {
			return err
		}

		if uu.GroupID != 0 {
			ug := UserGroup{GroupID: uu.GroupID}
			if err := tx.
				Model(&UserGroup{}).
				Where("id = ?", userID).
				Updates(&ug).Error; err != nil {
				return err
			}
		}
		if uu.RoleID != 0 {
			ur := UserRole{RoleID: uu.RoleID}
			if err := tx.
				Model(&UserRole{}).
				Where("id = ?", userID).
				Updates(&ur).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	if uu.RoleID != 0 {
		if _, casbinErr := auth.Casbin.RemoveGroupingPolicy(fmt.Sprintf("user::%d", userID), fmt.Sprintf("role::%d", or.RoleID)); casbinErr != nil {
			return casbinErr
		}

		if _, casbinErr := auth.Casbin.AddGroupingPolicy(fmt.Sprintf("user::%d", userID), fmt.Sprintf("role::%d", uu.RoleID)); casbinErr != nil {
			return casbinErr
		}
	}

	return nil
}

func Delete(userID uint64) error {
	if _, err := Get(userID); err != nil {
		return err
	}

	u := User{Status: 2}

	err := db.Mysql.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Where("id = ?", userID).
			Delete(&u).Error; err != nil {
			return err
		}

		if err := tx.
			Where("id = ?", userID).
			Delete(&UserGroup{}).Error; err != nil {
			return err
		}

		if err := tx.
			Where("id = ?", userID).
			Delete(&UserRole{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	if _, casbinErr := auth.Casbin.RemoveFilteredGroupingPolicy(0, fmt.Sprintf("user::%d", userID)); casbinErr != nil {
		return casbinErr
	}

	return nil
}

func List() (users []UserInfo, err error) {
	const q = `
		SELECT u.id, u.username, u.email, u.status, u.created_at, u.updated_at, g.id AS group_id, g.group_name, r.id AS role_id, r.role_name
		FROM users AS u
         LEFT JOIN user_groups ug ON ug.user_id = u.id
         LEFT JOIN groups g ON g.id = ug.group_id
         LEFT JOIN user_roles ur ON ur.user_id = u.id
         LEFT JOIN roles r ON r.id = ur.role_id
		WHERE u.deleted_at IS NULL
		ORDER BY u.id`

	err = db.Mysql.
		Raw(q).
		Scan(&users).Error
	return users, err
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
