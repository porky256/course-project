package repository

import (
	"github.com/porky256/course-project/internal/models"
	"time"
)

type DatabaseRepo interface {
	InsertReservation(res *models.Reservation) (int, error)
	GetReservationByID(id int) (*models.Reservation, error)
	GetAllReservations() ([]models.Reservation, error)
	GetNewReservations() ([]models.Reservation, error)
	UpdateReservation(ur models.Reservation) error
	UpdateReservationProcessed(id, processed int) error
	DeleteReservationByID(id int) error

	InsertRoom(room *models.Room) (int, error)
	GetRoomByID(id int) (*models.Room, error)
	GetAllRooms() ([]models.Room, error)
	LookForAvailabilityOfRoom(start, end time.Time, roomID int) (bool, error)
	AvailabilityOfAllRooms(start, end time.Time) ([]models.Room, error)

	InsertUser(user *models.User) (int, error)
	GetUserByID(id int) (*models.User, error)

	InsertRestriction(res *models.Restriction) (int, error)

	InsertRoomRestriction(rmres *models.RoomRestriction) (int, error)
	AddSingleDayRoomRestriction(roomID, restrictionID int, start time.Time) (int, error)
	GetRoomRestrictionsByRoomIdWithinDates(roomID int, start, end time.Time) ([]models.RoomRestriction, error)
	DeleteRoomRestrictionByID(id int) error

	Authenticate(email, passwordSample string) (int, string, error)
}
