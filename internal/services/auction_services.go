package services

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Messagekind int

const (
	// Requests
	PlaceBid Messagekind = iota

	// Success
	SuccessfullyPlacedBid

	// Errors
	FailedToPlaceBid

	// Info
	NewBidPlaced
	AuctionFinished
)

type Message struct {
	UserId  uuid.UUID
	Message string
	Amount  float64
	Kind    Messagekind
}

type AuctionLobby struct {
	sync.Mutex
	Rooms map[uuid.UUID]*AuctionRoom
}

type AuctionRoom struct {
	Id         uuid.UUID
	Context    context.Context
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
	Clients    map[uuid.UUID]*Client

	BidsService BidsService
}

func (r *AuctionRoom) registerClient(client *Client) {
	slog.Info("New user connected", "Client", client)
	r.Clients[client.UserId] = client
}

func (r *AuctionRoom) unregisterClient(client *Client) {
	slog.Info("User disconnected", "Client", client)

	delete(r.Clients, client.UserId)
}

func (r *AuctionRoom) broadcastMessage(message Message) {
	slog.Info("New message received", "RoomID", r.Id, "Message", message, "Userid", message.UserId)

	switch message.Kind {
	case PlaceBid:
		bid, err := r.BidsService.PlaceBid(r.Context, r.Id, message.UserId, message.Amount)

		if err != nil {
			if errors.Is(err, ErrBidIsTooLow) {
				if client, ok := r.Clients[message.UserId]; ok {
					client.Send <- Message{
						Kind:    FailedToPlaceBid,
						Message: ErrBidIsTooLow.Error(),
					}
				}
			}
			return
		}

		if client, ok := r.Clients[message.UserId]; ok {
			client.Send <- Message{
				Kind:    SuccessfullyPlacedBid,
				Message: "Your Bid was placed successfully",
			}
		}

		for id, client := range r.Clients {
			newBidMessage := Message{
				Message: "A new bid was placed",
				Amount:  bid.BidAmount,
				Kind:    NewBidPlaced,
			}

			if id != message.UserId {
				client.Send <- newBidMessage
			}
		}
	}
}

func (r *AuctionRoom) Run() {
	slog.Info("Auction has started", "auctionID", r.Id)

	defer func() {
		close(r.Broadcast)
		close(r.Register)
		close(r.Unregister)
	}()

	for {
		select {
		case client := <-r.Register:
			r.registerClient(client)
		case client := <-r.Unregister:
			r.unregisterClient(client)
		case message := <-r.Broadcast:
			r.broadcastMessage(message)
		case <-r.Context.Done():
			slog.Info("Auction has ended", "auctionID", r.Id)

			for _, client := range r.Clients {
				client.Send <- Message{
					Kind:    AuctionFinished,
					Message: "Auction has been finished",
				}
			}

			return
		}
	}
}

func NewAuctionRoom(ctx context.Context, id uuid.UUID, bidService BidsService) *AuctionRoom {
	return &AuctionRoom{
		Id:          id,
		Broadcast:   make(chan Message),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Clients:     make(map[uuid.UUID]*Client),
		Context:     ctx,
		BidsService: bidService,
	}
}

type Client struct {
	Room   *AuctionRoom
	Conn   *websocket.Conn
	Send   chan Message
	UserId uuid.UUID
}

func NewClient(room *AuctionRoom, conn *websocket.Conn, userId uuid.UUID) *Client {
	return &Client{
		Room:   room,
		Conn:   conn,
		Send:   make(chan Message, 512),
		UserId: userId,
	}
}
