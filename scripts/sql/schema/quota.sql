
-- The table of user quota.
CREATE TABLE IF NOT EXISTS `user_quota` (
    -- The user id.
    `user_id` VARCHAR(65) NOT NULL,

    `workspace` JSON,
    `member` JSON,
    `custom_role` JSON,
    `stream_job` JSON,
    `sync_job` JSON,
    `data_source` JSON,
    `udf` JSON,
    `file` JSON,
    `flink_cluster` JSON,
    `network` JSON,
    PRIMARY KEY (`user_id`)
) ENGINE=InnoDB COMMENT='The quota limit for user level.';

-- The table of workspace quota.
CREATE TABLE IF NOT EXISTS `workspace_quota` (
    `space_id` CHAR(20) NOT NULL,

    `workspace` JSON,
    `member` JSON,
    `custom_role` JSON,
    `stream_job` JSON,
    `sync_job` JSON,
    `data_source` JSON,
    `udf` JSON,
    `file` JSON,
    `flink_cluster` JSON,
    `network` JSON,
    PRIMARY KEY (`space_id`)
) ENGINE=InnoDB COMMENT='The quota limit for workspace level.';

-- The table of user quota.
CREATE TABLE IF NOT EXISTS `member_quota` (
    `space_id` CHAR(20) NOT NULL,
    -- The user id of member.
    `user_id` VARCHAR(65) NOT NULL,

    `workspace` JSON,
    `member` JSON,
    `custom_role` JSON,
    `stream_job` JSON,
    `sync_job` JSON,
    `data_source` JSON,
    `udf` JSON,
    `file` JSON,
    `flink_cluster` JSON,
    `network` JSON,
    PRIMARY KEY (`space_id`, `user_id`)
) ENGINE=InnoDB COMMENT='The quota limit for member level.';
