package mgodb

import (
	"github.com/smalltree0/beego_goblog/RS"
	db "github.com/smalltree0/com/mongo"
	"gopkg.in/mgo.v2/bson"
)

const COUNTERNAME = "content"

type Content struct { // 文章
	ID      int32  // 文章唯一ID:t_id
	Content string // 内容
}

// Topic
func getNextArticleID() int32 {
	return db.NextVal(COUNTERNAME)
}

func createArticle(id int32, content string) int {
	newTopic := Content{ID: id, Content: content}
	err := db.Insert(DB, C_TOPIC, newTopic)
	if err != nil {
		return RS.RS_create_failed
	}
	return RS.RS_success
}

func deleteArticle(art_id int32) int {
	err := db.Remove(DB, C_TOPIC, bson.M{"id": art_id})
	if err != nil {
		return RS.RS_delete_failed
	}
	return RS.RS_success
}

func updateArticle(art_id int32, content string) int {
	c := getArticleComplete(art_id)
	c.Content = content
	if err := db.Update(DB, C_TOPIC, bson.M{"id": art_id}, *c); err != nil {
		return RS.RS_update_failed
	}
	return RS.RS_success
}

func getArticleComplete(art_id int32) *Content {
	var c Content
	db.FindOne(DB, C_TOPIC, bson.M{"id": art_id}, &c)
	return &c
}

func getArticleCompleteAll() *[]Content {
	var c []Content
	db.FindOne(DB, C_TOPIC, bson.M{}, &c)
	return &c
}
