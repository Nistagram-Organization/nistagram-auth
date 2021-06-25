package registration_request

import (
	"github.com/Nistagram-Organization/agent-shared/src/utils/rest_error"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/gender"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"github.com/Nistagram-Organization/nistagram-shared/src/proto"
	"net/mail"
	"net/url"
	"strings"
	"time"
)

const (
	USER  = "user"
	AGENT = "agent"
)

type RegistrationRequest struct {
	user.User
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (r *RegistrationRequest) Validate() rest_error.RestErr {
	if strings.TrimSpace(r.Username) == "" {
		return rest_error.NewBadRequestError("Username cannot be empty")
	}
	if len(strings.TrimSpace(r.Password)) < 8 {
		return rest_error.NewBadRequestError("Password must be at least 8 characters long")
	}
	if strings.TrimSpace(r.Name) == "" {
		return rest_error.NewBadRequestError("Name cannot be empty")
	}
	if strings.TrimSpace(r.Surname) == "" {
		return rest_error.NewBadRequestError("Surname cannot be empty")
	}
	if strings.TrimSpace(r.Phone) == "" {
		return rest_error.NewBadRequestError("Phone cannot be empty")
	}
	if r.Gender != gender.Male && r.Gender != gender.Female {
		return rest_error.NewBadRequestError("Gender can be 'male' or 'female'")
	}
	if time.Unix(r.BirthDate, 0).After(time.Now()) {
		return rest_error.NewBadRequestError("Birth date must be in the past")
	}
	if r.Role != USER && r.Role != AGENT {
		return rest_error.NewBadRequestError("Role can be 'user' or 'agent'")
	}
	if _, err := mail.ParseAddress(r.Email); err != nil {
		return rest_error.NewBadRequestError("Invalid email address")
	}
	if r.Role == AGENT {
		if _, err := url.ParseRequestURI(r.Website); err != nil {
			return rest_error.NewBadRequestError("Invalid website url")
		}
	}
	return nil
}

func (r *RegistrationRequest) ToUserMessage() *proto.UserMessage {
	registrationMessage := proto.UserMessage{
		Username:  r.Username,
		Password:  r.Password,
		Name:      r.Name,
		Surname:   r.Surname,
		BirthDate: r.BirthDate,
		Public:    r.Public,
		Taggable:  r.Taggable,
		Active:    true,
		Email:     r.Email,
		Website:   r.Website,
		Phone:     r.Phone,
	}

	if r.Gender == gender.Male {
		registrationMessage.Gender = proto.UserMessage_MALE
	} else {
		registrationMessage.Gender = proto.UserMessage_FEMALE
	}

	if r.Role == USER {
		registrationMessage.Role = proto.Role_USER
	} else {
		registrationMessage.Role = proto.Role_AGENT
	}

	return &registrationMessage
}
