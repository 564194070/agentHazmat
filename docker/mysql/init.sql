CREATE DATABASE IF NOT EXISTS agent_user DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE agent_user;

CREATE TABLE IF NOT EXISTS users (
  id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  deleted_at DATETIME(3) NULL,
  username   VARCHAR(64)  NOT NULL,
  password   VARCHAR(255) NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_users_username (username),
  KEY idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;