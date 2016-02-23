package models

import (
	"fmt"
	// "reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/deepzz0/go-common/log"
	db "github.com/deepzz0/go-common/mongo"
	"github.com/deepzz0/goblog/RS"
	"github.com/deepzz0/goblog/helper"
	"gopkg.in/mgo.v2/bson"
)

const OnePageCount = 15

// topic URL ＝ "/2016/01/02/id.html"
type Topic struct {
	ID         int32
	Author     string
	CreateTime time.Time
	EditTime   time.Time
	Title      string
	CategoryID string
	TagIDs     []string
	Content    []rune

	PCategory *Category `bson:"-"`
	PTags     []*Tag    `bson:"-"`
}

type TopicMgr struct {
	lock            sync.Mutex
	Topics          map[int32]*Topic // userid --> *User
	IDs             INT32
	GroupByCategory map[string]Topics
	GroupByTag      map[string]Topics
}

func NewTopic() *Topic {
	return &Topic{ID: NextVal(), CreateTime: time.Now(), EditTime: time.Now(), Author: Blogger.UserName}
}
func NextVal() int32 {
	return db.NextVal(C_TOPIC_ID)
}

type INT32 []int32

func (t INT32) Len() int           { return len(t) }
func (t INT32) Less(i, j int) bool { return t[i] > t[j] }
func (t INT32) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

type Topics []*Topic

func (ts Topics) Len() int           { return len(ts) }
func (ts Topics) Less(i, j int) bool { return ts[i].ID > ts[j].ID }
func (ts Topics) Swap(i, j int)      { ts[i], ts[j] = ts[j], ts[i] }

func NewTM() *TopicMgr {
	return &TopicMgr{Topics: make(map[int32]*Topic), GroupByCategory: make(map[string]Topics), GroupByTag: make(map[string]Topics)}
}

var TMgr = NewTM()

func (m *TopicMgr) loadTopics() {
	var topics []*Topic
	err := db.FindAll(DB, C_TOPIC, bson.M{"author": Blogger.UserName}, &topics)
	if err != nil {
		panic(err)
	}
	length := len(topics)
	m.IDs = make([]int32, 0, length)
	for _, topic := range topics {
		category := Blogger.GetCategoryByID(topic.CategoryID)
		if category != nil {
			topic.CategoryID = "default"
			category = Blogger.GetCategoryByID(topic.CategoryID)
		}
		topic.PCategory = category
		m.GroupByCategory[topic.CategoryID] = append(m.GroupByCategory[topic.CategoryID], topic)
		for i, id := range topic.TagIDs {
			if tag := Blogger.GetTagByID(id); tag != nil {
				topic.PTags = append(topic.PTags, tag)
				m.GroupByTag[id] = append(m.GroupByTag[id], topic)
			} else {
				topic.TagIDs = append(topic.TagIDs[:i], topic.TagIDs[i+1:]...)
			}
		}
		m.Topics[topic.ID] = topic
		m.IDs = append(m.IDs, topic.ID)
	}
	sort.Sort(m.IDs)
	for k, _ := range m.GroupByCategory {
		sort.Sort(m.GroupByCategory[k])
	}
	for k, _ := range m.GroupByTag {
		sort.Sort(m.GroupByTag[k])
	}
}

func (m *TopicMgr) UpdateTopics() int {
	for _, topic := range m.Topics {
		err := db.Update(DB, C_TOPIC, bson.M{"id": topic.ID}, topic)
		if err != nil {
			log.Error(err)
			return RS.RS_update_failed
		}
	}
	return RS.RS_success
}

func (m *TopicMgr) GetTopic(id int32) *Topic {
	return m.Topics[id]
}
func (m *TopicMgr) GetTopicsByPage(page int) ([]*Topic, int) {
	var ts []*Topic
	if page <= 0 {
		return ts, -1
	}
	length := len(m.IDs)
	remainPage := getPage(length) - page
	if remainPage >= 0 {
		index := page * OnePageCount
		for i := index - OnePageCount; i < length && i < index; i++ {
			ts = append(ts, m.Topics[m.IDs[i]])
		}
		return ts, remainPage
	}
	return ts, -1
}
func (m *TopicMgr) GetTopicsByCatgory(categoryID string, page int) ([]*Topic, int) {
	if page <= 0 {
		return make([]*Topic, 0), -1
	}
	topics := m.GroupByCategory[categoryID]
	length := len(topics)
	remainPage := getPage(length) - page
	if remainPage >= 0 {
		var start, end int
		end = page * OnePageCount
		start = end - OnePageCount
		if end > length {
			end = length
		}
		if start < 0 {
			start = 0
		}
		return topics[start:end], remainPage
	}
	return make([]*Topic, 0), -1
}
func (m *TopicMgr) GetTopicsByTag(tagID string, page int) ([]*Topic, int) {
	if page <= 0 {
		return make([]*Topic, 0), -1
	}
	topics := m.GroupByTag[tagID]
	length := len(topics)
	remainPage := getPage(length) - page

	if remainPage >= 0 {
		var start, end int
		end = page * OnePageCount
		start = end - OnePageCount
		if end > length {
			end = length
		}
		if start < 0 {
			start = 0
		}
		return topics[start:end], remainPage
	}
	return make([]*Topic, 0), -1
}
func getPage(length int) int {
	page := length / OnePageCount
	if length%OnePageCount > 0 {
		page++
	}
	return page
}

