package mgodb

import (
	// "fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/smalltree0/beego_goblog/RS"
	"github.com/smalltree0/beego_goblog/helper"
	"github.com/smalltree0/com/log"
	db "github.com/smalltree0/com/mongo"
	"gopkg.in/mgo.v2/bson"
)

// 账户
type User struct {
	/* 注册数据 */
	UserName string // 账户ID
	PassWord string // 账户密码
	Email    string // email
	/* 资料补充 */
	Sex       int    // 1:男， 2：女
	PNumber   int64  // 手机号
	Address   string // 住址
	Education string // 教育
	RealName  string // 真实姓名
	/* 自动记录 */
	CreateTime time.Time // 创建时间
	LoginTime  time.Time // 登录时间
	LogoutTime time.Time // 登出时间

	/* 个性设置 */
	BlogName  string // 博客名
	Motto     string // 座右铭
	Introduce string // 个人简介
	Bg        string // 图片背景
	/* 博客信息 */
	Tags       []Tag      // 标签
	Categories []Category // 分类
	Links      []Link     // 友情链接
	/* 文章id */
	Previews PRE  // 文章基础预览
	NeedSort bool // 是否该排序
}

type UserMgr struct {
	lock  sync.Mutex
	Users map[string]*User // userid >> *User
}

func NewUM() *UserMgr { return &UserMgr{Users: make(map[string]*User)} }

var UMgr = NewUM()

func init() {
	UMgr.loadUsers()
	go schedule()
}

func schedule() {
	tk := time.NewTicker(time.Hour)
	for {
		select {
		case <-tk.C:
			UMgr.UpdateUsers()
		}
	}
}

func (m *UserMgr) loadUsers() {
	var users []*User
	err := db.FindAll(DB, C_USER, nil, &users)
	if err != nil {
		panic(err)
	}
	for _, u := range users {
		u.Sort()
		m.Users[u.UserName] = u
	}
}

func (m *UserMgr) RegisterUser(name, passwd, email string) int {
	m.lock.Lock()
	defer m.lock.Unlock()
	name = strings.ToLower(name)
	passwd = strings.ToLower(passwd)
	if ok := db.KeyIsExsit(DB, C_USER, "username", name); ok {
		return RS.RS_user_exist
	}
	newpasswd := helper.EncryptPasswd(name, passwd)
	user := User{
		UserName:   name,
		PassWord:   newpasswd,
		Email:      email,
		CreateTime: time.Now(),
	}
	err := db.Insert(DB, C_USER, user)
	if err != nil {
		log.Error(err)
		return RS.RS_failed
	}
	m.Users[name] = &user
	return RS.RS_success
}

func (m *UserMgr) FoundPass(name, email string) int {
	if user, found := m.Users[name]; !found {
		return RS.RS_user_inexistence
	} else {
		log.Debug(user.UserName)
	}

	// 发送邮件
	// ....
	return RS.RS_success
}

func (m *UserMgr) LoginUser(name, passwd string) int {
	user := m.Users[name]
	if user == nil {
		return RS.RS_user_inexistence
	}
	if user.PassWord != helper.EncryptPasswd(name, passwd) {
		return RS.RS_password_error
	}
	user.LoginTime = time.Now()
	return RS.RS_success
}

func (m *UserMgr) LogoutUser(name string) int {
	user := m.Users[name]
	if user == nil {
		return RS.RS_user_inexistence
	}
	user.LogoutTime = time.Now()
	db.Update(DB, C_USER, bson.M{"accountid": name}, *m.Users[name])
	return RS.RS_success
}

func (m *UserMgr) Get(name string) *User {
	return m.Users[name]
}

func (m *UserMgr) UpdateUsers() int {
	for _, u := range m.Users {
		err := db.Update(DB, C_USER, bson.M{"username": u.UserName}, *u)
		if err != nil {
			return RS.RS_update_failed
		}
	}
	return RS.RS_success
}

