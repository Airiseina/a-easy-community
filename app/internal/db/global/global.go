package global

import (
	"commmunity/app/internal/db/msq"
	"commmunity/app/internal/db/red"
)

var (
	User      msq.UserData  = msq.NewGorm(msq.ConnectMysql())
	UserRedis red.UserRedis = red.NewRedis(red.ConnectRedis())
	Post      msq.PostData  = msq.NewGorm(msq.ConnectMysql())
	PostRedis red.PostRedis = red.NewRedis(red.ConnectRedis())
)
