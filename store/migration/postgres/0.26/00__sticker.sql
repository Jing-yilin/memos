-- sticker_pack table
CREATE TABLE sticker_pack (
  id SERIAL PRIMARY KEY,
  display_name VARCHAR(255) NOT NULL,
  description TEXT DEFAULT '',
  enabled BOOLEAN DEFAULT TRUE,
  created_ts BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
  updated_ts BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT
);

-- sticker table
CREATE TABLE sticker (
  id SERIAL PRIMARY KEY,
  pack_id INTEGER NOT NULL,
  display_name VARCHAR(255) NOT NULL,
  content BYTEA,
  external_link TEXT DEFAULT '',
  type VARCHAR(255) NOT NULL,
  size BIGINT DEFAULT 0,
  tags TEXT DEFAULT '',
  created_ts BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
  FOREIGN KEY (pack_id) REFERENCES sticker_pack(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX idx_sticker_pack_id ON sticker(pack_id);
CREATE INDEX idx_sticker_pack_enabled ON sticker_pack(enabled);
