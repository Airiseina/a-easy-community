package global

import (
	"commmunity/app/internal/db/msq"
	"commmunity/app/internal/db/red"
)

var (
	User  msq.UserData  = msq.NewGorm(msq.ConnectMysql())
	Redis red.UserRedis = red.NewRedis(red.ConnectRedis())
)
