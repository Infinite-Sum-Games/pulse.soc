package types

import (
	"strings"

	"github.com/IAmRiteshKoushik/pulse/pkg"
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

/*
* Request and response for login request. During the response, the cookies are
* setup as well for further authentication requests
 */
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	// Trim space
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)

	return v.ValidateStruct(r,
		v.Field(r.Email, v.Required, is.Email),
		v.Field(r.Password,
			v.Required,
			v.Length(8, 50),
			v.Match(pkg.Rex.PasswordRegex),
		),
	)
}

// TODO: Need to send more data
type LoginResponse struct {
	Message   string `json:"message"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

/*
* Request and response for customer registrations. During the response, a
* 5-minute valid token is setup which is added in the cookie jar. And this can
* be used to further add in client-side validations if the OTP is being
* entered really late!
 */
type CustomerRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *CustomerRegisterRequest) Validate() error {
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)

	return v.ValidateStruct(r,
		v.Field(r.Email, v.Required, is.Email),
		v.Field(r.Password,
			v.Required,
			v.Length(9, 50),
			v.Match(pkg.Rex.PasswordRegex),
		),
	)
}
