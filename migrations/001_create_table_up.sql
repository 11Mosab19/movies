CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    story_line TEXT NOT NULL,
    producer_name TEXT NOT NULL,
    category TEXT NOT NULL,
    release_date DATE,
    status TEXT NOT NULL CHECK (length(btrim(status)) > 0),
    imdb_rating NUMERIC(3,1) CHECK (imdb_rating BETWEEN 0 AND 10),
    internal_rating NUMERIC(3,1) NOT NULL DEFAULT 0.0
        CHECK (internal_rating BETWEEN 0 AND 10)
);

CREATE TABLE stars (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    home_land TEXT NOT NULL,
    date_of_birth DATE NOT NULL
);

CREATE TABLE movie_stars (
    movie_id INT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    star_id INT NOT NULL REFERENCES stars(id) ON DELETE CASCADE,
    PRIMARY KEY (movie_id, star_id)
);

CREATE INDEX movies_producer_name_idx ON movies (producer_name);
CREATE INDEX stars_name_idx ON stars (name);