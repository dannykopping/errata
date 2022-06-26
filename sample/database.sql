PRAGMA foreign_keys= OFF;

BEGIN TRANSACTION;

CREATE TABLE users
(
    email         text not null,
    password_hash text not null,
    spam          int(1) default 0,
    abuse         int(1) default 0
);
INSERT INTO users
VALUES ('valid@email.com', '5f4dcc3b5aa765d61d8327deb882cf99', 0, 0);
INSERT INTO users
VALUES ('abuse@email.com', '5f4dcc3b5aa765d61d8327deb882cf99', 0, 1);
INSERT INTO users
VALUES ('spam@email.com', '5f4dcc3b5aa765d61d8327deb882cf99', 1, 0);

COMMIT;