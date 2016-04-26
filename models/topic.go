package models

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/deepzz0/go-common/log"
	db "github.com/deepzz0/go-common/mongo"
	"github.com/deepzz0/goblog/RS"
	"github.com/deepzz0/goblog/helper"
	"github.com/russross/blackfriday"
	"gopkg.in/mgo.v2/bson"
)

const OnePageCount = 15

func init() {
	go scheduleTopic()
}

func scheduleTopic() {
	t := time.NewTicker(time.Minute * 10)
	for {
		select {
		case <-t.C:
			TMgr.DoDelete(time.Now())
		}
	}
}

// topic URL ＝ "/2016/01/02/id.html"
type Topic struct {
	ID         int32
	Author     string
	CreateTime time.Time
	EditTime   time.Time
	Title      string
	CategoryID string
	TagIDs     []string
	Content    string
	NeedDelete time.Time // 开始删除时间,超过48小时永久删除

	Preview   string    `bson:"-"`
	PCategory *Category `bson:"-"`
	PTags     []*Tag    `bson:"-"`
	URL       string    `bson:"-"`
	Time      string    `bson:"-"`
}

type TopicMgr struct {
	lock            sync.Mutex
	Topics          map[int32]*Topic // userid --> *User
	IDs             INT32
	GroupByCategory map[string]Topics
	GroupByTag      map[string]Topics
	DeleteTopics    Topics
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
		if !topic.NeedDelete.IsZero() {
			m.DeleteTopics = append(m.DeleteTopics, topic)
			continue
		}
		category := Blogger.GetCategoryByID(topic.CategoryID)
		if category == nil {
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
		m.DoTopicUpdate(topic)
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

func (m *TopicMgr) DoTopicUpdate(topic *Topic) {
	topic.Content = string(blackfriday.MarkdownCommon([]byte(topic.Content)))
	reg, _ := regexp.Compile(`\</\w{1,3}\>`)
	index := reg.FindAllStringIndex(topic.Content, 10)
	x := index[len(index)-1]
	topic.Preview = string(blackfriday.MarkdownCommon([]byte(topic.Content[:x[len(x)-1]])))

	topic.URL = fmt.Sprintf("%s/%d.html", topic.CreateTime.Format(helper.Layout_y_m_d), topic.ID)
	topic.Time = topic.CreateTime.Format(helper.Layout_y_m_d)
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

func (m *TopicMgr) LoadTopic(id int32) (*Topic, error) {
	var t *Topic
	err := db.FindOne(DB, C_TOPIC, bson.M{"id": id}, &t)
	if err != nil {
		return nil, err
	}
	return t, nil
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

func (m *TopicMgr) GetTopicsSearch(search string) []*Topic {
	var topics []*Topic
	for _, v := range m.Topics {
		if strings.Contains(strings.ToLower(v.Title), search) {
			topics = append(topics, v)
		}
	}
	return topics
}

func getPage(length int) int {
	page := length / OnePageCount
	if length%OnePageCount > 0 {
		page++
	}
	return page
}

func (m *TopicMgr) AddTopic(topic *Topic) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if err := db.Insert(DB, C_TOPIC, topic); err != nil {
		return err
	}
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
			newtag.Extra = "/tag/" + id
			newtag.Text = id
			m.GroupByTag[id] = append(m.GroupByTag[id], topic)
			sort.Sort(m.GroupByTag[id])
			topic.PTags = append(topic.PTags, newtag)
			Blogger.AddTag(newtag)
		}
	}
	m.Topics[topic.ID] = topic
	m.IDs = append(m.IDs, topic.ID)
	sort.Sort(m.IDs)

	m.DoTopicUpdate(topic)
	return nil
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
		if category == nil {
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
				newtag.Extra = "/tag/" + id
				newtag.Text = id
				m.GroupByTag[id] = append(m.GroupByTag[id], topic)
				topic.PTags = append(topic.PTags, newtag)
				Blogger.AddTag(newtag)
			}
			sort.Sort(m.GroupByTag[id])
		}
	}
	topic.EditTime = time.Now()
	if err := db.Update(DB, C_TOPIC, bson.M{"id": topic.ID}, topic); err != nil {
		return err
	}
	m.DoTopicUpdate(topic)
	return nil
}

func (m *TopicMgr) DelTopic(id int32) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if topic := m.GetTopic(id); topic != nil && topic.NeedDelete.IsZero() {
		topic.NeedDelete = time.Now()
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
		m.DeleteTopics = append(m.DeleteTopics, topic)
		delete(m.Topics, id)
		return nil
	}
	return fmt.Errorf("Topic id=%d not found in cache.", id)
}

func (m *TopicMgr) RestoreTopic(topic *Topic) int {
	if topic.NeedDelete.IsZero() {
		return RS.RS_notin_trash
	}
	topic.NeedDelete = time.Time{}
	err := m.ModTopic(topic, topic.CategoryID, strings.Join(topic.TagIDs, ","))
	if err != nil {
		return RS.RS_undo_falied
	}
	return RS.RS_success
}

func (m *TopicMgr) DoDelete(t time.Time) {
	for _, topic := range m.DeleteTopics {
		if topic.NeedDelete.AddDate(0, 0, 2).Before(t) {
			db.Remove(DB, C_TOPIC, bson.M{"id": topic.ID})
		}
	}
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
