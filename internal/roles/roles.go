package roles

import (
	"fmt"
	"go-admin/internal/auth"
	"go-admin/internal/utils"
	"go-admin/internal/utils/db"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var (
	ErrRoleNotFound = errors.New("Role not found.")

	ErrRoleExists = errors.New("Role already exist.")
)

func Create(cr CreateRole) error {
	if err := db.Mysql.
		Where("role_name = ?", cr.RoleName).
		Find(&Role{}).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return ErrRoleExists
		}
	}

	r := Role{RoleName: cr.RoleName}

	err := db.Mysql.Create(&r).Error
	return err
}

func Get(roleID string) (*GetRoleInfo, error) {
	var r GetRoleInfo

	if err := db.Mysql.
		Model(&Role{}).
		Where("id = ?", roleID).
		Scan(&r).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, ErrRoleNotFound
		}
	}

	return &r, nil
}

func Update(roleID string, ur UpdateRole) error {
	if _, err := Get(roleID); err != nil {
		return err
	}

	updateRole := make(map[string]interface{})
	if ur.RoleName != "" {
		updateRole["role_name"] = ur.RoleName
	}

	err := db.Mysql.
		Model(Role{}).
		Where("id = ?", roleID).
		Updates(updateRole).Error
	return err
}

func Delete(roleID string) error {
	if _, err := Get(roleID); err != nil {
		return err
	}

	err := db.Mysql.
		Where("id = ?", roleID).
		Delete(&Role{}).Error
	if _, err := auth.Casbin.RemoveFilteredGroupingPolicy(0, fmt.Sprintf("role::%s", roleID)); err != nil {
		return err
	}
	if _, err := auth.Casbin.RemoveFilteredGroupingPolicy(1, fmt.Sprintf("role::%s", roleID)); err != nil {
		return err
	}

	return err
}

func List() ([]GetRoleInfo, error) {
	var roles []GetRoleInfo

	if err := db.Mysql.
		Model(Role{}).
		Scan(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func GetRoleMenus(roleID string) ([]RoleMenu, error) {
	var rm []RoleMenu

	if err := db.Mysql.
		Model(&RoleMenu{}).
		Where("role_id = ?", roleID).
		Scan(&rm).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, ErrRoleNotFound
		}
	}

	return rm, nil
}

func UpdateRoleMenu(roleID string, new []uint) error {
	rm, err := GetRoleMenus(roleID)
	old := make([]uint, 0)
	for _, menu := range rm {
		old = append(old, menu.MenuID)
	}

	delRoleMenu := utils.Difference(old, new)
	addRoleMenu := utils.Difference(new, old)
	for _, id := range delRoleMenu {
		db.Mysql.
			Where("role_id = ?", roleID).
			Where("menu_id = ?", id).
			Delete(&RoleMenu{})
		if _, err := auth.Casbin.RemoveGroupingPolicy(fmt.Sprintf("role::%s", roleID), fmt.Sprintf("action::%d", id)); err != nil {
			return err
		}
	}
	ridInt, _ := strconv.Atoi(roleID)
	rid := uint(ridInt)
	for _, mid := range addRoleMenu {
		rm := RoleMenu{RoleID: rid, MenuID: mid}
		db.Mysql.Create(&rm)
		if _, err := auth.Casbin.AddGroupingPolicy(fmt.Sprintf("role::%d", rid), fmt.Sprintf("action::%d", mid)); err != nil {
			return err
		}
	}
	return err
}
