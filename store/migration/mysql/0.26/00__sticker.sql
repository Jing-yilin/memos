-- sticker_pack table
CREATE TABLE sticker_pack (
  id INT PRIMARY KEY AUTO_INCREMENT,
  display_name VARCHAR(255) NOT NULL,
  description TEXT,
  enabled BOOLEAN DEFAULT TRUE,
  created_ts BIGINT NOT NULL DEFAULT (UNIX_TIMESTAMP()),
  updated_ts BIGINT NOT NULL DEFAULT (UNIX_TIMESTAMP())
);

-- sticker table  
CREATE TABLE sticker (
  id INT PRIMARY KEY AUTO_INCREMENT,
  pack_id INT NOT NULL,
  display_name VARCHAR(255) NOT NULL,
  content LONGBLOB,
  external_link TEXT,
  type VARCHAR(255) NOT NULL,
  size BIGINT DEFAULT 0,
  tags TEXT,
  created_ts BIGINT NOT NULL DEFAULT (UNIX_TIMESTAMP()),
  FOREIGN KEY (pack_id) REFERENCES sticker_pack(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX idx_sticker_pack_id ON sticker(pack_id);
CREATE INDEX idx_sticker_pack_enabled ON sticker_pack(enabled);
