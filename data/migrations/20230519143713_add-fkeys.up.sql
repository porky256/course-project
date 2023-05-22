ALTER TABLE reservations
    ADD CONSTRAINT fk_reservations_room_id
        FOREIGN KEY (room_id)
            REFERENCES rooms(id)
            ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE room_restrictions
    ADD CONSTRAINT fk_restrictions_room_id
        FOREIGN KEY (room_id)
            REFERENCES rooms(id)
            ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE room_restrictions
    ADD CONSTRAINT fk_restrictions_restriction_id
        FOREIGN KEY (restriction_id)
            REFERENCES restrictions(id)
            ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE room_restrictions
    ADD CONSTRAINT fk_restrictions_reservation_id
        FOREIGN KEY (reservation_id)
            REFERENCES reservations(id)
            ON DELETE CASCADE ON UPDATE CASCADE;