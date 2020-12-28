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
	dbName := "product"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		log.Panicln("Error in DB Connection ", err.Error())
	}
	return db
}

type Product struct {
	Id                  int64  `json:"id"`
	Product_name        string `json:"product_name"`
	Storage_center      string `json:"storage_center"`
	Product_description string `json:"product_description"`
}

func HomePage(c *gin.Context) {
	c.JSON(200, gin.H{
		"Project Name": "Product Storage Center",
	})
}

func CreateProduct(c *gin.Context) {
	// Initlize Database
	db := Conn()
	defer db.Close()
	var data Product
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
	result, err := db.ExecContext(c, "insert into productdata (product_name,storage_center,product_description) values (?,?,?)", data.Product_name, data.Storage_center, data.Product_description)
	result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Invalid Parameter"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{

		"message": "product saved",
	})
	return
}

func UpdateProduct(c *gin.Context) {
	// Initlize Database
	db := Conn()
	defer db.Close()
	id := c.Param("id")
	var data Product
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
	result, err := db.ExecContext(c, "update productdata set product_name = ? , storage_center = ? , product_description = ? where id = ? ", data.Product_name, data.Storage_center, data.Product_description, id)
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

		"message": "updated product",
	})
	return
}

func DeleteProduct(c *gin.Context) {
	// Initlize Database
	db := Conn()
	defer db.Close()
	// Store Id
	id := c.Param("id")
	log.Println("Id is : ", id)
	result, err := db.ExecContext(c, "delete from productdata where id = ?", id)
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

		"message": "product deleted",
	})
	return
}

func AllProduct(c *gin.Context) {
	var product Product
	// Initlize Database
	db := Conn()
	defer db.Close()
	result, err := db.Query("select *from productdata")
	defer result.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Server error"})
		return
	}
	for result.Next() {
		err := result.Scan(&product.Id, &product.Product_name, &product.Product_description, &product.Product_description)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Server error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": product,
		})
	}
	return
}

func SingleProduct(c *gin.Context) {
	// Initlize Database
	db := Conn()
	defer db.Close()
	// Store Id
	id := c.Param("id")
	log.Println("Id is : ", id)
	result, err := db.Query("select *from productdata where id = ?", id)
	defer result.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Server error"})
		return
	}
	for result.Next() {
		var product Product
		err := result.Scan(&product.Id, &product.Product_name, &product.Product_description, &product.Product_description)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Server error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": product,
		})
	}
	return
}

func RequestHandler() {
	r := gin.Default()
	r.GET("/", HomePage)
	r.GET("/all", AllProduct)
	r.GET("/product/:id", SingleProduct)
	r.POST("/create", CreateProduct)
	r.PUT("/update/:id", UpdateProduct)
	r.DELETE("/delete/:id", DeleteProduct)
	r.Run()
}

func main() {
	Conn()
	RequestHandler()
}
