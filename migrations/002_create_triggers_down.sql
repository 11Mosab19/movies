DROP TRIGGER IF EXISTS users_update_at_trigger ON users;
DROP TRIGGER IF EXISTS users_password_changed_trigger ON users;
DROP TRIGGER IF EXISTS update_movie_rating_trigger ON reviews;
DROP TRIGGER IF EXISTS trg_watch_list_updated_at ON watch_list;
DROP TRIGGER IF EXISTS reviews_update_at_trigger ON reviews;

DROP FUNCTION IF EXISTS update_at_users_update();
DROP FUNCTION IF EXISTS revoke_refresh_tokens_on_password_change();
DROP FUNCTION IF EXISTS update_movie_internal_rating();
DROP FUNCTION IF EXISTS update_watch_list_updated_at();
DROP FUNCTION IF EXISTS update_reviews_updated_at();
