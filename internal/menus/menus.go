package menus

import (
	"fmt"
	"go-admin/internal/auth"
	"go-admin/internal/utils/db"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var (
	ErrMenuNotFound = errors.New("Menu not found.")

	ErrMenuExists = errors.New("Menu already exist.")
)

func Create(cm CreateMenu) error {
	var m = Menu{}
	if err := db.Mysql.
		Where("menu_name = ?", cm.MenuName).
		Where("menu_path = ?", cm.MenuPath).
		Where("method = ?", cm.Method).
		Find(&m).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return err
		}
	}

	if m.ID != 0 {
		return ErrMenuExists
	}

	nm := Menu{MenuName: cm.MenuName, MenuPath: cm.MenuPath, MenuType: cm.MenuType, Method: cm.Method, Icon: cm.Icon, ParentID: cm.ParentID, OrderID: cm.OrderID}

	err := db.Mysql.Create(&nm).Error
	if cm.MenuType == 2 || cm.MenuType == 3 {
		_, err := auth.Casbin.AddPolicy(fmt.Sprintf("action::%d", nm.ID), nm.MenuPath, nm.Method)
		return err
	}
	return err
}

func Get(menuID uint) (*MenuInfo, error) {
	var m MenuInfo

	if err := db.Mysql.
		Model(&Menu{}).
		Where("id = ?", menuID).
		Scan(&m).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, ErrMenuNotFound
		}
	}

	return &m, nil
}

func GetTree(menuID uint) (*MenuInfo, error) {
	var m MenuInfo

	if err := db.Mysql.
		Model(&Menu{}).
		Where("id = ?", menuID).
		Scan(&m).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, ErrMenuNotFound
		}
	}
	funs := make([]*MenuInfo, 0)
	var err error
	funs, err = FunListForPid(menuID)
	if err != nil {
		return nil, err
	}
	m.Funs = funs
	return &m, nil
}

func Update(menuID uint, um UpdateMenu) error {
	oldMenu, err := Get(menuID)
	if err != nil {
		return err
	}

	updateMenu := make(map[string]interface{})
	if um.MenuName != "" {
		updateMenu["menu_name"] = um.MenuName
	}
	if um.MenuPath != "" {
		updateMenu["menu_path"] = um.MenuPath
	}
	if um.MenuType != 0 {
		updateMenu["menu_type"] = um.MenuType
	}
	if um.Method != "" {
		updateMenu["method"] = um.Method
	}
	if um.Icon != "" {
		updateMenu["icon"] = um.Icon
	}
	if um.OrderID != 999999 {
		updateMenu["order_id"] = um.OrderID
	}
	updateMenu["parent_id"] = um.ParentID

	if err := db.Mysql.
		Model(Menu{}).
		Where("id = ?", menuID).
		Updates(updateMenu).Error; err != nil {
		return err
	}

	newMenu, err1 := Get(menuID)
	if err1 != nil {
		return err1
	}

	if oldMenu.MenuType == 2 || oldMenu.MenuType == 3 {
		if _, err := auth.Casbin.RemovePolicy(fmt.Sprintf("action::%d", oldMenu.ID), oldMenu.MenuPath, oldMenu.Method); err != nil {
			return err
		}
	}

	if newMenu.MenuType == 2 || newMenu.MenuType == 3 {
		if _, err := auth.Casbin.AddPolicy(fmt.Sprintf("action::%d", newMenu.ID), newMenu.MenuPath, newMenu.Method); err != nil {
			return err
		}
	}

	return nil
}

