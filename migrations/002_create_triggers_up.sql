CREATE OR REPLACE FUNCTION update_at_users_update()
RETURNS TRIGGER AS $$
BEGIN 
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_update_at_trigger
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_at_users_update();
--
CREATE OR REPLACE FUNCTION update_movie_internal_rating()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE movies
    SET internal_rating = COALESCE(
        (
            SELECT AVG(rating)
            FROM reviews
            WHERE movie_id = COALESCE(NEW.movie_id, OLD.movie_id)
        ),
        0.0
    )
    WHERE id = COALESCE(NEW.movie_id, OLD.movie_id);

    RETURN NULL;
END;
$$;

CREATE TRIGGER update_movie_rating_trigger
AFTER INSERT OR UPDATE OR DELETE
ON reviews
FOR EACH ROW
EXECUTE FUNCTION update_movie_internal_rating();