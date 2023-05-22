ALTER TABLE reservations
    DROP CONSTRAINT fk_reservations_room_id;

ALTER TABLE room_restrictions
    DROP CONSTRAINT fk_restrictions_room_id;

ALTER TABLE room_restrictions
    DROP CONSTRAINT fk_restrictions_restriction_id;

ALTER TABLE room_restrictions
    DROP CONSTRAINT fk_restrictions_reservation_id;
