INSERT INTO restrictions (id,restriction_name) VALUES
                                     (1,'Reservation'),
                                     (2,'Owner Block')
                                        ON CONFLICT (id) DO UPDATE SET restriction_name=EXCLUDED.restriction_name;