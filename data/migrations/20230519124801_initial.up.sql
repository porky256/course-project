CREATE TABLE IF NOT EXISTS users (
    id           SERIAL NOT NULL PRIMARY KEY,
    first_name   VARCHAR(256) NOT NULL DEFAULT '',
    last_name    VARCHAR(256) NOT NULL DEFAULT '',
    email        VARCHAR(256) NOT NULL,
    password     VARCHAR(60) NOT NULL,
    access_level INTEGER NOT NULL DEFAULT 1,
    created_at   TIMESTAMP NOT NULL DEFAULT now(),
    updated_at   TIMESTAMP NOT NULL DEFAULT now()
    );

CREATE TABLE IF NOT EXISTS restrictions (
    id               SERIAL NOT NULL PRIMARY KEY,
    restriction_name VARCHAR(256) NOT NULL DEFAULT '',
    created_at       TIMESTAMP NOT NULL DEFAULT now(),
    updated_at       TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS rooms (
    id         SERIAL NOT NULL PRIMARY KEY,
    room_name  VARCHAR(256) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS reservations (
    id         SERIAL NOT NULL PRIMARY KEY,
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
    id                SERIAL NOT NULL PRIMARY KEY,
    start_date        DATE NOT NULL,
    end_date          DATE NOT NULL,
    room_id           INTEGER,
    reservation_id    INTEGER,
    restriction_id    INTEGER,
    created_at        TIMESTAMP NOT NULL DEFAULT now(),
    updated_at        TIMESTAMP NOT NULL DEFAULT now()
);
