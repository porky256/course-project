CREATE UNIQUE INDEX users_email_idx ON users (email);

CREATE INDEX room_restrictions_start_date_end_date_idx ON room_restrictions (start_date,end_date);