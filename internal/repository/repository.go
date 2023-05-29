package repository

import (
	"github.com/porky256/course-project/internal/models"
	"time"
)

type DatabaseRepo interface {
	InsertReservation(res *models.Reservation) (int, error)
	InsertRoom(room *models.Room) (int, error)
	InsertUser(user *models.User) (int, error)
	InsertRestriction(res *models.Restriction) (int, error)
	InsertRoomRestriction(rmres *models.RoomRestriction) (int, error)
	LookForAvailabilityOfRoom(start, end time.Time, roomID int) (bool, error)
	AvailabilityOfAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (*models.Room, error)
	GetUserByID(id int) (*models.User, error)
	Authenticate(email, passwordSample string) (int, string, error)
}
