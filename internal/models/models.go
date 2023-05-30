package models

import (
	"time"
)

type User struct {
	ID          int `bun:",pk,autoincrement"`
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time `bun:",nullzero"`
	UpdatedAt   time.Time `bun:",nullzero"`
}

type Room struct {
	ID        int       `bun:",pk,autoincrement"`
	Name      string    `bun:"room_name"`
	CreatedAt time.Time `bun:",nullzero"`
	UpdatedAt time.Time `bun:",nullzero"`
}

type Restriction struct {
	ID              int `bun:",pk,autoincrement"`
	RestrictionName string
	CreatedAt       time.Time `bun:",nullzero"`
	UpdatedAt       time.Time `bun:",nullzero"`
}

type Reservation struct {
	ID          int `bun:",pk,autoincrement"`
	FirstName   string
	LastName    string
	Email       string
	Phone       string
	StartDate   time.Time `bun:"type:Date"`
	EndDate     time.Time `bun:"type:Date"`
	RoomID      int
	IsProcessed int
	CreatedAt   time.Time `bun:",nullzero"`
	UpdatedAt   time.Time `bun:",nullzero"`
	Room        *Room     `bun:"rel:belongs-to,join:room_id=id"`
}

type RoomRestriction struct {
	ID            int       `bun:",pk,autoincrement"`
	StartDate     time.Time `bun:"type:Date"`
	EndDate       time.Time `bun:"type:Date"`
	RoomID        int
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time    `bun:",nullzero"`
	UpdatedAt     time.Time    `bun:",nullzero"`
	Room          *Room        `bun:"rel:belongs-to,join:room_id=id"`
	Reservation   *Reservation `bun:"rel:belongs-to,join:reservation_id=id"`
	Restriction   *Restriction `bun:"rel:belongs-to,join:restriction_id=id"`
}
