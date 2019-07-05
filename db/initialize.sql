CREATE USER `emojicord`@`%`;
GRANT INSERT,SELECT,UPDATE,DELETE ON `emojicord_db`.* TO `emojicord`@`%`;

CREATE DATABASE IF NOT EXISTS `emojicord_db`;

CREATE TABLE IF NOT EXISTS `emojicord_db`.`users` (
    `id`               VARCHAR(64)    NOT NULL,
    `username`         VARCHAR(32)    NOT NULL,
    `discriminator`    VARCHAR(4)     NOT NULL,
    `avatar`           VARCHAR(32)    NOT NULL,
    `locale`           VARCHAR(16),
    `created_at`       TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`       TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY ( `id` )
);

CREATE TABLE IF NOT EXISTS `emojicord_db`.`users__discord_tokens` (
    `user_id`          VARCHAR(64)     NOT NULL,
    `access_token`     VARCHAR(255)    NOT NULL,
    `token_type`       VARCHAR(255)    NOT NULL,
    `refresh_token`    VARCHAR(255)    NOT NULL,
    `expiry`           INT UNSIGNED    NOT NULL,

    PRIMARY KEY ( `user_id` ),

    CONSTRAINT `fk_users__discord_tokens__users`
        FOREIGN KEY ( `user_id` )
        REFERENCES `emojicord_db`.`users` ( `id` )
        ON DELETE CASCADE
);