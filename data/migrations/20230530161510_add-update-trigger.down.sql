DROP TRIGGER IF EXISTS row_mod_on_room_restrictions_trigger_ ON room_restrictions;
DROP TRIGGER IF EXISTS row_mod_on_reservations_trigger_ ON reservations;
DROP TRIGGER IF EXISTS row_mod_on_rooms_trigger_ ON rooms;
DROP TRIGGER IF EXISTS row_mod_on_restrictions_drops_trigger_ ON restrictions;
DROP TRIGGER IF EXISTS row_mod_on_users_trigger_ ON users;

DROP FUNCTION IF EXISTS update_row_modified_function_();