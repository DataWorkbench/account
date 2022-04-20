-- Workspace Info
CREATE TABLE IF NOT EXISTS `workspace` (
    -- Workspace ID, unique within a region
    `id` CHAR(20) NOT NULL,

    -- User ID of workspace owner
    `owner` VARCHAR(65) NOT NULL,

    -- Workspace name, unique within a region
    -- The max length of use set is 128. The system will be auto rename to <name>.<id> when deleted.
    -- Thus the VARCHAR should be define as 149 (128 + 20 + 1)
    `name` VARCHAR(149) NOT NULL,

    -- Workspace description
    `desc` VARCHAR(1024) CHARACTER SET utf8mb4 DEFAULT '' NOT NULL,

    -- Workspace status, 1 => "enabled", 2 => "disabled", 3 => "deleted"
    `status` TINYINT(1) UNSIGNED DEFAULT 1 NOT NULL,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    -- Timestamp of update time, Update when some changed, default value should be same as "created"
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`),
    UNIQUE INDEX unique_space_name (`owner`, `name`)

) ENGINE=InnoDB COMMENT='The workspace info';

