CREATE OR REPLACE FUNCTION update_at_users_update()
RETURNS TRIGGER AS $$
BEGIN 
    NEW.updated_at := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_update_at_trigger
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_at_users_update();

CREATE OR REPLACE FUNCTION revoke_refresh_tokens_on_password_change()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.hashed_password IS DISTINCT FROM OLD.hashed_password THEN
        DELETE FROM refresh_tokens WHERE user_id = NEW.id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_password_changed_trigger
BEFORE UPDATE OF hashed_password ON users
FOR EACH ROW
EXECUTE FUNCTION revoke_refresh_tokens_on_password_change();
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

---

CREATE OR REPLACE FUNCTION update_watch_list_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.watch_status IS DISTINCT FROM OLD.watch_status THEN
        NEW.updated_at := CURRENT_TIMESTAMP;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_watch_list_updated_at
BEFORE UPDATE ON watch_list
FOR EACH ROW
EXECUTE FUNCTION update_watch_list_updated_at();

---

CREATE OR REPLACE FUNCTION update_reviews_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER reviews_update_at_trigger
BEFORE UPDATE ON reviews
FOR EACH ROW
EXECUTE FUNCTION update_reviews_updated_at();
