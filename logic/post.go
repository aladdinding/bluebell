package logic

import (
	"bluebell_backend/dao/mysql"
	"bluebell_backend/models"
	"bluebell_backend/pkg/snowflake"
	"fmt"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

func CreatePost(post *models.Post) (err error) {
	// 生成帖子ID
	postID, err := snowflake.GetID()
	if err != nil {
		zap.L().Error("snowflake.GetID() failed", zap.Error(err))
		return
	}
	post.PostID = postID
	// 创建帖子
	if err := mysql.CreatePost(post); err != nil {
		zap.L().Error("mysql.CreatePost(&post) failed", zap.Error(err))
		return err
	}
	community, err := mysql.GetCommunityNameByID(fmt.Sprint(post.CommunityID))
	if err != nil {
		zap.L().Error("mysql.GetCommunityNameByID failed", zap.Error(err))
		return err
	}

	// TODO redis???
	if err := redis.CreatePost(
		fmt.Sprint(post.PostID),
		fmt.Sprint(post.AuthorId),
		post.Title,
		TruncateByWords(post.Content, 120),
		community.CommunityName); err != nil {
		zap.L().Error("redis.CreatePost failed", zap.Error(err))
		return err
	}
	return

}

func GetPost(postID string) (post *models.ApiPostDetail, err error) {
	post, err = mysql.GetPostByID(postID)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(postID) failed", zap.String())
	}
}
