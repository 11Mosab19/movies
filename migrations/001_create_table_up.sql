CREATE TYPE "roles" AS ENUM ('admin','user');
CREATE TYPE "status" AS ENUM ('released','upcoming');
CREATE TYPE "watch_status" AS ENUM ('watched','listed');
CREATE TYPE "actor_role" AS ENUM ('star','guest');


CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    title TEXT UNIQUE NOT NULL,
    story_line TEXT NOT NULL,
    producer_name TEXT NOT NULL,
    category TEXT NOT NULL,
    release_date DATE,
    poster_url TEXT NOT NULL,
    status status NOT NULL,
    imdb_rating NUMERIC(3,1) CHECK (imdb_rating BETWEEN 0 AND 10),
    internal_rating NUMERIC(3,1) NOT NULL DEFAULT 0.0 CHECK (internal_rating BETWEEN 0 AND 10)
);

CREATE TABLE actors (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    home_land TEXT NOT NULL,
    date_of_birth DATE NOT NULL,
    UNIQUE (name, home_land,date_of_birth)
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    full_name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    role roles NOT NULL DEFAULT 'user'
);

CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    review TEXT NOT NULL,
    rating NUMERIC(3,1) NOT NULL CHECK(rating BETWEEN 0 AND 10),
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    movie_id INT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, movie_id)
);

CREATE TABLE watch_list (
    movie_id INT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    watch_status watch_status NOT NULL DEFAULT 'listed',
    set_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (movie_id, user_id)
);

CREATE TABLE movie_actors (
    movie_id INT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    actor_id INT NOT NULL REFERENCES actors(id) ON DELETE CASCADE,
    role actor_role NOT NULL,
    PRIMARY KEY (movie_id, actor_id)
);

CREATE TABLE refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX movies_producer_name_idx ON movies(producer_name);
CREATE INDEX actors_name_idx ON actors(name);
CREATE INDEX users_email_idx ON users(email);
CREATE INDEX movies_title_idx ON movies(title);
CREATE INDEX movies_imdb_rating_idx ON movies(imdb_rating);
CREATE INDEX watch_list_user_id_idx ON watch_list(user_id);
