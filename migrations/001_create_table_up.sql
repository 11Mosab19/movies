CREATE TYPE "roles" AS ENUM ('admin','user');


CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    title TEXT UNIQUE NOT NULL,
    story_line TEXT NOT NULL,
    producer_name TEXT NOT NULL,
    category TEXT NOT NULL,
    release_date DATE,
    poster_url TEXT NOT NULL,
    status TEXT NOT NULL CHECK (length(btrim(status)) > 0),
    imdb_rating NUMERIC(3,1) CHECK (imdb_rating BETWEEN 0 AND 10),
    internal_rating NUMERIC(3,1) NOT NULL DEFAULT 0.0 CHECK (internal_rating BETWEEN 0 AND 10)
);

CREATE TABLE stars (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    home_land TEXT NOT NULL,
    date_of_birth DATE NOT NULL
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    full_name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    role roles DEFAULT 'user'
);

CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    review TEXT ,
    rating NUMERIC(3,1) NOT NULL CHECK(rating BETWEEN 0 AND 10),
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    movie_id INT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    UNIQUE (user_id, movie_id)
);

CREATE TABLE watch_list (
    movie_id INT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    user_id REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (movie_id, user_id)
);


CREATE TABLE movie_stars (
    movie_id INT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    star_id INT NOT NULL REFERENCES stars(id) ON DELETE CASCADE,
    PRIMARY KEY (movie_id, star_id)
);

CREATE INDEX movies_producer_name_idx ON movies(producer_name);
CREATE INDEX stars_name_idx ON stars(name);
CREATE INDEX users_email_idx ON users(email);
CREATE INDEX movies_title_idx ON movies(title);
CREATE INDEX movies_imdb_rating_idx ON movies(imdb_rating);