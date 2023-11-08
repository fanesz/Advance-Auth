package handler

import (
	"advanceauth/backend/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
)

func BindAndStructValidator(c *gin.Context, req interface{}, structValidate bool) error {
	if err := c.ShouldBindJSON(req); err != nil {
		Error(c, http.StatusBadRequest, utils.REQ_WRONG_BODY_FORMAT, err.Error())
		return err
	}

	if structValidate {
		val := validator.New()
		if err := val.Struct(req); err != nil {
			Error(c, http.StatusBadRequest, utils.REQ_FIELD_ERROR, err.Error())
			return err
		}
	}

	return nil
}

func QueryValidator(query *gorm.DB, c *gin.Context, count bool) (bool, int64) {
	if query.Error != nil {
		Error(c, http.StatusInternalServerError, utils.DB_QUERY_ERROR, query.Error.Error())
		return false, -1
	}
	if !count {
		return true, -1
	}
	var result int64
	if query.Count(&result); result == 0 {
		return false, 0
	}
	return true, result
}
