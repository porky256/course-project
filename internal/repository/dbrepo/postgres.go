package dbrepo

import (
	"context"
	"github.com/porky256/course-project/internal/models"
	"time"
)

const (
	queryTimeout = 3 * time.Second
)

func (pdb *postgressDB) InsertReservation(res *models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	var newID int
	err := pdb.DB.NewInsert().Model(res).Returning("id").Scan(ctx, &newID)
	return newID, err
}

func (pdb *postgressDB) InsertRoom(room *models.Room) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	var newID int
	err := pdb.DB.NewInsert().Model(room).Returning("id").Scan(ctx, &newID)
	return newID, err
}

func (pdb *postgressDB) InsertUser(user *models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	var newID int
	err := pdb.DB.NewInsert().Model(user).Returning("id").Scan(ctx, &newID)
	return newID, err
}

func (pdb *postgressDB) InsertRestriction(res *models.Restriction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	var newID int
	err := pdb.DB.NewInsert().Model(res).Returning("id").Scan(ctx, &newID)
	return newID, err
}

func (pdb *postgressDB) InsertRoomRestriction(rmres *models.RoomRestriction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	var newID int
	err := pdb.DB.NewInsert().Model(rmres).Returning("id").Scan(ctx, &newID)
	return newID, err
}

func (pdb *postgressDB) GetRoom(id int) (*models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	res := new(models.Room)
	err := pdb.DB.NewSelect().Model(res).Where("id = ?", id).Scan(ctx)
	return res, err
}

func (pdb *postgressDB) LookForAvailabilityOfRoom(start, end time.Time, roomID int) (bool, error) {
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

func (pdb *postgressDB) AvailabilityOfAllRooms(start, end time.Time) ([]models.Room, error) {
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

func (pdb *postgressDB) GetRoomByID(id int) (*models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	room := new(models.Room)
	err := pdb.DB.NewSelect().Model(&room).Where("id==?", id).Scan(ctx)
	return room, err
}