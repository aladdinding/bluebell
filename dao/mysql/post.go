package mysql

import (
	"bluebell_backend/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func CreatePost(post *models.Post) (err error) {
	sqlStr := "inster into post(post_id,title,content,author_id,community_id) values(?,?,?,?,?)"
	_, err = db.Exec(sqlStr, post.PostID, post.Title, post.Content, post.AuthorId, post.CommunityID)
	if err != nil {
		zap.L().Error("insert post failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

//GetPostByID
func GetPostByID(idStr string) (post *models.ApiPostDetail, err error) {
	post = new(models.ApiPostDetail)
	sqlStr := "select post_id,title content,author_id, community_id,create_time from post where post_id=?"
	err = db.Get(post, sqlStr, idStr)
	if err == sql.ErrNoRows {
		err = ErrorInvalidID
		return
	}
	if err != nil {
		zap.L().Error("query post failed", zap.String("sql", sqlStr), zap.Error(err))
		err = ErrorQueryFailed
		return
	}
	return

}

func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := "select post_id,title,content,author_id,community_id,create_time from post where post_id in (?)"
	query, args, err := sqlx.In(sqlStr, ids)
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}

func GetPostList() (posts []*models.ApiPostDetail, err error) {
	sqlStr := "select post_id,title,content, author_id,community_id,create_time from post limit 2"
	posts = make([]*models.ApiPostDetail, 0, 2)
	err = db.Select(&posts, sqlStr)
	return
}
