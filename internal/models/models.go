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
	ID        int `bun:",pk,autoincrement"`
	RoomName  string
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
	ID        int `bun:",pk,autoincrement"`
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time `bun:"type:Date"`
	EndDate   time.Time `bun:"type:Date"`
	RoomID    int
	CreatedAt time.Time `bun:",nullzero"`
	UpdatedAt time.Time `bun:",nullzero"`
	Room      *Room     `bun:"-"`
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
	Room          *Room        `bun:"-"`
	Reservation   *Reservation `bun:"-"`
	Restriction   *Restriction `bun:"-"`
}
