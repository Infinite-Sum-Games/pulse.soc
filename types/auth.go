package types

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type RegisterUserRequest struct {
	Email      string `json:"email"`
	GhUsername string `json:"github_username"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
}

func (r *RegisterUserRequest) Validate() error {
	r.Email = strings.TrimSpace(r.Email)
	r.GhUsername = strings.TrimSpace(r.GhUsername)
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.MiddleName = strings.TrimSpace(r.MiddleName)
	r.LastName = strings.TrimSpace(r.LastName)

	err := v.ValidateStruct(r,
		v.Field(
			&r.Email,
			v.Required,
			is.EmailFormat,
			v.Match(regexp.MustCompile(`@cb.students.amrita.edu$`)),
		),
		v.Field(&r.GhUsername, v.Required, v.Length(3, 50)),
		v.Field(&r.FirstName, v.Required, v.Length(2, 50), is.Alpha),
		v.Field(&r.MiddleName, v.Length(0, 50), is.Alpha),
		v.Field(&r.LastName, v.Required, v.Length(1, 50), is.Alpha),
	)
	if err != nil {
		return err
	}

	// Check for Valid GitHub username
	url := fmt.Sprintf("https://api.github.com/users/%s", r.GhUsername)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("invalid github username")
	default:
		return fmt.Errorf("could not search for github username")
	}
}

type RegisterUserOtpVerifyRequest struct {
	Otp string `json:"otp"`
}

func (r *RegisterUserOtpVerifyRequest) Validate() error {
	r.Otp = strings.TrimSpace(r.Otp)

	return v.ValidateStruct(r,
		v.Field(&r.Otp, v.Required, v.Length(6, 6), is.Digit),
	)
}
