package api

import (
	"context"
	"net/http"

	"github.com/andresilvase/gobid/internal/jsonutils"
	"github.com/andresilvase/gobid/internal/services"
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

	product_id, err := api.ProductService.CreateProduct(r.Context(),
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

	ctx, _ := context.WithDeadline(context.Background(), data.AuctionEnd)

	auctionRoom := services.NewAuctionRoom(ctx, product_id, api.BidsService)

	go auctionRoom.Run()

	api.AuctionLobby.Lock()
	api.AuctionLobby.Rooms[product_id] = auctionRoom
	api.AuctionLobby.Unlock()

	jsonutils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"message":    "Auction has started successfully",
		"product_id": product_id,
	})

}
