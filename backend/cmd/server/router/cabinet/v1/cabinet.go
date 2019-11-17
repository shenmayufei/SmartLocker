package v1

import (
	v1 "SmartLocker/cmd/server/router/http/v1"
	"SmartLocker/e"
	"SmartLocker/model"
	"SmartLocker/service/auth"
	"SmartLocker/service/cache"
	"SmartLocker/service/node"
	"SmartLocker/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func PingPong(c *gin.Context) {
	// make sure cid is a number
	cid, err := strconv.Atoi(c.PostForm("cid"))
	if err != nil {
		c.JSON(http.StatusOK, v1.Wrap(e.InvalidParams, nil))
		return
	}

	name, status := cache.NodePingPong(strconv.Itoa(cid))
	if !status {
		c.JSON(http.StatusOK, v1.Wrap(e.Unauthorized, nil))
		return
	}

	c.JSON(http.StatusOK, v1.Wrap(e.Success, gin.H{"ping": "pong", "name": name}))
}

func Register(c *gin.Context) {
	name := c.PostForm("name")
	location := c.PostForm("location")
	token := c.PostForm("token")
	if name == "" || token == "" || location == "" {
		c.JSON(http.StatusOK, v1.Wrap(e.InvalidParams, nil))
		return
	}
	lockerNum, err := strconv.Atoi(c.PostForm("num"))
	if err != nil {
		c.JSON(http.StatusOK, v1.Wrap(e.InvalidParams, nil))
		return
	}

	status := cache.CheckToken(name, token)
	if !status {
		c.JSON(http.StatusOK, v1.Wrap(e.Unauthorized, nil))
		return
	}

	cid, status := node.RegisterCabinet(name, location, lockerNum)
	if !status {
		c.JSON(http.StatusOK, v1.Wrap(e.InternalError, nil))
		return
	}
	c.JSON(http.StatusOK, v1.Wrap(e.Success, cid))
}

func GenerateToken(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.JSON(http.StatusOK, v1.Wrap(e.InvalidParams, nil))
		return
	}

	token := c.PostForm("token")
	claim, err := auth.CheckToken(token)
	if err != e.Success {
		c.JSON(http.StatusOK, v1.Wrap(err, nil))
		return
	}

	if claim.Role != model.ADMIN {
		c.JSON(http.StatusOK, v1.Wrap(e.PermissionDenied, nil))
		return
	}

	randToken := util.RandString(20)
	status := cache.GenerateToken(name, randToken)
	if !status {
		c.JSON(http.StatusOK, v1.Wrap(e.RedisError, nil))
		return
	}
	c.JSON(http.StatusOK, v1.Wrap(e.Success, token))
	return
}