func (m *TopicMgr) AddTopic(topic *Topic, domain string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	category := Blogger.GetCategoryByID(topic.CategoryID)
	if category == nil {
		topic.CategoryID = "default"
		category = Blogger.GetCategoryByID(topic.CategoryID)
	}
	m.GroupByCategory[topic.CategoryID] = append(m.GroupByCategory[topic.CategoryID], topic)
	sort.Sort(m.GroupByCategory[topic.CategoryID])
	topic.PCategory = category
	category.addCount()
	for _, id := range topic.TagIDs {
		if tag := Blogger.GetTagByID(id); tag != nil {
			m.GroupByTag[id] = append(m.GroupByTag[id], topic)
			topic.PTags = append(topic.PTags, tag)
			tag.addCount()
		} else {
			newtag := NewTag()
			newtag.ID = id
			newtag.Node = &helper.Node{Type: "a", Extra: fmt.Sprintf("href='/tag/%s'", id), Text: id}
			newtag.addCount()
			Blogger.Tags[id] = newtag
			m.GroupByTag[id] = append(m.GroupByTag[id], topic)
			sort.Sort(m.GroupByTag[id])
			topic.PTags = append(topic.PTags, newtag)
		}
	}
	m.Topics[topic.ID] = topic
	m.IDs = append(m.IDs, topic.ID)
	sort.Sort(m.IDs)
	return db.Insert(DB, C_TOPIC, topic)
}

func (m *TopicMgr) CategoryGroupDeleteTopic(topic *Topic) {
	topics := m.GroupByCategory[topic.CategoryID]
	for i, t := range topics {
		if t != nil && t == topic {
			Blogger.ReduceCategoryCount(topic.CategoryID)
			if m.GroupByCategory[topic.CategoryID] != nil {
				m.GroupByCategory[topic.CategoryID] = append(m.GroupByCategory[topic.CategoryID][:i], m.GroupByCategory[topic.CategoryID][i+1:]...)
			}
		}
	}
}

func (m *TopicMgr) TagGroupDeleteTopic(id string, topic *Topic) {
	topics := m.GroupByTag[id]
	for i, t := range topics {
		if t != nil && t == topic {
			Blogger.ReduceTagCount(id)
			if m.GroupByTag[id] != nil {
				m.GroupByTag[id] = append(m.GroupByTag[id][:i], m.GroupByTag[id][i+1:]...)
			}
		}
	}
}

func (m *TopicMgr) ModTopic(topic *Topic, catgoryID string, tags string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if topic.CategoryID != catgoryID {
		m.CategoryGroupDeleteTopic(topic)
		category := Blogger.GetCategoryByID(catgoryID)
		if category != nil {
			topic.CategoryID = "default"
			category = Blogger.GetCategoryByID(catgoryID)
		}
		topic.CategoryID = catgoryID
		topic.PCategory = category
		m.GroupByCategory[catgoryID] = append(m.GroupByCategory[catgoryID], topic)
		sort.Sort(m.GroupByCategory[catgoryID])
		category.addCount()
	}
	if tags == "" {
		topic.TagIDs = make([]string, 0)
		topic.PTags = make([]*Tag, 0)
	} else {
		tagIDS := strings.Split(tags, ",")
		for _, id := range topic.TagIDs {
			m.TagGroupDeleteTopic(id, topic)
		}
		topic.TagIDs = make([]string, 0, len(tagIDS))
		topic.PTags = make([]*Tag, 0, len(tagIDS))
		for _, id := range tagIDS {
			topic.TagIDs = append(topic.TagIDs, id)
			if tag := Blogger.GetTagByID(id); tag != nil {
				topic.PTags = append(topic.PTags, tag)
				tag.addCount()
				m.GroupByTag[id] = append(m.GroupByTag[id], topic)
			} else {
				newtag := NewTag()
				newtag.ID = id
				newtag.Node = &helper.Node{Type: "a", Extra: fmt.Sprintf("href='/tag/%s'", id), Text: id}
				newtag.addCount()
				Blogger.Tags[id] = newtag
				m.GroupByTag[id] = append(m.GroupByTag[id], topic)
				topic.PTags = append(topic.PTags, newtag)
			}
			sort.Sort(m.GroupByTag[id])
		}
	}
	return db.Update(DB, C_TOPIC, bson.M{"id": topic.ID}, topic)
}

func (m *TopicMgr) DelTopic(id int32) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if topic := m.GetTopic(id); topic != nil {
		if err := db.Remove(DB, C_TOPIC, bson.M{"id": id}); err != nil {
			return err
		}
		if topic.CategoryID != "" {
			m.CategoryGroupDeleteTopic(topic)
		}
		for _, id := range topic.TagIDs {
			m.TagGroupDeleteTopic(id, topic)
		}
		for i, id := range m.IDs {
			if id == topic.ID {
				m.IDs = append(m.IDs[:i], m.IDs[i+1:]...)
			}
		}
		delete(m.Topics, id)
		return nil
	}
	return fmt.Errorf("Topic id=%d not found in cache.", id)
}

// -----------------------------------------------------------------
func IsHaveTag(id string, ids []string) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}
