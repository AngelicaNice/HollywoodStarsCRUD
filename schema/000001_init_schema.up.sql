CREATE TABLE actors (
  id          serial      NOT NULL UNIQUE,
  name        varchar(20) NOT NULL,
  surname     varchar(20) NOT NULL,
  sex         varchar(6)  NOT NULL,
  birth_year  integer     NOT NULL,
  birth_place varchar(15) NOT NULL,
  rest_year   integer,
  language    varchar(15)
);

CREATE TABLE follows (
  follow_id         serial      NOT NULL UNIQUE,
  following_user_id integer     NOT NULL,
  followed_actor_id integer     NOT NULL,
  created_at        timestamp   NOT NULL DEFAULT (now())
);

CREATE TABLE users (
  id            serial      NOT NULL UNIQUE,
  nickname      varchar(20) NOT NULL UNIQUE,
  email         varchar(30) NOT NULL UNIQUE,
  password      varchar(90) NOT NULL,
  registered_at timestamp   NOT NULL DEFAULT (now())
);

CREATE INDEX ON actors (surname);

ALTER TABLE follows ADD FOREIGN KEY (following_user_id) REFERENCES users (id);

ALTER TABLE follows ADD FOREIGN KEY (followed_actor_id) REFERENCES actors (id);
