package user

import (
	"github.com/guidoxie/knife/cmd/demo/internal/dao/user"
	"github.com/guidoxie/knife/pkg/app/http"
	"github.com/guidoxie/knife/pkg/errcode"
	"github.com/guidoxie/knife/pkg/log"
)

type controller struct {
	p user.User
}

func (u controller) Add(c *http.Context) {
	param := struct {
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	if err := c.BindAndValid(&param); err != nil {
		c.OutParamErr(err)
		return
	}
	if _, err := u.p.Add(param.Account, param.Password); err != nil {
		log.Errorf("user.Add err: %v", err)
		c.OutErr(errcode.UserAddError)
		return
	}
	c.OutSuccess(nil)
}
