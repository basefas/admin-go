package menus

import (
	"admin-go/internal/auth"
	"admin-go/internal/roles"
	"admin-go/internal/utils/db"
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	ErrMenuNotFound = errors.New("未找到菜单")
	ErrMenuExists   = errors.New("菜单已存在")
	ErrMenuBind     = errors.New("存在未删除绑定关系")
)

func Create(cm CreateMenu) error {
	var m Menu
	if err := db.Mysql.
		Where("name = ?", cm.Name).
		Where("path = ?", cm.Path).
		Where("method = ?", cm.Method).
		Find(&m).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	if m.ID != 0 {
		return ErrMenuExists
	}

	nm := Menu{
		Name:     cm.Name,
		Path:     cm.Path,
		Type:     cm.Type,
		Method:   cm.Method,
		Icon:     cm.Icon,
		ParentID: cm.ParentID,
		OrderID:  cm.OrderID}

	if err := db.Mysql.Create(&nm).Error; err != nil {
		return err
	}

	if cm.Type == 2 || cm.Type == 3 {
		if _, casbinErr := auth.Casbin.AddPolicy(fmt.Sprintf("action::%d", nm.ID), nm.Path, nm.Method); casbinErr != nil {
			return casbinErr
		}
	}

	return nil
}

func GetInfo(menuID uint64) (menuInfo MenuInfo, err error) {
	if err = db.Mysql.
		Model(&Menu{}).
		Where("id", menuID).
		Take(&menuInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return menuInfo, ErrMenuNotFound
		} else {
			return menuInfo, err
		}
	}
	return menuInfo, nil
}

func Get(menuID uint64) (menu Menu, err error) {
	if err = db.Mysql.Take(&menu, menuID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return menu, ErrMenuNotFound
		} else {
			return menu, err
		}
	}
	return menu, nil
}

func GetTree(menuID uint64) (menuInfo MenuInfo, err error) {
	menuInfo, err = GetInfo(menuID)
	if err != nil {
		return menuInfo, err
	}
	funs := make([]MenuInfo, 0)
	funs, err = FunListForPid(menuID)
	if err != nil {
		return menuInfo, err
	}
	menuInfo.Funs = funs
	return menuInfo, nil
}

func Update(menuID uint64, um UpdateMenu) error {
	oldMenu, err := Get(menuID)
	if err != nil {
		return err
	}

	updateMenu := make(map[string]interface{})
	if um.Name != "" {
		updateMenu["name"] = um.Name
	}
	if um.Path != "" {
		updateMenu["path"] = um.Path
	}
	if um.Type != 0 {
		updateMenu["type"] = um.Type
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

	if dbErr := db.Mysql.
		Model(&Menu{}).
		Where("id = ?", menuID).
		Updates(updateMenu).Error; dbErr != nil {
		return dbErr
	}

	newMenu, err1 := Get(menuID)
	if err1 != nil {
		return err1
	}

	if oldMenu.Type == 2 || oldMenu.Type == 3 {
		if _, casbinErr := auth.Casbin.RemovePolicy(fmt.Sprintf("action::%d", oldMenu.ID), oldMenu.Path, oldMenu.Method); casbinErr != nil {
			return casbinErr
		}
	}

	if newMenu.Type == 2 || newMenu.Type == 3 {
		if _, casbinErr := auth.Casbin.AddPolicy(fmt.Sprintf("action::%d", newMenu.ID), newMenu.Path, newMenu.Method); casbinErr != nil {
			return casbinErr
		}
	}

	return nil
}

func Delete(menuID uint64) error {
	if _, err := Get(menuID); err != nil {
		return err
	}

	rm := make([]roles.RoleMenu, 0)

	if err := db.Mysql.
		Model(&roles.RoleMenu{}).
		Where("menu_id = ?", menuID).
		Find(&rm).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	if len(rm) > 0 {
		return ErrMenuBind
	}

	if err := db.Mysql.Delete(&Menu{}, menuID).Error; err != nil {
		return err
	}

	if _, err := auth.Casbin.RemoveFilteredPolicy(0, fmt.Sprintf("action::%d", menuID)); err != nil {
		return err
	}

	return nil
}

func List() (menus []MenuInfo, err error) {
	menus = make([]MenuInfo, 0)
	if err = db.Mysql.
		Model(&Menu{}).
		Find(&menus).Error; err != nil {
		return nil, err
	}

	return menus, nil
}

func FunListForPid(pid uint64) (menus []MenuInfo, err error) {
	menus = make([]MenuInfo, 0)
	if err = db.Mysql.
		Model(&Menu{}).
		Where("type = ?", 3).
		Where("parent_id = ?", pid).
		Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func System(userID uint64) ([]MenuInfo, error) {
	return TreeForPid(0, userID)
}

func Tree() ([]MenuInfo, error) {
	return TreeForPid(0, 0)
}

func TreeForPid(pid, userID uint64) (menus []MenuInfo, err error) {
	l, listErr := list(userID)
	if listErr != nil {
		return menus, listErr
	}
	ml := make([]MenuInfo, 0)
	fl := make([]MenuInfo, 0)
	root := make([]MenuInfo, 0)
	menus = make([]MenuInfo, 0)
	for _, item := range l {
		if item.Type == 1 || item.Type == 2 {
			ml = append(ml, item)
		}
		if item.Type == 3 {
			fl = append(fl, item)
		}
		if item.ParentID == pid {
			root = append(root, item)
		}
	}
	for _, menu := range root {
		menu.Children = getMenuForPid(menu.ID, ml, fl)
		menu.Funs = getFunForPid(menu.ID, fl)
		menus = append(menus, menu)
	}
	return menus, err
}

func list(userID uint64) (menus []MenuInfo, err error) {
	if userID == 0 {
		if err = db.Mysql.
			Model(&Menu{}).
			Order("order_id").
			Find(&menus).Error; err != nil {
			return menus, err
		}
	} else {
		const q = `
		SELECT m.*
		FROM ag_menu AS m
		  LEFT JOIN ag_role_menu AS rm ON rm.menu_id = m.id
		  LEFT JOIN ag_user_role AS ur ON rm.role_id = ur.role_id
		WHERE rm.deleted_at IS NULL
		  AND ur.deleted_at IS NULL
		  AND m.deleted_at IS NULL
		  AND user_id = ?
		ORDER BY order_id`

		if err = db.Mysql.
			Raw(q, userID).
			Scan(&menus).Error; err != nil {
			return menus, err
		}
	}
	return menus, nil
}

func getMenuForPid(pid uint64, menuList, funList []MenuInfo) []MenuInfo {
	cl := make([]MenuInfo, 0)
	for _, menu := range menuList {
		if menu.ParentID == pid {
			cl = append(cl, menu)
		}
	}

	children := make([]MenuInfo, 0)
	for i, menu := range cl {
		cl[i].Children = getMenuForPid(menu.ID, menuList, funList)
		cl[i].Funs = getFunForPid(menu.ID, funList)
		children = append(children, cl[i])
	}

	return children
}

func getFunForPid(pid uint64, funs []MenuInfo) []MenuInfo {
	fs := make([]MenuInfo, 0)
	for i, fun := range funs {
		if pid == fun.ParentID {
			funs[i].Children = make([]MenuInfo, 0)
			funs[i].Funs = make([]MenuInfo, 0)
			fs = append(fs, funs[i])
		}
	}
	return fs
}
