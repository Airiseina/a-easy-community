package model

const (
	RoleUser  = 0 // 普通用户
	RoleAdmin = 1 // 管理员
)

type User struct {
	ID           uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Account      string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Hash         string `gorm:"not null"`
	Name         string `gorm:"type:varchar(50);not null"`
	Introduction string `gorm:"type:varchar(100);not null"`
	Role         int    `gorm:"type:tinyint;default:0;comment:角色 0:普通用户 1:管理员"`
	Avatar       string `gorm:"type:varchar(255);default:'';comment:头像URL"`
	IsMuted      bool   `gorm:"default:false;comment:是否禁言"`
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
}
