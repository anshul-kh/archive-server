package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExecRequestBody struct {
	Lang string `json:"lang"`
	Code string `json:"code"`
}

func executionHandler(ctx *gin.Context) {

	var body ExecRequestBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, ok := GlobalLangHandler.exts[body.Lang]
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": "false",
			"message": "Language Is Not Supported Yet..",
		})
		return
	}

	job_id, err := GlobalJobMapper.CreateJob(body.Lang, body.Code)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": "false",
			"message": "Error Creating A New Job",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"job_id":  job_id,
	})
}

func resultHandler(ctx *gin.Context) {

	job_id := ctx.Param("job_id")

	result, err := GlobalJobMapper.GetJobResult(job_id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": "false",
			"message": err.Error(),
		})
		return
	}

	var _type byte
	if len(result) > 0 {
		_type = []byte(result)[0]
	}

	switch _type {
	case SUCCESS:
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"result":  result[1:],
			"type":    "success",
		})
		return
	case ERROR:
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"result":  result[1:],
			"type":    "error",
		})
		return
	case PENDING:
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"result":  "(pending...)",
			"type":    "pending",
		})
		return
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"result":  "Internal Server Error",
			"type":    "error",
		})
	}

}
