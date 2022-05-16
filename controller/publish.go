package controller

import (
	"fmt"
	"net/http"
	"path/filepath"
	"github.com/gin-gonic/gin"
	"douyin/message"
	"douyin/service"
)

type VideoListResponse struct {
	message.Response
	VideoList []message.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.Query("token")
	// check user token
	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, message.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	// get uploaded video file
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, message.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	// get video address, save video into address
	filename := filepath.Base(data.Filename)
	var user message.User = usersLoginInfo[token]
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, message.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	// write video message into database
	service.SaveVideoMessage(saveFile, user.Id)

	c.JSON(http.StatusOK, message.Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	// check user token
	token := c.Query("token")
	user, exist := usersLoginInfo[token]
	if !exist {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: message.Response{StatusCode: 1,
				StatusMsg: "User doesn't exist",
			},
		})
		return
	}
	//get user video list
	videolist, err := service.GetUserVideoList(user.Id)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: message.Response{StatusCode: 1,
				StatusMsg: "Error occers while getting video list",
			},
		})
		return
	}

	//return user video list
	c.JSON(http.StatusOK, VideoListResponse{
		Response: message.Response{
			StatusCode: 0,
		},
		VideoList: videolist,
	})
}
