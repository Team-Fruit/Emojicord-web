CREATE USER `emojicord`@`%`;
GRANT INSERT,SELECT,UPDATE,DELETE ON `emojicord_db`.* TO `emojicord`@`%`;

CREATE DATABASE IF NOT EXISTS `emojicord_db`;
