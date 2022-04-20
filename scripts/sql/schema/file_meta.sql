CREATE TABLE IF NOT EXISTS `file`
(
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- File ID, unique within a region
    `id` CHAR(20) NOT NULL,

    -- Id of Parent Directory. pid is "" means root(`/`).
    `pid` CHAR(20) NOT NULL,

    -- IsDirectory represents this job whether a directory.
    `is_directory` BOOL,

    -- File Name, Unique within a workspace
    -- The max length of use set is 128. The system will be auto rename to <name>.<id> when deleted.
    -- Thus the VARCHAR should be define as 149 (128 + 20 + 1)
    `name` VARCHAR(149) NOT NULL,

    -- File description.
    `desc` VARCHAR(1024) CHARACTER SET utf8mb4 DEFAULT '' NOT NULL,

    -- File status. 1 => "deleted", 2 => "enabled"
    `status` TINYINT(1) UNSIGNED NOT NULL,

    -- File Size.
    `size` BIGINT(50),

    -- MD5 value of file data encoded in hexadecimal.
    `etag` CHAR(32),

    -- The version of this file.
    `version` CHAR(16) NOT NULL,

    -- Who created this file.
    `created_by` varchar(65),

    -- Timestamp of create time.
    `created` BIGINT(20) UNSIGNED NOT NULL,

    -- Timestamp of update time.
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`),
    INDEX unique_file_name (`space_id`, `pid`, `name`)
) ENGINE=InnoDB COMMENT='The file schema';
