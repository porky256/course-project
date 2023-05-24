INSERT INTO rooms (id,room_name) VALUES
                                  (1,'General''s Quarters'),
                                  (2,'Major''s Suite')
                                 ON CONFLICT (id) DO UPDATE SET room_name=EXCLUDED.room_name;