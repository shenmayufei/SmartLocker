package v1

import (
	"SmartLocker/e"
	"SmartLocker/service/auth"
	"SmartLocker/service/face"
	"SmartLocker/service/task"
	"SmartLocker/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func RecognizeFace(c *gin.Context) {
	url := verifyImgAndSave(c)
	if url == "" {
		return
	}

	name := face.Recognize(url)
	go checkTask(name)

	c.JSON(http.StatusOK, Wrap(e.Success, name))
}

func checkTask(name string) {
	t, b := task.GetBackendTask(name)
	if !b {
		return
	}

	for _, i := range t {
		s := strings.Split(i, "-")
		if len(s) != 2 {
			continue
		}
		// id->position s[1]->cid
		id, _ := strconv.Atoi(s[0])
		task.AddClientTask(s[1], 0, id)
	}
	task.ConsumeBackendTask(name)
}

func RegisterFace(c *gin.Context) {
	token := c.PostForm("token")
	claim, err := auth.CheckToken(token)
	if err != e.Success {
		c.JSON(http.StatusOK, Wrap(err, nil))
		return
	}

	url := verifyImgAndSave(c)
	if url == "" {
		return
	}

	s := face.InsertImage(claim.Username, url)
	if !s {
		c.JSON(http.StatusOK, Wrap(e.RegistrationFailed, nil))
		return
	}

	c.JSON(http.StatusOK, Wrap(e.Success, nil))
}

func verifyImgAndSave(c *gin.Context) string {
	file, err := c.FormFile("img")
	if err != nil {
		c.JSON(http.StatusOK, Wrap(e.UploadFailed, nil))
		return ""
	}

	// TODO check content-type locally
	ext, errInt := util.VerifyImg(file.Header.Get("Content-Type"))
	if errInt != e.Success {
		c.JSON(http.StatusOK, Wrap(errInt, nil))
		return ""
	}

	filename := fmt.Sprintf("%s.%s",
		util.EncodeSha256(file.Filename+util.RandString(5)), ext)

	if err := c.SaveUploadedFile(file,
		fmt.Sprintf("./resources/uploads/%s", filename)); err != nil {
		c.JSON(http.StatusOK, Wrap(e.UploadFailed, nil))
		return ""
	}

	url := fmt.Sprintf("./resources/uploads/%s", filename)
	return url
}
