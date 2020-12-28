package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func Conn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "Hacker"
	dbName := "article"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		log.Panicln("Error in DB Connection ", err.Error())
	}
	return db
}

type Article struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
}

func HomePage(c *gin.Context) {
	c.JSON(200, gin.H{
		"Project Name": "Article Management System",
	})
}

func CreateArticle(c *gin.Context) {
	// Initlize Database
	db := Conn()
	defer db.Close()
	var data Article
	//Reading Request Body
	RequestBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Invalid Parameter"})
		return
	}
	//Convert Request Body into Json Formate
	json.Unmarshal(RequestBody, &data)
	log.Println("Data : ", &data)
	result, err := db.ExecContext(c, "insert into articledata (title,description,author) values (?,?,?)", data.Title, data.Description, data.Author)
	result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Invalid Parameter"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{

		"message": "article saved",
	})
	return
}

func UpdateArticle(c *gin.Context) {
	// Initlize Database
	db := Conn()
	defer db.Close()
	id := c.Param("id")
	var data Article
	//Reading Request Body
	RequestBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Invalid Parameter"})
		return
	}
	//Convert Request Body into Json Formate
	json.Unmarshal(RequestBody, &data)
	log.Println("Update Data : ", &data)
	result, err := db.ExecContext(c, "update articledata set title = ? , description = ? , author = ? where id = ? ", data.Title, data.Description, data.Author, id)
	rows, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Invalid Parameter"})
		return
	}
	if rows != 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not Found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{

		"message": "updated article",
	})
	return
}

func DeleteArticle(c *gin.Context) {
	// Initlize Database
	db := Conn()
	defer db.Close()
	// Store Id
	id := c.Param("id")
	log.Println("Id is : ", id)
	result, err := db.ExecContext(c, "delete from articledata where id = ?", id)
	rows, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Invalid Parameter"})
		return
	}
	if rows != 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not Found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{

		"message": "article deleted",
	})
	return
}

func AllArticle(c *gin.Context) {
	var article Article
	// Initlize Database
	db := Conn()
	defer db.Close()
	result, err := db.Query("select *from articledata")
	defer result.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Server error"})
		return
	}
	for result.Next() {
		err := result.Scan(&article.Id, &article.Title, &article.Description, &article.Author)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Server error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": article,
		})
	}
	return
}

func SingleArticle(c *gin.Context) {
	// Initlize Database
	db := Conn()
	defer db.Close()
	// Store Id
	id := c.Param("id")
	log.Println("Id is : ", id)
	result, err := db.Query("select *from articledata where id = ?", id)
	defer result.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Server error"})
		return
	}
	for result.Next() {
		var article Article
		err := result.Scan(&article.Id, &article.Title, &article.Description, &article.Author)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Server error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": article,
		})
	}
	return
}

func RequestHandler() {
	r := gin.Default()
	r.GET("/", HomePage)
	r.GET("/all", AllArticle)
	r.GET("/article/:id", SingleArticle)
	r.POST("/create", CreateArticle)
	r.PUT("/update/:id", UpdateArticle)
	r.DELETE("/delete/:id", DeleteArticle)
	r.Run()
}

func main() {
	Conn()
	RequestHandler()
}
