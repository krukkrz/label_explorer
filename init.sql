CREATE TABLE IF NOT EXISTS releases_count (
  id SERIAL PRIMARY KEY,
  label_id INTEGER NOT NULL,
  style_genre TEXT NOT NULL,
  artist_name TEXT NOT NULL,
  release_count INTEGER NOT NULL DEFAULT 0,
  UNIQUE(label_id, style_genre, artist_name)
);