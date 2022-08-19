package router

import (
	"net/http"

	"gotribe/compiler/lib"

	"github.com/gin-gonic/gin"
)

// RunController 脚本运行API
type RunController struct {
	*BaseController
}

var Run = &RunController{}

// RequestExecParams 编译器参数校验
type RequestExecParams struct {
	Code  string `form:"code" json:"code" binding:"required"`
	Input string `form:"input" json:"input"`
	Lang  string `form:"lang" json:"lang" binding:"required"`
}

// Exec post接口
func (p *RunController) Exec(c *gin.Context) {

	params := RequestExecParams{}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, Error(err.Error(), 1001, false))
		return
	}

	if len(params.Code) > 1024*400 {
		c.JSON(http.StatusBadRequest, Error("提交的代码太长，最多允许400KB", 1002, false))
		return
	}
	langexists, _ := lib.LangExists(params.Lang)
	if !langexists {
		c.JSON(http.StatusBadRequest, Error("暂时不支持该语言", 1002, false))
		return
	}
	tpl := lib.Run(params.Lang)
	output := lib.DockerRun(tpl.Image, params.Code, tpl.File, tpl.Cmd, tpl.Timeout, tpl.Memory)
	// 返回数据
	data := make(map[string]string)
	data["stdout"] = output
	data["stderr"] = ""

	c.JSON(http.StatusOK, Success("success", data))
}