//-------------------------------user-------------------------------------
func (u *User) ChangePass(oldpass, newpass string) int {
	if u.PassWord != helper.EncryptPasswd(u.UserName, oldpass) {
		return RS.RS_password_error
	}
	u.PassWord = helper.EncryptPasswd(u.UserName, newpass)
	return RS.RS_success
}

func (u *User) ChangeInfo(args map[string]interface{}) int {
	sex := args["sex"].(int)
	pn := args["phonenumber"].(int64)
	addr := args["address"].(string)
	ed := args["education"].(string)
	rn := args["realname"].(string)

	u.Sex = sex
	u.PNumber = pn
	u.Address = addr
	u.Education = ed
	u.RealName = rn
	return RS.RS_success
}

type Preview struct {
	ID         int32     // 文章ID
	Title      string    // 标题
	Content    string    // 预览内容
	Group      string    // 分类
	Tags       []Tag     // 标签
	Views      int       // 浏览数
	IsDrafts   bool      `json:"-"` // 是否草稿
	TopTime    time.Time `json:"-"` // 置顶时间(0:不置顶, >0:置顶)
	CreateTime time.Time // 创建时间
	UpdateTime time.Time // 更新时间
}

type Link struct {
	Addr string // 链接地址
	Name string // 链接名
}

type Tag struct { // 标签
	Name string
}

type Category struct { // 分类
	Name  string
	Count int64
}

// preview
type PRE []*Preview

func (p PRE) Less(i, j int) bool {
	if p[i].TopTime.After(p[j].TopTime) {
		return true
	} else if p[i].CreateTime.After(p[j].CreateTime) {
		return true
	}
	return false
}
func (p PRE) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p PRE) Len() int {
	return len(p)
}

func (u *User) Sort() {
	if u.NeedSort {
		sort.Sort(u.Previews)
		u.NeedSort = false
	}
}

const pagecount = 15

func (u *User) GetSimplePreview(id int32) *Preview {
	u.Sort()
	for _, pre := range u.Previews {
		if pre.ID == id {
			return pre
		}
	}
	return nil
}

func getInt(length int, count int) (maxpage int) {
	maxpage = length / count
	over := length % count
	if over > 0 {
		maxpage += 1
	}
	return
}

func (u *User) GetPreviews(page int) ([]*Preview, int) {
	u.Sort()
	length := len(u.Previews)
	if length <= pagecount {
		return u.Previews, 1
	}
	maxpage := getInt(length, pagecount)

	index := page * pagecount
	end := length - index
	start := end - pagecount
	if start < 0 {
		start = 0
	}
	return u.Previews[start:end], maxpage
}

func (u *User) GetPreviewByCategory(name string, page int) ([]*Preview, int) {
	u.Sort()
	pw := []*Preview{}
	for _, v := range u.Previews {
		if v.Group == name {
			pw = append(pw, v)
		}
	}

	length := len(pw)
	if length <= pagecount {
		return pw, 1
	}
	maxpage := getInt(length, pagecount)
	index := page * pagecount
	end := length - index
	start := end - pagecount
	if start < 0 {
		start = 0
	}
	return pw[start:end], maxpage
}

func (u *User) GetPreviewByTag(name string, page int) ([]*Preview, int) {
	u.Sort()
	pw := []*Preview{}
	for _, v := range u.Previews {
		for _, val := range v.Tags {
			if name == val.Name {
				pw = append(pw, v)
				break
			}
		}
	}
	length := len(pw)
	if length <= pagecount {
		return pw, 1
	}
	maxpage := getInt(length, pagecount)
	index := page * pagecount
	end := length - index
	start := end - pagecount
	if start < 0 {
		start = 0
	}
	return pw[start:end], maxpage
}

func (u *User) TopPreview(id int32) {
	for _, v := range u.Previews {
		if v.ID == id {
			v.TopTime = time.Now()
			break
		}
	}
	u.NeedSort = true
}

