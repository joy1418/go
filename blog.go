package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Blog struct {
	Id         int    `json:"id" form:"id"`
	Title      string `json:"title" form:"title"`
	Author     string `json:"author" form:"author"`
	Content    string `json:"content" form:"content"`
	CreateTime string `json:"create_time" form:"create_time"`
}

//发表新博客方法
func (b *Blog) AddBlog() (id int64, err error) {
	rs, err := db.Exec("INSERT INTO blog(title,author,content,create_time) VALUES(?,?,?,?)", b.Title, b.Author, b.Content, b.CreateTime)
	if err != nil {
		return
	}
	id, err = rs.LastInsertId()
	return
}

//获取文章列表方法
func (b *Blog) GetBlogs() (blogs []Blog, err error) {
	blogs = make([]Blog, 0)
	rows, err := db.Query("SELECT * FROM blog")
	defer rows.Close()

	if err != nil {
		return
	}
	for rows.Next() {
		var blog Blog
		rows.Scan(&blog.Id, &blog.Title, &blog.Author, &blog.Content, &blog.CreateTime)
		blogs = append(blogs, blog)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

//获取某篇文章
func (b *Blog) GetBlog() (blog Blog, err error) {
	err = db.QueryRow("SELECT * FROM blog WHERE id=?", b.Id).Scan(&blog.Id, &blog.Title, &blog.Author, &blog.Content, &blog.CreateTime)
	return
}

//修改某篇文章
func (b *Blog) UpdBlog() (ra int64, err error) {
	stmt, err := db.Prepare("UPDATE blog SET title=?,author=?,content=?,create_time=?")
	defer stmt.Close()
	if err != nil {
		return
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	rs, err := stmt.Exec(b.Title, b.Author, b.Content, t)
	if err != nil {
		return
	}
	ra, err = rs.RowsAffected()
	return
}

//删除某篇文章
func (b *Blog) DelBlog() (ra int64, err error) {
	rs, err := db.Exec("DELETE FROM blog WHERE id=?", b.Id)
	if err != nil {
		return
	}
	ra, err = rs.RowsAffected()
	return
}

//handler函数
func IndexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "blog",
	})
}

func AddBlogHandler(c *gin.Context) {
	title := c.Request.FormValue("title")
	author := c.Request.FormValue("author")
	content := c.Request.FormValue("content")
	create_time := time.Now().Format("2006-01-02 15:04:05")
	b := Blog{Title: title, Author: author, Content: content, CreateTime: create_time}
	ra, err := b.AddBlog()
	if err != nil {
		log.Fatalln(err)
	}
	msg := fmt.Sprintf("insert successful %d", ra)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}

func GetBlogsHandler(c *gin.Context) {
	var b Blog
	blogs, err := b.GetBlogs()
	if err != nil {
		log.Fatalln(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"blogs": blogs,
	})
}

func GetBlogHandler(c *gin.Context) {
	cid := c.Param("id")
	id, err := strconv.Atoi(cid)
	if err != nil {
		log.Fatalln(err)
	}
	b := Blog{Id: id}
	blog, err := b.GetBlog()
	if err != nil {
		log.Fatalln(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"blog": blog,
	})
}

func UpdBlogHandler(c *gin.Context) {
	cid := c.Param("id")
	id, err := strconv.Atoi(cid)
	if err != nil {
		log.Fatalln(err)
	}
	b := Blog{Id: id}
	err = c.Bind(&b)
	if err != nil {
		log.Fatalln(err)
	}
	ra, err := b.UpdBlog()
	if err != nil {
		log.Fatalln(err)
	}
	msg := fmt.Sprintf("Update blog %d successful %d", b.Id, ra)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}

func DelBlogHandler(c *gin.Context) {
	cid := c.Param("id")
	id, err := strconv.Atoi(cid)
	if err != nil {
		log.Fatalln(err)
	}
	b := Blog{Id: id}
	ra, err := b.DelBlog()
	if err != nil {
		log.Fatalln(err)
	}
	msg := fmt.Sprintf("delete %d successful %d", b.Id, ra)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/go?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}
	router := gin.Default()

	router.GET("/", IndexHandler)
	router.POST("/blog", AddBlogHandler)
	router.GET("/blogs", GetBlogsHandler)
	router.GET("/blog/:id", GetBlogHandler)
	router.PUT("/blog/:id", UpdBlogHandler)
	router.DELETE("/blog/:id", DelBlogHandler)
	router.Run()

}
