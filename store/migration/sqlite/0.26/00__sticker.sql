-- sticker_pack table
CREATE TABLE sticker_pack (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  display_name TEXT NOT NULL,
  description TEXT DEFAULT '',
  enabled BOOLEAN DEFAULT 1,
  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  updated_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now'))
);

-- sticker table
CREATE TABLE sticker (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  pack_id INTEGER NOT NULL,
  display_name TEXT NOT NULL,
  content BLOB,
  external_link TEXT DEFAULT '',
  type TEXT NOT NULL,
  size INTEGER DEFAULT 0,
  tags TEXT DEFAULT '',
  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  FOREIGN KEY(pack_id) REFERENCES sticker_pack(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX idx_sticker_pack_id ON sticker(pack_id);
CREATE INDEX idx_sticker_pack_enabled ON sticker_pack(enabled);
