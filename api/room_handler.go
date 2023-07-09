package api

import (
	"hotel-reservation/db"
	"time"

	"github.com/gofiber/fiber/v2"
)

type BookRoomParams struct {
	FromDate time.Time `json:"fromDate"`
	TillDate time.Time `json:"tillDate"`
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store,
	}
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	return nil
}
