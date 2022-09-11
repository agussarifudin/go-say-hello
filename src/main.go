package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// func ById(c *gin.Context){
// 	c.JSON(200,gin.H{
// 		"message":"Hello World",
// 	})}


func PostHomePage(c *gin.Context){
	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err.Error())
	}
	c.JSON(200,gin.H{
		"message":string(value),
	})
}

func QueryStrings(c *gin.Context){
	name := c.Query("name") // /query?name=test&age=23
	age := c.Query("age")

	c.JSON(200,gin.H{
		"name":name,
		"age":age,
	})
}
func PathParameters(c *gin.Context){
	name := c.Param("name")
	age := c.Param("age")

	c.JSON(200, gin.H{
		"name":name,
		"age":age,
	})
}



func main(){

	db,err := sql.Open("mysql","root:Kmzwa8awaaas@tcp(127.0.0.1:3306)/golang")

	err = db.Ping()
	if err != nil{
		panic("Gagal ke database")
	}
	defer db.Close()
	fmt.Println("hello world")

	r := gin.Default()


	type Users struct{
		Id int `json:"id"`
		User_nik string `json:"user_nik"`
		User_name string `json:"user_name"`
		User_address string `json:"user_address"`
	}

	r.GET("/:id",func(c *gin.Context){
		var(
			users Users
			result gin.H
		)
		id := c.Param("id")
		row := db.QueryRow("select id,user_nik,user_name,user_address from users where id = ?;",id)
		err = row.Scan(&users.Id,&users.User_nik,&users.User_name,&users.User_address)

		if err != nil {
			result = gin.H{
				"hasil":"tidak ada",
				"Jumlah":0,}
			}else{
				result =gin.H{
					"HASIL":users,
					"jumlah":1,
				}
			}
			c.JSON(http.StatusOK,result)
		
	})

	r.GET("/",func (c *gin.Context){
		var(
			users Users
			userSlice []Users
		)
		rows , err := db.Query("select * from users ;")
		if err != nil{
			fmt.Print(err.Error())
		}
		for rows.Next(){
			err = rows.Scan(&users.Id,&users.User_name,&users.User_nik,&users.User_address)
			userSlice = append(userSlice,users)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK,gin.H{
			"hasil":userSlice,
			"jmlh":len(userSlice),
		})
	})

	r.POST("/post",func (c *gin.Context){
		var buffer bytes.Buffer
		id := c.PostForm("id")
		user_name := c.PostForm("user_name")
		user_nik := c.PostForm("user_nik")
		user_address := c.PostForm("user_address")
		stmt , err := db.Prepare("insert into users (id,user_nik,user_name,user_address) values (?,?,?,?);")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(id,user_nik,user_name,user_address)

		if err != nil {
			fmt.Print(err.Error())
		}

		buffer.WriteString(user_name)
		buffer.WriteString(" ")
		buffer.WriteString(user_nik)
		defer stmt.Close()

		data := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"pesan ": fmt.Sprintf(" berhasil menambahkan user %s",data),
		})
	})

	r.PUT("/",func (c *gin.Context){
		var buffer bytes.Buffer
		id := c.PostForm("id")
		user_nik := c.PostForm("user_nik")
		user_name := c.PostForm("user_name")
		user_address := c.PostForm("user_address")
		stmt , err := db.Prepare("update users set user_nik = ?, user_name=?,user_address=? where id = ?;")

		if err != nil {
			fmt.Print(err.Error())
		}
		_,err = stmt.Exec(user_nik,user_name,user_address,id)
		if err != nil{
			fmt.Print(err.Error())
		}

		buffer.WriteString(user_nik)
		buffer.WriteString(" ")
		buffer.WriteString(user_name)
		buffer.WriteString(" ")
		buffer.WriteString(user_address)
		defer stmt.Close()

		data := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"pesan ": fmt.Sprintf("Berhasil merubah menjadi %s",data),
		})
	})

	r.DELETE("/",func (c *gin.Context){

		id := c.PostForm("id")
		query , err := db.Prepare("delete from users where id = ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_,err = query.Exec(id)
		if err != nil{
			fmt.Print(err.Error())
		}
		
		
		defer query.Close()
	
		c.JSON(http.StatusOK, gin.H{
			"pesan":fmt.Sprintf("Berhasil menghapus data"),
		})
	})

	r.GET("/query",QueryStrings)
	r.POST("/",PostHomePage)
	r.GET("/path/:name/:age",PathParameters)
	r.Run()
}


