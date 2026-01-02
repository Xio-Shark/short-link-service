-- 最小可用表结构
CREATE TABLE IF NOT EXISTS links (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  code VARCHAR(16) NOT NULL,
  original_url VARCHAR(2048) NOT NULL,
  expire_at BIGINT NOT NULL DEFAULT 0,
  status TINYINT NOT NULL DEFAULT 1,
  created_at BIGINT NOT NULL,
  updated_at BIGINT NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_code (code),
  KEY idx_expire_at (expire_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS visits (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  link_id BIGINT UNSIGNED NOT NULL,
  ip VARCHAR(64) NOT NULL,
  user_agent VARCHAR(255) NOT NULL,
  created_at BIGINT NOT NULL,
  PRIMARY KEY (id),
  KEY idx_link_id (link_id),
  KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
