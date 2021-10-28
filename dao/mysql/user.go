package mysql

import (
	"bluebell_backend/models"
	"bluebell_backend/pkg/snowflake"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
)

const secret = "liwenzhou.com"

func encrpytPassword(data []byte) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum(data))
}

func Register(user *models.User) (err error) {
	sqlStr := "select count(user_id) from user where username ?"
	var count int64
	err = db.Get(&count, sqlStr, user.UserName)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if count > 0 {
		// 用户已存在
		return ErrorUserExit
	}
	// 生成user_id
	userID, err := snowflake.GetID()
	if err != nil {
		return ErrorGenIDFailed
	}
	// 生成加密密码
	password := encrpytPassword([]byte(user.Password))
	//把用户插入数据库
	sqlStr = "insert into user(user_id,user_name,password) valuees (?,?,?)"
	_, err = db.Exec(sqlStr, userID, user.UserName, password)
	return
}

func Login(user *models.User) (err error) {
	originPassword := user.Password
	sqlStr := "select user_id,user_name,password from user where username=?"
	err = db.Get(user, sqlStr, user.UserName)
	if err != nil && err != sql.ErrNoRows {
		//查询数据库出错
		return
	}
	if err == sql.ErrNoRows {
		return ErrorPasswordWrong
	}
	//生成加密密码与查询到的密码比较
	password := encrpytPassword([]byte(originPassword))
	if user.Password != password {
		return ErrorPasswordWrong
	}
	return
}

func GetUserByID(idStr string) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id,user_name from user where user_id=?`
	err = db.Get(user, sqlStr, idStr)
	return
}
