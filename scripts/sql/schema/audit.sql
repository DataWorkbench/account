-- Workspace Operation Audit
CREATE TABLE IF NOT EXISTS `audit` (
    -- Only used to query sort by
    `id` BIGINT(20) NOT NULL AUTO_INCREMENT,

    -- The user id of who execute this operation.
    `user_id`  VARCHAR(65) NOT NULL,

    -- The workspace id
    `space_id` VARCHAR(24),

    -- The type of operation permission,  1 => "Write", 2 => "Read".
    `perm_type` TINYINT(1) UNSIGNED NOT NULL,

    -- The operation of user behavior.
    `api_name` VARCHAR(128) NOT NULL,

    -- The operation state, 1 => "Success", 2 => "Failed".
    `state` TINYINT(1) UNSIGNED DEFAULT 0 NOT NULL,

    -- Timestamp of time of when accessed.
    `created` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`),
    -- Index to list all records of the specified user.
    INDEX mul_list_audit_by_user_id(`user_id`),
    -- Index to lists all records of the specified workspace id.
    INDEX mul_list_audit_by_space_id(`space_id`)
) ENGINE=InnoDB COMMENT='The workspace operation opaudit record';
