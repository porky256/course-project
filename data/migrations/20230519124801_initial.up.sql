CREATE TABLE IF NOT EXISTS users (
    id           INTEGER PRIMARY KEY,
    first_name   VARCHAR(256) NOT NULL DEFAULT '',
    last_name    VARCHAR(256) NOT NULL DEFAULT '',
    email        VARCHAR(256) NOT NULL,
    password     VARCHAR(60) NOT NULL,
    access_level INTEGER NOT NULL DEFAULT 1,
    created_at   TIMESTAMP NOT NULL DEFAULT now(),
    updated_at   TIMESTAMP NOT NULL DEFAULT now()
    );

CREATE TABLE IF NOT EXISTS restrictions (
    id               INTEGER PRIMARY KEY,
    restriction_name VARCHAR(256) NOT NULL DEFAULT '',
    created_at       TIMESTAMP NOT NULL DEFAULT now(),
    updated_at       TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS rooms (
    id         INTEGER PRIMARY KEY,
    room_name  VARCHAR(256) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS reservations (
    id         INTEGER PRIMARY KEY,
    first_name VARCHAR(256) NOT NULL DEFAULT '',
    last_name  VARCHAR(256) NOT NULL DEFAULT '',
    email      VARCHAR(256) NOT NULL UNIQUE,
    phone      VARCHAR(256) NOT NULL DEFAULT '',
    start_date DATE NOT NULL,
    end_date   DATE NOT NULL,
    room_id    INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS room_restrictions (
    id                INTEGER PRIMARY KEY,
    start_date        DATE NOT NULL,
    end_date          DATE NOT NULL,
    room_id           INTEGER,
    reservation_id    INTEGER,
    restriction_id    INTEGER,
    created_at        TIMESTAMP NOT NULL DEFAULT now(),
    updated_at        TIMESTAMP NOT NULL DEFAULT now()
);
