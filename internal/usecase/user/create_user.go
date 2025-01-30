package user

import (
	"context"

	"github.com/andresilvase/gobid/internal/validator"
)

type CreateUserReq struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

func (req CreateUserReq) Valid(ctx context.Context) validator.Evaluator {

	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(req.UserName), "user_name", "this field is required")

	eval.CheckField(validator.NotBlank(req.Email), "email", "this field is required")

	eval.CheckField(
		validator.Matches(req.Email, validator.EmailRx),
		"email", "this field must be a valid email",
	)

	eval.CheckField(validator.NotBlank(req.Bio), "bio", "this field is required")

	eval.CheckField(
		validator.MinChars(req.Bio, 10),
		"bio", "this field must have a length betwwen 10 and 255",
	)

	eval.CheckField(
		validator.MaxChars(req.Bio, 255),
		"bio", "this field must have a length betwwen 10 and 255",
	)

	eval.CheckField(
		validator.MinChars(req.Password, 8),
		"password", "this field must have at least 8 chars",
	)

	return eval
}
