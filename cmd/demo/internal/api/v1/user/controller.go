package user

import (
	"github.com/guidoxie/knife/cmd/demo/internal/dao/user"
	"github.com/guidoxie/knife/pkg/app/http"
)

type controller struct {
	p user.User
}

func (u controller) Add(c *http.Context) {
	param := struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}{}

	if err := c.BindJSON(&param); err != nil {
		return
	}
	if _, err := u.p.Add(param.Account, param.Password); err != nil {

		return
	}
	c.OutSuccess(nil)
}
