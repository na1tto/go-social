CREATE TABLE IF NOT EXISTS posts(
  id BIGSERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  user_id bigint NOT NULL,
  content text NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);