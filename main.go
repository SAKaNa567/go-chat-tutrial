package main

import "github.com/gin-gonic/gin"

func main() {
    
    //routing 
    m := gin.Default()
    m.POST("/register",controller.UserAddController)
    m.Run(":8080")
}
