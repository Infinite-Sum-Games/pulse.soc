package types

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/IAmRiteshKoushik/pulse/cmd"
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
	client := &http.Client{Timeout: 10 * time.Second}
	maxRetries := 3

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		// Add PAT to the Authorization header
		req.Header.Set("Authorization", "Bearer "+cmd.AppConfig.GitHub.PAT)
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		remaining := resp.Header.Get("X-RateLimit-Remaining")
		resetTime := resp.Header.Get("X-RateLimit-Reset")
		if remaining == "0" {
			resetUnix, _ := time.Parse(time.RFC3339, resetTime)
			waitDuration := time.Until(resetUnix)
			cmd.Log.Warn(
				fmt.Sprintf("Rate limit exceeded. Waiting for %v seconds before retrying...\n",
					waitDuration.Seconds(),
				),
			)
			continue
		}

		switch resp.StatusCode {
		case http.StatusOK:
			// Username found
			cmd.Log.Info("GitHub username verified successfully")
			return nil
		case http.StatusNotFound:
			// Username not found
			return fmt.Errorf("Invalid GitHub username")
		case http.StatusTooManyRequests:
			// Username search failed due to rate-limit
			// Exponential backoff with retries
			waitTime := time.Duration(attempt*attempt) * time.Second
			time.Sleep(waitTime)
			continue
		default:
			return fmt.Errorf("Unexpected status code. Could not find GitHub username.")
		}
	}
	return fmt.Errorf("Max retries exceeded for validating GitHub username.")
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

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginUserRequest) Validate() error {
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)

	return v.ValidateStruct(r,
		v.Field(&r.Email, v.Required, is.EmailFormat),
		v.Field(&r.Password, v.Required, v.Length(8, 130)),
	)
}