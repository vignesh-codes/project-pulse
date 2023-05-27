package utils

import (
	"fmt"
	"net/http"
	"pulse-service/logger"
	"pulse-service/utils/response"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

var node *snowflake.Node

func InitSnowflakeNode() error {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		logger.Logger.Fatal("Error creating Snowflake node:", zap.Any(logger.KEY_ERROR, err))
		return err
	}
	return nil
}

func GenerateUUID(str string) int64 {
	return node.Generate().Int64()
}

func BindJSON(c *gin.Context, req interface{}) bool {
	if c.ContentType() != "application/json" {
		msg := fmt.Sprintf("%s only accepts Content-Type application/json", c.FullPath())
		reply := response.UnsupportedMediaType(msg)
		c.JSON(reply.Status(), reply)
		return false
	}

	if err := c.ShouldBind(req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {

			var errList []string
			for _, err := range errs {
				e := fmt.Sprintf("field=%s tag=%s required=%s kind=%s value=%v", err.Field(), err.Tag(), err.Param(), err.Kind(), err.Value())
				errList = append(errList, e)
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"status_code": http.StatusBadRequest,
				"message":     "Validation errors",
				"errors":      errList,
			})
			return false
		}

		if err.Error() == "EOF" {
			status := response.BadRequest("Empty Body")
			c.JSON(status.Status(), status)
			return false
		}

		fallback := response.InternalServerError(logger.BindJSONtoStruct, "BindJSON::fallback", err)
		c.JSON(fallback.Status(), fallback)
		return false
	}

	return true
}
