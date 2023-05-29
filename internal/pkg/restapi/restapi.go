package restapi

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/model"
	"github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/service"
)

type server struct {
	userService service.User
}

func NewServer(factory service.Factory) *server {
	userService, _ := factory.CreateUserService()

	return &server{userService: userService}
}

func (srv *server) getIndex(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello")
}

func (srv *server) postUser(c *gin.Context) {
	var createUserRequest model.CreateUserRequest

	if err := c.BindJSON(&createUserRequest); err != nil {
		log.Printf("Binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := createUserRequest.Validate(); err != nil {
		log.Printf("Validation error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	createUserResponse, err := srv.userService.CreateUser(createUserRequest)
	if err != nil {
		log.Printf("CreateUser error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, createUserResponse)
}

func (srv *server) postUserLogin(c *gin.Context) {
	var loginUserRequest model.LoginUserRequest

	if err := c.BindJSON(&loginUserRequest); err != nil {
		log.Printf("Binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	loginUserResponse, err := srv.userService.LoginUser(loginUserRequest)
	if err != nil {
		log.Printf("LoginUser error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, loginUserResponse)
}

func setupRouter() *gin.Engine {

	factory := service.Factory(&service.FactoryInMemory{})
	srv := NewServer(factory)

	router := gin.Default()

	router.GET("/", srv.getIndex)
	router.POST("/v1/user", srv.postUser)
	router.POST("/v1/user/login", srv.postUserLogin)

	return router
}

func Run(addr string) error {
	router := setupRouter()
	return router.Run(addr)
}
