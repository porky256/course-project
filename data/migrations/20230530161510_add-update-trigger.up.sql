CREATE OR REPLACE FUNCTION update_row_modified_function_()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$
    language PLPGSQL;

CREATE TRIGGER row_mod_on_users_trigger_ BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE PROCEDURE update_row_modified_function_();

CREATE TRIGGER row_mod_on_restrictions_drops_trigger_ BEFORE UPDATE ON restrictions
    FOR EACH ROW EXECUTE PROCEDURE update_row_modified_function_();

CREATE TRIGGER row_mod_on_rooms_trigger_ BEFORE UPDATE ON rooms
    FOR EACH ROW EXECUTE PROCEDURE update_row_modified_function_();

CREATE TRIGGER row_mod_on_reservations_trigger_ BEFORE UPDATE ON reservations
    FOR EACH ROW EXECUTE PROCEDURE update_row_modified_function_();

CREATE TRIGGER row_mod_on_room_restrictions_trigger_ BEFORE UPDATE ON room_restrictions
    FOR EACH ROW EXECUTE PROCEDURE update_row_modified_function_();