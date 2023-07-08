package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hotel-reservation/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func makeTestUser() *types.User
func TestAuthenticate(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HanldeAuthenticate)

	params := AuthParams{
		Email:    "james@foo.com",
		Password: "supersecurepassword",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status of 200 but got %d", resp.StatusCode)
	}
	var authResp AuthResponse
	if err = json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Error(err)
	}
	fmt.Println(resp)
}
