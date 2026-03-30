package user

import (
	"forum-gateway/util"
	"strconv"

	. "forum-gateway/handler"
	"forum/log"
	"forum/model"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetRole ... 设置权限
// @Summary 设置权限
// @Description 前端无需接此接口，casbin只会在启动时读取一遍数据库，因此改权限不能直接修改数据库，故写一个接口
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Success 200 {object} handler.Response
// @Router /auth/set_role/{id} [post]
func SetRole(c *gin.Context) {
	log.Info("User SetRole function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	userId := c.MustGet("userId").(uint32)

	targetId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, nil, err.Error(), GetLine())
		return
	}

	var req AddRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	// 属于该角色才能赋予别人该角色
	ok, err := model.HasRole(userId, req.Role)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "", GetLine())
		return
	}

	err = model.AddRole("user", uint32(targetId), req.Role)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	SendResponse(c, nil, nil)
}
