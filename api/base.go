package api

import "github.com/gin-gonic/gin"

type APIInstance interface {
	Logic(c *gin.Context) Result
}

type DataList struct {
	Count int64       `json:"count"`
	List  interface{} `json:"list"`
}

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func FORM(c *gin.Context, val APIInstance) {
	if err := c.ShouldBind(val); err != nil {
		c.JSON(200, err.Error())
		return
	}
	c.JSON(200, val.Logic(c))
}

func JSON(c *gin.Context, val APIInstance) {
	if err := c.ShouldBindJSON(val); err != nil {
		c.JSON(200, err.Error())
		return
	}
	c.JSON(200, val.Logic(c))
}

func Header(c *gin.Context, val APIInstance) {
	if err := c.ShouldBindHeader(val); err != nil {
		c.JSON(200, err.Error())
		return
	}
	c.JSON(200, val.Logic(c))
}

func XML(c *gin.Context, val APIInstance) {
	if err := c.ShouldBindXML(val); err != nil {
		c.JSON(200, err.Error())
		return
	}
	c.JSON(200, val.Logic(c))
}

func URL(c *gin.Context, val APIInstance) {
	if err := c.ShouldBindUri(val); err != nil {
		c.JSON(200, err.Error())
		return
	}
	c.JSON(200, val.Logic(c))
}
