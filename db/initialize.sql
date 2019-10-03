CREATE USER `emojicord`@`%` IDENTIFIED BY 'password';
GRANT INSERT,SELECT,UPDATE,DELETE ON `emojicord_db`.* TO `emojicord`@`%`;

CREATE DATABASE IF NOT EXISTS `emojicord_db`;

CREATE TABLE IF NOT EXISTS `emojicord_db`.`users` (
    `id`               VARCHAR(64)    NOT NULL,
    `username`         VARCHAR(32)    NOT NULL,
    `discriminator`    VARCHAR(4)     NOT NULL,
    `avatar`           VARCHAR(34)    NOT NULL,
    `locale`           VARCHAR(16),
    `created_at`       TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `last_login`       TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY ( `id` )
);

CREATE TABLE IF NOT EXISTS `emojicord_db`.`users__discord_tokens` (
    `user_id`          VARCHAR(64)     NOT NULL,
    `access_token`     VARCHAR(255)    NOT NULL,
    `token_type`       VARCHAR(6)      NOT NULL,
    `refresh_token`    VARCHAR(255)    NOT NULL,
    `expiry`           TIMESTAMP       NOT NULL,

    PRIMARY KEY ( `user_id` ),

    CONSTRAINT `fk_users__discord_tokens__users`
        FOREIGN KEY ( `user_id` )
        REFERENCES `emojicord_db`.`users` ( `id` )
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `emojicord_db`.`discord_guilds` (
    `id`                  VARCHAR(64)     NOT NULL,
    `name`                VARCHAR(100)    NOT NULL,
    `icon`                VARCHAR(255)    NOT NULL,
    `is_bot_exists`       BOOLEAN         NOT NULL,

    PRIMARY KEY ( `id` )
);

CREATE TABLE IF NOT EXISTS `emojicord_db`.`users__discord_guilds` (
    `user_id`        VARCHAR(64)     NOT NULL,
    `guild_id`       VARCHAR(64)     NOT NULL,
    `is_owner`       BOOLEAN         NOT NULL,
    `permissions`    INT UNSIGNED    NOT NULL,
    `can_invite`     BOOLEAN         NOT NULL,
    
    PRIMARY KEY ( `user_id`, `guild_id` ),

    CONSTRAINT `fk__users__discord_guilds__users`
        FOREIGN KEY ( `user_id` )
        REFERENCES `emojicord_db`.`users` ( `id` )
        ON DELETE CASCADE,
    CONSTRAINT `fk__users__discord_guilds__discord_guilds`
        FOREIGN KEY ( `guild_id` )
        REFERENCES `emojicord_db`.`discord_guilds` ( `id` )
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `emojicord_db`.`discord_emojis` (
    `id`                   VARCHAR(64)   NOT NULL,
    `guild_id`             VARCHAR(64)   NOT NULL,
    `name`                 VARCHAR(32)   NOT NULL,
    `is_animated`          BOOLEAN       NOT NULL,
    `user_id`              VARCHAR(64)   NOT NULL,
    `user_username`        VARCHAR(32)   NOT NULL,
    `user_discriminator`   VARCHAR(4)    NOT NULL,
    `user_avatar`          VARCHAR(34)   NOT NULL,

    PRIMARY KEY ( `id` ),

    CONSTRAINT `fk__discord_guilds__discord_emojis`
        FOREIGN KEY ( `guild_id` )
        REFERENCES `emojicord_db`.`discord_guilds` ( `id` )
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `emojicord_db`.`users__discord_emojis` (
    `user_id`      VARCHAR(64)   NOT NULL,
    `emoji_id`     VARCHAR(64)   NOT NULL,
    `is_enabled`   BOOLEAN       NOT NULL,

    PRIMARY KEY ( `user_id`, `emoji_id` ),

    CONSTRAINT `fk__users__discord_emojis__users`
        FOREIGN KEY ( `user_id` )
        REFERENCES `emojicord_db`.`users` ( `id` )
        ON DELETE CASCADE,
    CONSTRAINT `fk__users__discord_emojis__discord_emojis`
        FOREIGN KEY ( `emoji_id` )
        REFERENCES `emojicord_db`.`discord_emojis` ( `id` )
        ON DELETE CASCADE
);