package groups

import (
	"fmt"
	"go-admin/internal/auth"
	"go-admin/internal/utils/db"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var (
	ErrGroupNotFound = errors.New("Group not found.")

	ErrGroupExists = errors.New("Group already exist.")

	ErrGroupHasUser = errors.New("Group has user, delete user first.")
)

func Create(cg CreateGroup) error {
	if err := db.Mysql.
		Where("group_name = ?", cg.GroupName).
		Find(&Group{}).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return ErrGroupExists
		}
	}

	g := Group{GroupName: cg.GroupName}
	gr := GroupRole{GroupID: g.ID, RoleID: cg.RoleID}

	tx := db.Mysql.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&g).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&gr).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	if _, err := auth.Casbin.AddGroupingPolicy(fmt.Sprintf("group::%d", g.ID), fmt.Sprintf("role::%d", gr.RoleID)); err != nil {
		return err
	}

	return nil
}

func Get(groupID string) (*GetGroupInfo, error) {
	var g GetGroupInfo

	if err := db.Mysql.
		Model(&Group{}).
		Where("id = ?", groupID).
		Scan(&g).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, ErrGroupNotFound
		}
	}

	return &g, nil
}

func Update(groupID string, ug UpdateGroup) error {
	if _, err := Get(groupID); err != nil {
		return err
	}

	updateGroup := make(map[string]interface{})

	if ug.GroupName != "" {
		updateGroup["group_name"] = ug.GroupName
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
		Model(&Group{}).
		Where("id = ?", groupID).
		Updates(updateGroup).Error; err != nil {
		tx.Rollback()
		return err
	}

	if ug.RoleID != 0 {
		gr := GroupRole{RoleID: ug.RoleID}

		if err := tx.
			Model(&GroupRole{}).
			Where("group_id = ?", groupID).
			Updates(&gr).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	if ug.RoleID != 0 {
		if _, err := auth.Casbin.RemoveGroupingPolicy(fmt.Sprintf("group::%s", groupID), fmt.Sprintf("role::%d", ug.RoleID)); err != nil {
			return err
		}

		if _, err := auth.Casbin.AddGroupingPolicy(fmt.Sprintf("group::%s", groupID), fmt.Sprintf("role::%d", ug.RoleID)); err != nil {
			return err
		}
	}

	return nil
}

func Delete(groupID string) error {
	if _, err := Get(groupID); err != nil {
		return err
	}

	u, err := GetUserForGroup(groupID)
	if err != nil {
		return err
	}

	if len(u) > 0 {
		return ErrGroupHasUser
	}

	if err := db.Mysql.
		Where("id = ?", groupID).
		Delete(&Group{}).Error; err != nil {
		return err
	}

	if _, err := auth.Casbin.RemoveFilteredGroupingPolicy(0, fmt.Sprintf("group::%s", groupID)); err != nil {
		return err
	}

	return nil
}

func List() ([]GetGroupInfo, error) {
	var groups []GetGroupInfo

	if err := db.Mysql.
		Select("g.id, g.group_name, g.created_at, g.updated_at, r.id AS role_id, r.role_name").
		Table("group_ AS g").
		Joins("LEFT JOIN group_role AS gr ON gr.group_id = g.id").
		Joins("LEFT JOIN role_ AS r ON r.id = gr.role_id").
		Where("g.deleted_at IS NULL").
		Where("gr.deleted_at IS NULL").
		Where("r.deleted_at IS NULL").
		Order("g.id ASC").
		Scan(&groups).Error; err != nil {
		return nil, err
	}

	return groups, nil
}

func GetUserForGroup(groupID string) ([]User, error) {
	var u []User

	if err := db.Mysql.
		Select("u.id, u.username").
		Table("user_group AS ug").
		Joins("LEFT JOIN user_ AS u ON ug.user_id = u.id").
		Where("ug.deleted_at IS NULL").
		Where("u.deleted_at IS NULL").
		Where("group_id =?", groupID).
		Order("u.id ASC").
		Scan(&u).Error; err != nil {
		return nil, err
	}

	return u, nil
}