func (u *User) CancelTopview(id int32) {
	for _, v := range u.Previews {
		if v.ID == id {
			v.TopTime = time.Time{}
			break
		}
	}
	u.NeedSort = true
}

func (u *User) AddPreview(title, content, group, tag string, isdrafts, istop bool) int {
	tag_slice := strings.Split(tag, "|")
	var tags []Tag
	for _, v := range tag_slice {
		tags = append(tags, Tag{v})
	}
	runes := []rune(content)
	pre_content := ""
	count := 0
	for i, v := range runes {
		if v == 10 { // 10 代表回车
			count++
			if count >= 8 {
				pre_content = string(runes[:i]) + " ......"
				break
			}
		}
	}
	if count < 8 {
		pre_content = string(runes)
	}

	id := getNextArticleID()
	var toptime time.Time
	if istop {
		toptime = time.Now()
	}
	var pre = &Preview{
		ID:         id,
		Title:      title,
		Content:    pre_content,
		Group:      group,
		Tags:       tags,
		IsDrafts:   isdrafts,
		TopTime:    toptime,
		CreateTime: time.Now(),
	}

	u.Previews = append(u.Previews, pre)
	u.addCategoryCount(pre.Group)
	createArticle(id, content)

	u.NeedSort = true
	return RS.RS_success
}

func (u *User) UpdatePreview(art_id int32, title, content, group string, tags []Tag) int {
	for _, v := range u.Previews {
		if v.ID == art_id {
			v.UpdatePreview(art_id, title, content, group, tags)
		}
	}
	return RS.RS_success
}

func (u *User) DelPreview(art_id int32) int {
	for i := len(u.Previews) - 1; i >= 0; i-- {
		if pw := u.Previews[i]; art_id == pw.ID {
			u.delCategoryCount(pw.Group)
			u.Previews = append(u.Previews[:i], u.Previews[i+1:]...)
		}
	}
	return RS.RS_success
}

func (u *User) GetContent(art_id int32) *Content {
	return getArticleComplete(art_id)
}

// category
func (u *User) addCategoryCount(name string) {
	for i, v := range u.Categories {
		if v.Name == name {
			u.Categories[i].Count++
		}
	}
}

func (u *User) delCategoryCount(name string) {
	for i, v := range u.Categories {
		if v.Name == name {
			if v.Count > 0 {
				u.Categories[i].Count--
			}
		}
	}
}

func (u *User) GetCategories() []Category {
	return u.Categories
}

func (u *User) AddCategory(name string) {
	c := Category{Name: name}
	u.Categories = append(u.Categories, c)
}

func (u *User) DelCategory(name string) {
	for i, v := range u.Categories {
		if v.Name == name {
			u.Categories = append(u.Categories[:i], u.Categories[i+1:]...)
		}
	}
}

// tag
func (u *User) AddTag(name string) {
	t := Tag{Name: name}
	u.Tags = append(u.Tags, t)
}

func (u *User) GetTags() []Tag {
	return u.Tags
}

func (u *User) DelTag(name string) {
	for i, v := range u.Tags {
		if v.Name == name {
			u.Tags = append(u.Tags[:i], u.Tags[i+1:]...)
		}
	}
}

//
func (u *User) ModBlogName(name string) {
	u.BlogName = name
}

func (u *User) ModMotto(motto string) {
	u.Motto = motto
}

//-----------------------------preview------------------------------
func (pre *Preview) UpdatePreview(art_id int32, title, content, group string, tags []Tag) {
	updateArticle(art_id, content)
	pre.Title = title
	runes := []rune(content)
	pre_content := ""
	count := 0
	for i, v := range runes {
		if v == 10 { // 10 代表回车
			count++
			if count >= 8 {
				pre_content = string(runes[:i]) + " ......"
				break
			}
		}
	}
	if count < 8 {
		pre_content = string(runes)
	}
	pre.Content = pre_content
	pre.Group = group
	pre.Tags = tags
	pre.UpdateTime = time.Now()
}
