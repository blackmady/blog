package model

import (
	"blog/conf"
	"log"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"xorm.io/xorm"

	// 数据库驱动
	_ "github.com/go-sql-driver/mysql"
)

// Db 数据库操作句柄
var Db *xorm.Engine

func Init() {
	// 初始化数据库操作的 Xorm
	db, err := xorm.NewEngine("mysql", conf.Dsn.Dsn())
	if err != nil {
		log.Fatalln("数据库 dsn:", err.Error())
	}
	if err = db.Ping(); err != nil {
		log.Fatalln("数据库 ping:", err.Error())
	}
	db.SetMaxIdleConns(conf.Xorm.Idle)
	db.SetMaxOpenConns(conf.Xorm.Open)
	// 是否显示sql执行的语句
	db.ShowSQL(conf.Xorm.Show)
	db.ShowExecTime(conf.Xorm.Show)
	if conf.Xorm.Cache.Enable {
		// 设置xorm缓存
		cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), conf.Xorm.Cache.Count)
		db.SetDefaultCacher(cacher)
	}
	if conf.Xorm.Sync {
		err := db.Sync2(new(User), new(Cate), new(Tag), new(Post), new(PostTag), new(Opts))
		if err != nil {
			log.Fatalln("数据库 sync:", err.Error())
		}
	}
	Db = db
	//缓存
	initMap()
}

// Page 分页基本数据
type Page struct {
	Pi   int    `json:"pi" form:"pi" query:"pi"`       //分页页码
	Ps   int    `json:"ps" form:"ps" query:"ps"`       //分页大小
	Mult string `json:"mult" form:"mult" query:"mult"` //多条件信息
}

// Trim 去除空白字符
func (p *Page) Trim() string {
	p.Mult = strings.TrimSpace(p.Mult)
	return p.Mult
}

// JwtClaims jwt
type JwtClaims struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Num  string `json:"num"`
	Role Role   `json:"role"`
	jwt.StandardClaims
}

// Naver 上下页
type Naver struct {
	Prev string
	Next string
}

// State 统计信息
type State struct {
	Post int `json:"post"`
	Page int `json:"page"`
	Cate int `json:"cate"`
	Tag  int `json:"tag"`
}

// Collect 统计信息
func Collect() (*State, bool) {
	mod := &State{}
	has, _ := Db.SQL(`SELECT * FROM(SELECT COUNT(id) as post FROM post WHERE type=0)as a ,(SELECT COUNT(id) as page FROM post WHERE type=1) as b, (SELECT COUNT(id) as cate FROM cate) as c, (SELECT COUNT(id) as tag FROM tag) as d`).Get(mod)
	return mod, has
}
