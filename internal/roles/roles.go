package roles

import (
	"admin-go/internal/auth"
	"admin-go/internal/groups"
	"admin-go/internal/users"
	"admin-go/internal/utils"
	"admin-go/internal/utils/db"
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	ErrRoleNotFound = errors.New("角色不存在")

	ErrRoleExists = errors.New("角色已存在")

	ErrRoleBind = errors.New("存在未删除绑定关系")
)

func Create(cr CreateRole) error {
	r := Role{}

	if err := db.Mysql.
		Where("role_name = ?", cr.RoleName).
		Find(&r).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	if r.ID != 0 {
		return ErrRoleExists
	}

	nr := Role{RoleName: cr.RoleName}

	err := db.Mysql.Create(&nr).Error
	return err
}

func Get(roleID uint64) (role Role, err error) {
	if err = db.Mysql.Take(&role, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return role, ErrRoleNotFound
		} else {
			return role, err
		}
	}
	return role, nil
}

func GetInfo(roleID uint64) (roleInfo RoleInfo, err error) {
	if err = db.Mysql.
		Model(&Role{}).
		Where("id", roleID).
		Take(&roleInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return roleInfo, ErrRoleNotFound
		} else {
			return roleInfo, err
		}
	}
	return roleInfo, nil
}

func Update(roleID uint64, ur UpdateRole) error {
	if _, err := Get(roleID); err != nil {
		return err
	}

	updateRole := make(map[string]interface{})
	if ur.RoleName != "" {
		updateRole["role_name"] = ur.RoleName
	}

	err := db.Mysql.
		Model(&Role{}).
		Where("id = ?", roleID).
		Updates(updateRole).Error
	return err
}

func Delete(roleID uint64) error {
	if _, err := Get(roleID); err != nil {
		return err
	}

	var ur []users.UserRole

	if err := db.Mysql.
		Model(&users.UserRole{}).
		Where("role_id = ?", roleID).
		Find(&ur).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	if len(ur) > 0 {
		return ErrRoleBind
	}

	var gr []groups.GroupRole

	if err := db.Mysql.
		Model(&groups.GroupRole{}).
		Where("role_id = ?", roleID).
		Find(&gr).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	if len(gr) > 0 {
		return ErrRoleBind
	}

	if err := db.Mysql.Delete(&Role{}, roleID).Error; err != nil {
		return err
	}
	if _, err := auth.Casbin.RemoveFilteredGroupingPolicy(0, fmt.Sprintf("role::%d", roleID)); err != nil {
		return err
	}
	//if _, err := auth.Casbin.RemoveFilteredGroupingPolicy(1, fmt.Sprintf("role::%d", roleID)); err != nil {
	//	return err
	//}

	return nil
}

func List() (roles []RoleInfo, err error) {

	if err = db.Mysql.
		Model(&Role{}).
		Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func GetRoleMenus(roleID uint64) (roleMenu []RoleMenu, err error) {
	if err = db.Mysql.
		Model(&RoleMenu{}).
		Where("role_id = ?", roleID).
		Find(&roleMenu).Error; err != nil {
		return nil, err
	}

	return roleMenu, nil
}

func UpdateRoleMenu(roleID uint64, new []uint64) error {
	rm, err := GetRoleMenus(roleID)
	if err != nil {
		return err
	}
	old := make([]uint64, 0)
	for _, menu := range rm {
		old = append(old, menu.MenuID)
	}

	delRoleMenu := utils.Difference(old, new)
	addRoleMenu := utils.Difference(new, old)
	for _, mid := range delRoleMenu {
		db.Mysql.
			Where("role_id = ?", roleID).
			Where("menu_id = ?", mid).
			Delete(&RoleMenu{})
		if _, casbinErr := auth.Casbin.RemoveGroupingPolicy(fmt.Sprintf("role::%d", roleID), fmt.Sprintf("action::%d", mid)); casbinErr != nil {
			return casbinErr
		}
	}

	for _, mid := range addRoleMenu {
		nrm := RoleMenu{RoleID: roleID, MenuID: mid}
		db.Mysql.Create(&nrm)
		if _, casbinErr := auth.Casbin.AddGroupingPolicy(fmt.Sprintf("role::%d", roleID), fmt.Sprintf("action::%d", mid)); casbinErr != nil {
			return casbinErr
		}
	}
	return err
}
