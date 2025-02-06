package api

import (
	"net/http"

	"github.com/andresilvase/gobid/internal/jsonutils"
	"github.com/andresilvase/gobid/internal/store/pgstore"
	"github.com/andresilvase/gobid/internal/usecase/product"
	"github.com/google/uuid"
)

func (api *Api) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[product.CreateProductReq](r)

	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	userId, ok := api.Sessions.Get(r.Context(), "AuthenticatedUserId").(uuid.UUID)

	if !ok {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected error, try again later",
		})
		return
	}

	id, err := api.ProductService.CreateProduct(r.Context(),
		pgstore.CreateProductParams{
			SellerID:    userId,
			ProductName: data.ProductName,
			Description: data.Description,
			Baseprice:   data.Baseprice,
			AuctionEnd:  data.AuctionEnd,
		},
	)

	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "failed to create product auction, try again later",
		})
		return
	}

	jsonutils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"message":    "product created successfully",
		"product_id": id,
	})

}
