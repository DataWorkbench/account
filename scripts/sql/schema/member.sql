
-- Workspace Member
CREATE TABLE IF NOT EXISTS `member` (
    -- Only used to query sort by
    `id` BIGINT(20) NOT NULL AUTO_INCREMENT,

    -- Workspace ID that the member belongs to.
    `space_id` VARCHAR(24) NOT NULL,

    -- Use account user-id as member id.
    `user_id` VARCHAR(65) NOT NULL,

    -- Member status, 1 => "Normal" 2 => "Deleted".
    `status` TINYINT(1) UNSIGNED DEFAULT 1 NOT NULL,

    -- Workspace description
    `desc` VARCHAR(1024) CHARACTER SET utf8mb4 DEFAULT '' NOT NULL,

    -- The id lists of system role. Multiple id separated by commas, eg: "ros-1,ros-2".
    `system_role_ids` VARCHAR(256) NOT NULL,

    -- The id lists of custom role. Multiple id separated by commas, eg: "roc-1,roc-2"
    -- A member can have up to 100 custom roles.
    `custom_role_ids` VARCHAR(2048) NOT NULL,

    -- User ID of created this member.
    `created_by` VARCHAR(65) NOT NULL,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    -- Timestamp of update time, Update when some changed, default value should be same as "created"
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`),
    UNIQUE key unique_user_id(`space_id`, `user_id`)

) ENGINE =Innodb COMMENT ='The workspace member';
