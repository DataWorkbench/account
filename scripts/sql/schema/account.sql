
CREATE TABLE `user` (
    `user_id` VARCHAR(50) NOT NULL,
    `name` VARCHAR(256) NOT NULL,
    `password` VARCHAR(255) DEFAULT '' NOT NULL,
    `email` VARCHAR(255) DEFAULT '' NOT NULL,
    `status` TINYINT NOT NULL,
    `role` TINYINT NOT NULL,
    `source` TINYINT(1) UNSIGNED DEFAULT 1 NOT NULL,
    `created` BIGINT(20) UNSIGNED NOT NULL,
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`user_id`),
    UNIQUE KEY (`name`)
) ENGINE=InnoDB;

CREATE TABLE `access_key` (
    `access_key_id` VARCHAR(64) NOT NULL,
    `secret_access_key` VARCHAR(256) DEFAULT '' NOT NULL,
    `owner` VARCHAR(50) NOT NULL,
    `name` VARCHAR(256) DEFAULT '' NOT NULL,
    `controller` TINYINT NOT NULL,
    `description` VARCHAR(512) DEFAULT '' NOT NULL,
    `status` TINYINT NOT NULL,
    `ip_white_list` VARCHAR(512) DEFAULT '' NOT NULL,
    `created` BIGINT(20) UNSIGNED NOT NULL,
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`access_key_id`)
) ENGINE=InnoDB;


CREATE TABLE `notification` (
    `owner` varchar(64) NOT NULL,
    `id` varchar(256) NOT NULL,
    `name` varchar(256) NOT NULL,
    `description` varchar(512) NOT NULL DEFAULT '',
    `email` varchar(64) NOT NULL,
    `created` bigint(20) unsigned NOT NULL,
    `updated` bigint(20) unsigned NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB;

