package api

import (
	aptly_password "github.com/bborbe/aptly_utils/password"
	aptly_url "github.com/bborbe/aptly_utils/url"
	aptly_user "github.com/bborbe/aptly_utils/user"
)

type Api struct {
	Url      aptly_url.Url
	User     aptly_user.User
	Password aptly_password.Password
}

func New(url string, user string, password string) Api {
	return Api{
		Url: aptly_url.Url(url), User: aptly_user.User(user), Password: aptly_password.Password(password)}
}
