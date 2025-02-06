package product

import (
	"context"
	"time"

	"github.com/andresilvase/gobid/internal/validator"
)

type CreateProductReq struct {
	ProductName string    `json:"product_name"`
	Description string    `json:"description"`
	Baseprice   float64   `json:"baseprice"`
	AuctionEnd  time.Time `json:"auction_end"`
}

const minAuctionDuration = 2 * time.Hour

func (req CreateProductReq) Valid(ctx context.Context) validator.Evaluator {
	var val validator.Evaluator

	val.CheckField(validator.NotBlank(req.ProductName), "product_name", "this field is required")
	val.CheckField(
		validator.MinChars(req.Description, 10) &&
			validator.MaxChars(req.Description, 255),
		"description", "description must have a length between 10 and 255",
	)
	val.CheckField(req.Baseprice > 0, "baseprice", "this field must be greater than zero")
	val.CheckField(time.Until(req.AuctionEnd) >= minAuctionDuration, "auction_end", "must to be at least two hours")

	return val
}
