DROP INDEX IF EXISTS actors_name_idx;
DROP INDEX IF EXISTS movies_producer_name_idx;
DROP INDEX IF EXISTS users_email_idx;
DROP INDEX IF EXISTS movies_title_idx;
DROP INDEX IF EXISTS movies_imdb_rating_idx;
DROP INDEX IF EXISTS watch_list_user_id_idx;


DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS watch_list;
DROP TABLE IF EXISTS movie_actors;
DROP TABLE IF EXISTS actors;
DROP TABLE IF EXISTS movies;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS actor_role;
DROP TYPE IF EXISTS watch_status;
DROP TYPE IF EXISTS status;
DROP TYPE IF EXISTS roles;
