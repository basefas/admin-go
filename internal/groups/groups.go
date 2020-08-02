package groups

import (
	"go-admin/internal/utils/db"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var (
	ErrGroupNotFound = errors.New("Group not found.")

	ErrGroupExists = errors.New("Group already exist.")
)

func Create(cg CreateGroup) error {
	if err := db.Mysql.
		Where("group_name = ?", cg.GroupName).
		Find(&Group{}).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return ErrGroupExists
		}
	}

	g := Group{GroupName: cg.GroupName, ParentID: cg.ParentID}

	err := db.Mysql.Create(&g).Error
	return err
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
	updateGroup["parent_id"] = ug.ParentID

	err := db.Mysql.
		Debug().
		Model(Group{}).
		Where("id = ?", groupID).
		Updates(updateGroup).Error
	return err
}

func Delete(groupID string) error {
	if _, err := Get(groupID); err != nil {
		return err
	}

	err := db.Mysql.
		Where("id = ?", groupID).
		Delete(&Group{}).Error
	return err
}

func List() ([]GetGroupInfo, error) {
	var groups []GetGroupInfo
	sql := `
		SELECT g.id, g.group_name, g.parent_id, g.created_at, g.updated_at, n.num AS head_count
		FROM group_ AS g
		LEFT JOIN (
			SELECT g.id, count(*) AS num
			FROM group_ AS g
			LEFT JOIN user_group ug ON g.id = ug.group_id
			GROUP BY g.id) AS n
		ON g.id= n.id
		WHERE g.deleted_at IS NULL
		`
	err := db.Mysql.
		Raw(sql).
		Scan(&groups).Error
	if err != nil {
		return nil, err
	}

	return groups, nil
}
