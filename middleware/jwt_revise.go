package middleware

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

var jwtServ = service.ServiceGroupApp.SystemServiceGroup.JwtService

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取token  同时存储claim到上下文中
		token, err := utils.GetToken(c)
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenExpired):
				response.NoAuth(utils.TokenExpired.Error(), c)
			case errors.Is(err, jwt.ErrTokenMalformed):
				response.NoAuth(utils.TokenMalformed.Error(), c)
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				response.NoAuth(utils.TokenSignatureInvalid.Error(), c)
			case errors.Is(err, jwt.ErrTokenNotValidYet):
				response.NoAuth(utils.TokenNotValidYet.Error(), c)
			case errors.Is(err, utils.TokenInvalid):
				response.NoAuth(utils.TokenInvalid.Error(), c)
			default:
				response.NoAuth(err.Error(), c)
			}
			c.Abort()
			return
		}
		//判断token是否在黑名单
		isBlacklist := jwtServ.IsBlacklist(token)
		if isBlacklist {
			response.NoAuth("账户异地登录或令牌失效", c)
			//令牌不合格，将清除掉
			utils.ClearToken(c)
			c.Abort()
			return
		}

		// 已登录用户被管理员禁用 需要使该用户的jwt失效 此处比较消耗性能 如果需要 请自行打开
		// 用户被删除的逻辑 需要优化 此处比较消耗性能 如果需要 请自行打开

		//if user, err := userService.FindUserByUuid(claims.UUID.String()); err != nil || user.Enable == 2 {
		//	_ = jwtService.JsonInBlacklist(system.JwtBlacklist{Jwt: token})
		//	response.FailWithDetailed(gin.H{"reload": true}, err.Error(), c)
		//	c.Abort()

		claims, err := utils.GetClaims(c)
		if claims.ExpiresAt.Unix()-time.Now().Unix() < claims.BufferTime { //判断剩余时间是否小于缓冲时间
			dr, _ := utils.ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime) //解析时间
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(dr))       //重新设置过期时间
			j := utils.NewJWT()
			newToken, _ := j.CreateTokenByOldToken(token, *claims) //通过老token创建新的token
			newClaims, _ := j.ParseToken(newToken)                 //解析获取新的claims
			c.Header("new-token", newToken)
			c.Header("new-expires-at", strconv.FormatInt(newClaims.ExpiresAt.Unix(), 10))
			utils.SetToken(c, newToken, int(dr.Seconds()))
			if global.GVA_CONFIG.System.UseMultipoint {
				// 记录新的活跃jwt
				_ = jwtService.SetRedisJWT(newToken, newClaims.Username)
			}
		}
		c.Next()

		if newToken, exists := c.Get("new-token"); exists {
			c.Header("new-token", newToken.(string))
		}
		if newExpiresAt, exists := c.Get("new-expires-at"); exists {
			c.Header("new-expires-at", newExpiresAt.(string))
		}

	}
}
