ALTER TABLE actors
DROP COLUMN IF EXISTS image_public_id,
DROP COLUMN IF EXISTS image_url;

ALTER TABLE movies
DROP COLUMN IF EXISTS poster_public_id;
