package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func initRouter() {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	//e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
	//	Skipper:    skipper,
	//	SigningKey: []byte("goblog"),
	//}))

	//注册登录
	{
		e.GET("/", func(c echo.Context) error {
			return c.File("html/index.html")
		})
		e.GET("/register", func(c echo.Context) error {
			return c.File("html/register.html")
		})
		e.GET("/login", func(c echo.Context) error {
			return c.File("html/login.html")
		})
		e.POST("/register", postRegister)

		e.POST("/login", postLogin)
	}
	//完善个人信息 restful
	{

		e.GET("/infermation", func(c echo.Context) error {

			return c.File("html/infermation.html")
		})

		e.PUT("/users/:id", putInfermation)
		//获取验证码
		e.PUT("/users/passwd", putPasswd)
		//找回密码
		e.PUT("/users/pd", putPasswd1)

	}
	//文章接口
	ga := e.Group("/article", middleware.JWT("123456"))
	{
		ga.GET("", func(c echo.Context) error {
			return nil
		})
	}
	gu := e.Group("/user")
	{
		gu.GET("", func(c echo.Context) error {
			return nil
		})
	}
	go func() {
		if err := e.Start(":1323"); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func skipper(c echo.Context) bool {
	if c.Request().Method == http.MethodOptions {
		return true
	}
	if c.Path() == "/register" || c.Path() == "/login" || c.Path() == "/infermation" || c.Path() == "/" {
		return true
	}
	return false
}

func postRegister(c echo.Context) error {
	u := new(User)
	resp := new(Response)
	if err := c.Bind(u); err != nil {
		resp.Error = 1
		resp.Msg = "参数格式错误"
		resp.Data = err.Error()
		logrus.Error(err)
		return c.JSON(400, resp)
	}
	logrus.Infof("%+v\n", u)
	if err := postI(u); err != nil {
		resp.Error = 1
		resp.Msg = "注册失败"
		resp.Data = err.Error()
		logrus.Error(err)
		return c.JSON(400, resp)
	}
	resp.Error = 0
	resp.Msg = "注册成功"
	resp.Data = u
	return c.JSON(200, resp)
}
func postLogin(c echo.Context) error {
	u := new(User)
	resp := new(Response)
	if err := c.Bind(u); err != nil {
		resp.Error = 1
		resp.Msg = "参数错误"
		resp.Data = err.Error()
		logrus.Error(err)
		return c.JSON(400, resp)
	}
	logrus.Infof("%+v\n", u)
	if err := u.login(); err != nil {
		resp.Error = 1
		resp.Msg = "用户名/密码错误"
		resp.Data = err.Error()
		logrus.Error(err)
		return c.JSON(400, resp)
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = u.Username
	claims["id"] = u.ID
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(7*24)).Unix()
	t, err := token.SignedString([]byte("goblog"))
	if err != nil {
		resp.Error = 1
		resp.Msg = "token生成失败"
		resp.Data = err.Error()
		logrus.Error(err)
		return c.JSON(500, resp)
	}
	resp.Error = 0
	resp.Msg = "登录成功"
	resp.Data = t
	return c.JSON(200, resp)
}

func putInfermation(c echo.Context) error {
	u := new(User)
	resp := new(Response)
	id := c.Param("id")
	u.ID, _ = strconv.Atoi(id)
	logrus.Infof("%+v\n", u)
	if err := c.Bind(u); err != nil {
		resp.Error = 1
		resp.Msg = "参数格式错误"
		resp.Data = err.Error()
		logrus.Error(err)
		return c.JSON(400, resp)
	}
	logrus.Infof("%+v\n", u)
	if err := putI(u); err != nil {
		resp.Error = 1
		resp.Msg = "更新失败"
		resp.Data = err.Error()
		logrus.Error(err)
		return c.JSON(400, resp)
	}
	resp.Error = 0
	resp.Msg = "更新成功"
	resp.Data = u
	return c.JSON(200, resp)
}

func putPasswd(c echo.Context) error {
	cd := genCode()
	u := User{}
	sendEmail(cd, u.Email)
	mc[u.Email] = cd
	return c.JSON(200, 0)
}

func putPasswd1(c echo.Context) error {
	u := User{}
	cd := c.Param("code")
	v, ok := mc[u.Email]
	if !ok {
		return c.JSON(400, "验证码有误")
	}
	if v == cd {
		//继续让他修改密码
	}
	return c.JSON(200, 0)
}
