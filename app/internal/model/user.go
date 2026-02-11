package model

import "gorm.io/gorm"

const (
	RoleUser  = 0 // 普通用户
	RoleAdmin = 1 // 管理员
)

type User struct {
	gorm.Model
	Account     string      `gorm:"type:varchar(100);uniqueIndex;not null"`
	Hash        string      `gorm:"not null"`
	Role        int         `gorm:"type:tinyint;default:0;comment:角色 0:普通用户 1:管理员"`
	Followings  []*User     `gorm:"many2many:user_relations;joinForeignKey:follower_id;joinReferences:followed_id"`
	Followers   []*User     `gorm:"many2many:user_relations;joinForeignKey:followed_id;joinReferences:follower_id"`
	UserProfile UserProfile `gorm:"foreignKey:UserID" json:"user_profile"`
	Posts       []Post      `gorm:"foreignKey:UserID" json:"posts"`
}

type UserProfile struct {
	gorm.Model
	UserID       uint   `gorm:"uniqueIndex;not null"`
	Name         string `gorm:"type:varchar(50);not null"`
	Introduction string `gorm:"type:varchar(100);not null"`
	Avatar       string `gorm:"type:varchar(255);default:'';comment:头像URL"`
	IsMuted      bool   `gorm:"default:false;comment:是否禁言"`
}

type UserRelation struct {
	FollowerID uint `gorm:"uniqueIndex:idx_relation"`
	FollowedID uint `gorm:"uniqueIndex:idx_relation"`
	Follower   User `gorm:"foreignKey:FollowerID"`
	Followed   User `gorm:"foreignKey:FollowedID"`
}
type UserRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type UserPassword struct {
	Account        string `json:"account"`
	FirstPassWord  string `json:"first_password"`
	SecondPassWord string `json:"second_password"`
}

type UserInfoRequest struct {
	Account      string `json:"account"`
	Introduction string `json:"introduction"`
	Avatar       string `json:"avatar"`
	IsMuted      bool   `json:"is_muted"`
}
