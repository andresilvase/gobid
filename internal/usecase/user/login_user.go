package user

import (
	"context"

	"github.com/andresilvase/gobid/internal/validator"
)

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req LoginUserReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.Matches(req.Email, validator.EmailRx), "email", "email is invalid")
	eval.CheckField(validator.NotBlank(req.Password), "password", "password cannot be blank")

	return eval
}
