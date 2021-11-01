package redis

import (
	"github.com/go-redis/redis"
	"math"
	"time"
)

const (
	OneWeekInSeconds         = 7 * 24 * 3600
	VoteScore        float64 = 432
	PostPerAge               = 20
)

func PostVote(postID, userID string, v float64) (err error) {
	// 1. 帖子发布时间
	postTime := client.ZScore(KeyPostTimeZSet, postID).Val()
	if float64(time.Now().Unix())-postTime > OneWeekInSeconds {

		// 不允许投票了
		return ErrorVoteTimeExpire
	}
	key := KeyPostInfoHashPrefix + postID
	ov := client.ZScore(key, userID).Val() //获取当前分数

	diffAbs := math.Abs(ov - v)
	pipeline := client.TxPipeline()
	pipeline.ZAdd(key, redis.Z{
		Score:  v,
		Member: userID,
	})
	pipeline.ZIncrBy(KeyPostTimeZSet, VoteScore*diffAbs*v, postID)
	switch math.Abs(ov) - math.Abs(v) {
	case 1:
		pipeline.HIncrBy(KeyPostInfoHashPrefix+postID, "votes", -1)
	case 0:
	case -1:
		pipeline.HIncrBy(KeyPostInfoHashPrefix+postID, "votes", 1)
	default:

		return ErrorVoted

	}
	_, err = pipeline.Exec()
	return
}

func CreatePost(postID, userID, title, summary, communityName string) (err error) {
	now := float64(time.Now().Unix())
	votedKey := KeyPostVotedZSetPrefix + postID
	communityKey := KeyCommunityPostSetPrefix + communityName
	postInfo := map[string]interface{}{
		"title":    title,
		"summary":  summary,
		"post:id":  postID,
		"user:id":  userID,
		"time":     now,
		"votes":    1,
		"comments": 0,
	}
	//事务操作
	pipeline := client.TxPipeline()
	pipeline.ZAdd(votedKey, redis.Z{
		Score:  1,
		Member: userID,
	})

	pipeline.Expire(votedKey, time.Second*OneWeekInSeconds)

	pipeline.HMSet(KeyPostInfoHashPrefix+postID, postInfo)
	pipeline.ZAdd(KeyPostInfoHashPrefix, redis.Z{
		Score:  now + VoteScore,
		Member: postID,
	})
	pipeline.ZAdd(KeyPostTimeZSet, redis.Z{
		Score:  now,
		Member: postID,
	})
	pipeline.SAdd(communityKey, postID)
	_, err = pipeline.Exec()
	return

}
