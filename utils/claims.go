package utils

import (
	"errors"
	"net"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ClearToken(c *gin.Context) {
	// 增加cookie x-token 向来源的web添加
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}

	if net.ParseIP(host) != nil {
		c.SetCookie("x-token", "", -1, "/", "", false, false)
	} else {
		c.SetCookie("x-token", "", -1, "/", host, false, false)
	}
}

func SetToken(c *gin.Context, token string, maxAge int) {
	// 增加cookie x-token 向来源的web添加

	//从请求的主机信息中分离主机信息和端口号
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil { //不包含端口号，那么c.Request.Host就是主机信息，直接赋值给host即可
		host = c.Request.Host
	}

	if net.ParseIP(host) != nil { //判断字符串是否是有效的IP地址，是返回IP地址
		//maxAge：Cookie的最大存活时间，单位为秒 path："/" 表明cookie对哪些请求路径有效
		c.SetCookie("x-token", token, maxAge, "/", "", false, false)
	} else {
		c.SetCookie("x-token", token, maxAge, "/", host, false, false)
	}
}

// GetToken reWrite  这个用于中间件访问获取token，并进行一些初始化，主要目的是只运行一次
// 获取token，解析信息，存入上下文中  相较于源代码，这次直接全部返回去
func GetToken(c *gin.Context) (string, error) {
	token := GetTokenSingle(c)
	if token == "" {
		global.GVA_LOG.Error("未登录或非法访问")
		return "", errors.New("未登录或非法访问")
	}
	claims, err := GetClaims(c)
	if err != nil {
		global.GVA_LOG.Error("重新写入cookie token失败")
		return "", err
	} else {
		// 将获取的claim存入上下文中备用
		c.Set("claims", claims)
		//重新设置token
		SetToken(c, token, int((claims.ExpiresAt.Unix()-time.Now().Unix())/60))
	}
	return token, nil
}

func GetTokenSingle(c *gin.Context) string {
	token, _ := c.Cookie("x-token")
	if token == "" {
		token = c.Request.Header.Get("x-token")
	}
	return token
}

func GetClaims(c *gin.Context) (*systemReq.CustomClaims, error) {
	//判断上下文中是否有Claims
	if claims, exists := c.Get("claims"); exists {
		return claims.(*systemReq.CustomClaims), nil
	}
	//没有，获取token解析返回claims
	token := GetTokenSingle(c)
	j := NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		global.GVA_LOG.Error("jwt解析信息失败, 请检查请求头是否存在x-token且claims是否为规定结构")
	}
	return claims, err
}

// GetUserID 从Gin的Context中获取从jwt解析出来的用户ID
func GetUserID(c *gin.Context) uint {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.BaseClaims.ID
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.BaseClaims.ID
	}
}

// GetUserUuid 从Gin的Context中获取从jwt解析出来的用户UUID
func GetUserUuid(c *gin.Context) uuid.UUID {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return uuid.UUID{}
		} else {
			return cl.UUID
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.UUID
	}
}

// GetUserAuthorityId 从Gin的Context中获取从jwt解析出来的用户角色id
func GetUserAuthorityId(c *gin.Context) uint {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.AuthorityId
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.AuthorityId
	}
}

// GetUserInfo 从Gin的Context中获取从jwt解析出来的用户角色id
func GetUserInfo(c *gin.Context) *systemReq.CustomClaims {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return nil
		} else {
			return cl
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse
	}
}

// GetUserName 从Gin的Context中获取从jwt解析出来的用户名
func GetUserName(c *gin.Context) string {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return ""
		} else {
			return cl.Username
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.Username
	}
}

func LoginToken(user system.Login) (token string, claims systemReq.CustomClaims, err error) {
	j := NewJWT()
	claims = j.CreateClaims(systemReq.BaseClaims{
		UUID:        user.GetUUID(),
		ID:          user.GetUserId(),
		NickName:    user.GetNickname(),
		Username:    user.GetUsername(),
		AuthorityId: user.GetAuthorityId(),
	})
	token, err = j.CreateToken(claims)
	return
}