func Delete(menuID uint) error {
	if _, err := Get(menuID); err != nil {
		return err
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
		Where("id = ?", menuID).
		Delete(&Menu{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.
		Where("menu_id = ?", menuID).
		Delete(&RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}

	if _, err := auth.Casbin.RemoveFilteredPolicy(0, fmt.Sprintf("action::%d", menuID)); err != nil {
		return err
	}

	return nil
}

func DeleteTree(menuID uint) error {
	if _, err := Get(menuID); err != nil {
		return err
	}

	menus, err := List()
	if err != nil {
		return err
	}
	ids := make([]uint, 0)
	ids = append(ids, menuID)
	ids = append(ids, getChildrenID(menuID, menus)...)

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
		Where("id IN (?)", ids).
		Delete(&Menu{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.
		Where("menu_id IN (?)", ids).
		Delete(&RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	for _, id := range ids {
		if _, err := auth.Casbin.RemoveFilteredPolicy(0, fmt.Sprintf("action::%d", id)); err != nil {
			return err
		}
	}
	return nil
}

func getChildrenID(pid uint, menus []MenuInfo) []uint {
	childrenID := make([]uint, 0)
	for _, menu := range menus {
		if menu.ParentID == pid {
			childrenID = append(childrenID, menu.ID)
		}
	}

	for _, id := range childrenID {
		childrenID = append(childrenID, getChildrenID(id, menus)...)
	}
	return childrenID
}

func List() ([]MenuInfo, error) {
	menus := make([]MenuInfo, 0)

	err := db.Mysql.
		Model(&Menu{}).
		Scan(&menus).Error
	if err != nil {
		return nil, err
	}

	return menus, nil
}

func FunListForPid(pid uint) ([]*MenuInfo, error) {
	var menus []*MenuInfo

	err := db.Mysql.
		Model(&Menu{}).
		Where("menu_type = ?", 3).
		Where("parent_id = ?", pid).
		Scan(&menus).Error
	if err != nil {
		return nil, err
	}

	return menus, nil
}

func System(userID uint) ([]MenuInfo, error) {
	return TreeForPid(0, userID)
}

func Tree() ([]MenuInfo, error) {
	return TreeForPid(0, 0)
}

func TreeForPid(pid, userID uint) ([]MenuInfo, error) {
	l, err := list(userID)
	ml := make([]MenuInfo, 0)
	fl := make([]MenuInfo, 0)
	root := make([]MenuInfo, 0)
	for _, item := range l {
		if item.MenuType == 1 || item.MenuType == 2 {
			ml = append(ml, item)
		}
		if item.MenuType == 3 {
			fl = append(fl, item)
		}
		if item.ParentID == pid {
			root = append(root, item)
		}
	}
	menus := make([]MenuInfo, 0)
	for _, menu := range root {
		menu.Children = getMenuChildren(menu.ID, ml, fl)
		menu.Funs = getFunForPid(menu.ID, fl)
		menus = append(menus, menu)
	}
	return menus, err
}

func list(userID uint) ([]MenuInfo, error) {
	menus := make([]MenuInfo, 0)
	if userID == 0 {
		if err := db.Mysql.
			Model(&Menu{}).
			Order("order_id asc").
			Scan(&menus).Error; err != nil {
			return nil, err
		}
	} else {
		if err := db.Mysql.
			Select("m.*").
			Table("role_menu AS rm").
			Joins("LEFT JOIN user_role AS ur ON rm.role_id = ur.role_id").
			Joins("LEFT JOIN menu AS m ON rm.menu_id = m.id").
			Where("rm.deleted_at IS NULL").
			Where("ur.deleted_at IS NULL").
			Where("m.deleted_at IS NULL").
			Where("user_id =?", userID).
			Order("order_id asc").
			Scan(&menus).Error; err != nil {
			return nil, err
		}
	}

	return menus, nil
}

func getMenuChildren(pid uint, menuList, funList []MenuInfo) []*MenuInfo {
	cl := make([]MenuInfo, 0)
	for _, menu := range menuList {
		if menu.ParentID == pid {
			cl = append(cl, menu)
		}
	}

	children := make([]*MenuInfo, 0)
	for i, menu := range cl {
		cl[i].Children = getMenuChildren(menu.ID, menuList, funList)
		cl[i].Funs = getFunForPid(menu.ID, funList)
		children = append(children, &cl[i])
	}

	return children
}

func getFunForPid(pid uint, funs []MenuInfo) []*MenuInfo {
	fs := make([]*MenuInfo, 0)
	for i, fun := range funs {
		if pid == fun.ParentID {
			funs[i].Children = make([]*MenuInfo, 0)
			funs[i].Funs = make([]*MenuInfo, 0)
			fs = append(fs, &funs[i])
		}
	}
	return fs
}
