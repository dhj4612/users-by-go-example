package model

type UserPermission struct {
	UserId       int64 `gorm:"column:user_id" json:"userId"`
	PermissionId int64 `gorm:"column:permission_id" json:"permissionId"`
}

func (*UserPermission) TableName() string {
	return "user_permission"
}
