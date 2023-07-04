package db

const (
	DBNAME     = "hotel-reservation"
	TestDBNAME = "hotel-reservation-test"
	DBURI      = "mongodb://localhost:27017"
)

type Stotre struct {
	user  UserStore
	hotel HotelStore
	room  RoomStore
}
