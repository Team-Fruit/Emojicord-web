CREATE USER `emojicord`@`%` IDENTIFIED BY 'password';
GRANT INSERT,SELECT,UPDATE,DELETE ON `emojicord_db`.* TO `emojicord`@`%`;

CREATE DATABASE IF NOT EXISTS `emojicord_db`;

CREATE TABLE IF NOT EXISTS `emojicord_db`.`users` (
    `id`               BIGINT UNSIGNED         NOT NULL,
    `username`         VARCHAR(32)             NOT NULL,
    `discriminator`    SMALLINT(4) ZEROFILL    NOT NULL,
    `avatar`           VARCHAR(34)             NOT NULL,
    `locale`           VARCHAR(16),
    `created_at`       TIMESTAMP               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `last_login`       TIMESTAMP               NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY ( `id` )
);

CREATE TABLE IF NOT EXISTS `emojicord_db`.`users__discord_tokens` (
    `user_id`          BIGINT UNSIGNED     NOT NULL,
    `access_token`     VARCHAR(255)        NOT NULL,
    `token_type`       VARCHAR(6)          NOT NULL,
    `refresh_token`    VARCHAR(255)        NOT NULL,
    `expiry`           TIMESTAMP           NOT NULL,

    PRIMARY KEY ( `user_id` ),

    CONSTRAINT `fk_users__discord_tokens__users`
        FOREIGN KEY ( `user_id` )
        REFERENCES `emojicord_db`.`users` ( `id` )
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `emojicord_db`.`discord_guilds` (
    `id`                  BIGINT UNSIGNED    NOT NULL,
    `name`                VARCHAR(100)       NOT NULL,
    `icon`                VARCHAR(255)       NOT NULL,
    `is_bot_exists`       BOOLEAN            NOT NULL,

    PRIMARY KEY ( `id` )
);

CREATE TABLE IF NOT EXISTS `emojicord_db`.`users__discord_guilds` (
    `user_id`        BIGINT UNSIGNED    NOT NULL,
    `guild_id`       BIGINT UNSIGNED    NOT NULL,
    `is_owner`       BOOLEAN            NOT NULL,
    `permissions`    INT UNSIGNED       NOT NULL,
    `can_invite`     BOOLEAN            NOT NULL,
    
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

CREATE TABLE IF NOT EXISTS `emojicord_db`.`discord_emojis_users` (
    `id`               BIGINT UNSIGNED         NOT NULL,
    `username`         VARCHAR(32)             NOT NULL,
    `discriminator`    SMALLINT(4) ZEROFILL    NOT NULL,
    `avatar`           VARCHAR(34)             NOT NULL,

    PRIMARY KEY ( `id` )
);

CREATE TABLE IF NOT EXISTS `emojicord_db`.`discord_emojis` (
    `id`                   BIGINT UNSIGNED    NOT NULL,
    `guild_id`             BIGINT UNSIGNED    NOT NULL,
    `name`                 VARCHAR(32)        NOT NULL,
    `is_animated`          BOOLEAN            NOT NULL,
    `user_id`              BIGINT UNSIGNED    NOT NULL,
    `is_deleted`           BOOLEAN            NOT NULL DEFAULT 0,

    PRIMARY KEY ( `id` ),

    CONSTRAINT `fk__discord_guilds__discord_emojis`
        FOREIGN KEY ( `guild_id` )
        REFERENCES `emojicord_db`.`discord_guilds` ( `id` )
        ON DELETE CASCADE,
    CONSTRAINT `fk__discord_emojis_users__discord_emojis`
        FOREIGN KEY ( `user_id` )
        REFERENCES `emojicord_db`.`discord_emojis_users` ( `id` )
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `emojicord_db`.`users__discord_emojis` (
    `user_id`       BIGINT UNSIGNED    NOT NULL,
    `emoji_id`      BIGINT UNSIGNED    NOT NULL,
    `is_enabled`    BOOLEAN            NOT NULL,

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