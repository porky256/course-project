package dbrepo

import (
	"context"
	"errors"
	"github.com/porky256/course-project/internal/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	queryTimeout = 3 * time.Second
)

// InsertReservation inserts a reservation
func (pdb *postgresDB) InsertReservation(res *models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	var newID int
	err := pdb.DB.NewInsert().Model(res).Returning("id").Scan(ctx, &newID)
	return newID, err
}

// InsertRoom inserts a room
func (pdb *postgresDB) InsertRoom(room *models.Room) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	var newID int
	err := pdb.DB.NewInsert().Model(room).Returning("id").Scan(ctx, &newID)
	return newID, err
}

// InsertUser inserts an user
func (pdb *postgresDB) InsertUser(user *models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	var newID int
	err := pdb.DB.NewInsert().Model(user).Returning("id").Scan(ctx, &newID)
	return newID, err
}

// InsertRestriction inserts a restriction
func (pdb *postgresDB) InsertRestriction(res *models.Restriction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	var newID int
	err := pdb.DB.NewInsert().Model(res).Returning("id").Scan(ctx, &newID)
	return newID, err
}

// InsertRoomRestriction inserts a room restriction
func (pdb *postgresDB) InsertRoomRestriction(rmres *models.RoomRestriction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	var newID int
	err := pdb.DB.NewInsert().Model(rmres).Returning("id").Scan(ctx, &newID)
	return newID, err
}

// LookForAvailabilityOfRoom search if the room is available
func (pdb *postgresDB) LookForAvailabilityOfRoom(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	numberRows, err := pdb.DB.NewSelect().
		Table("room_restrictions").
		Where("room_id = ?", roomID).
		Where("end_date>?", start).
		Where("start_date<?", end).
		Count(ctx)
	return numberRows == 0, err
}

// AvailabilityOfAllRooms search for any room available on passed dates
func (pdb *postgresDB) AvailabilityOfAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	var rooms []models.Room
	subq := pdb.DB.NewSelect().
		Table("room_restrictions").
		Column("room_id").
		Where("end_date>?", start).
		Where("start_date<?", end)
	err := pdb.DB.NewSelect().
		Model((*models.Room)(nil)).
		Where("room.id not in (?)", subq).
		Scan(ctx, &rooms)
	return rooms, err
}

// GetRoomByID search for room by id
func (pdb *postgresDB) GetRoomByID(id int) (*models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	room := new(models.Room)
	err := pdb.DB.NewSelect().Model(room).Where("id=?", id).Scan(ctx)
	return room, err
}

// GetUserByID search for user by id
func (pdb *postgresDB) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	user := new(models.User)
	err := pdb.DB.NewSelect().Model(user).Where("id=?", id).Scan(ctx)
	return user, err
}

// Authenticate checks if user is known and password is correct
func (pdb *postgresDB) Authenticate(email, passwordSample string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	user := new(models.User)
	err := pdb.DB.NewSelect().Model(user).Where("email=?", email).Scan(ctx)

	if err != nil {
		return 0, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordSample))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return 0, "", errors.New("password is incorrect")
		} else {
			return 0, "", err
		}
	}

	return user.ID, user.Password, nil
}

// GetAllReservations search for all reservations
func (pdb *postgresDB) GetAllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	reservations := make([]models.Reservation, 0)
	err := pdb.DB.NewSelect().Model(&reservations).Relation("Room").Scan(ctx)

	return reservations, err
}

// GetNewReservations search for new reservations
func (pdb *postgresDB) GetNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	reservations := make([]models.Reservation, 0)
	err := pdb.DB.NewSelect().Model(&reservations).Relation("Room").Where("is_processed=0").Scan(ctx)

	return reservations, err
}

// GetReservationByID search for reservation by id
func (pdb *postgresDB) GetReservationByID(id int) (*models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	reservation := new(models.Reservation)
	err := pdb.DB.NewSelect().Model(reservation).Relation("Room").Where("reservation.id=?", id).Scan(ctx)

	return reservation, err
}

// UpdateReservation updates reservation
func (pdb *postgresDB) UpdateReservation(ur models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	_, err := pdb.DB.NewUpdate().Model(&ur).
		Column("first_name", "last_name", "email", "phone").
		WherePK().Exec(ctx)
	return err
}

// DeleteReservation deletes reservation
func (pdb *postgresDB) DeleteReservationByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	_, err := pdb.DB.NewDelete().Table("reservations").Where("id=?", id).Exec(ctx)
	return err
}

// UpdateReservationProcessed updates is_processed field in reservation
func (pdb *postgresDB) UpdateReservationProcessed(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	modelToUpdate := models.Reservation{
		ID:          id,
		IsProcessed: processed,
	}
	_, err := pdb.DB.NewUpdate().Model(&modelToUpdate).WherePK().Column("is_processed").Exec(ctx)
	return err
}
