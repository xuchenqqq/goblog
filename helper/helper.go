package helper

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/smalltree0/beego_goblog/RS"
)

// var Mydomain = "http://deepzz.com"

const (
	Layout_y_m_d      = "2006/01/02"
	Layout_y_m_d2     = "2006年01月02日"
	Layout_y_m_d_time = "2006/01/02 15:04:05"
)

// -------------------------- response --------------------------
const (
	WARNING = "warning"
	SUCCESS = "success"
	ALERT   = "alert"
	INFO    = "info"
)

type Response struct {
	Status int
	Data   interface{}
	Err    Error
}
type Error struct {
	Level string
	Msg   string
}
type Success struct {
	Level string
	Msg   string
}

func NewResponse() *Response {
	return &Response{Status: RS.RS_success}
}
func (resp *Response) Tips(level string, rs int) {
	resp.Err = Error{level, "code=" + fmt.Sprint(rs) + "|" + RS.Desc(rs)}
}
func (resp *Response) WriteJson(w http.ResponseWriter) {
	b, err := json.Marshal(resp)
	if err != nil {
		w.Write([]byte(`{Status:-1,Err:Error{Level:"alert",Msg:"code=-1|序列化失败！"}}`))
	} else {
		w.Write(b)
	}
}
func (resp *Response) Success() {
	resp.Status = RS.RS_success
	resp.Data = Success{Level: SUCCESS, Msg: "恭喜(●'◡'●)|操作成功。"}
}

// -------------------------- Node --------------------------
type Tostring interface {
	String() string
}

type Node struct {
	Type     string
	Class    string
	Extra    string
	Text     string
	Children []*Node
}

func (n *Node) String() string {
	html := "<" + n.Type
	if n.Class != "" {
		html += " class='" + n.Class + "'"
	}
	if n.Extra != "" {
		html += " " + n.Extra
	}
	html += ">"
	html += n.Text
	if len(n.Children) > 0 {
		for _, child := range n.Children {
			html += child.String()
		}
	}
	html += "</" + n.Type + ">"
	return html
}

type Group struct {
	Data interface{}
	Page int
}

// -------------------------- Entrypassword --------------------------
const (
	SALT = "$^*#,.><)(_+f*m"
)

// rand salt
func RandSalt() string {
	var salt = ""
	for i := 0; i < 4; i++ {
		rand := GetRand()
		salt += string(SALT[rand.Intn(len(SALT))])
	}
	return salt
}

// encrypt password
func EncryptPasswd(name, pass, salt string) string {
	salt1 := "%$@w"
	h := md5.New()
	io.WriteString(h, salt1)
	io.WriteString(h, name)
	io.WriteString(h, salt)
	io.WriteString(h, pass)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func VerifyPasswd(passwd, name, pass, salt string) bool {
	return passwd == EncryptPasswd(name, pass, salt)
}

// randseed
func GetRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// -------------------------- Entrypdata --------------------------
func DesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	crypted := make([]byte, len(origData))
	blockMode := cipher.NewCBCEncrypter(block, key)
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func DesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData := make([]byte, len(crypted))
	blockMode := cipher.NewCBCDecrypter(block, key)
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
