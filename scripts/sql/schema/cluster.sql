
CREATE TABLE IF NOT EXISTS `flink_cluster` (
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- Cluster ID, unique within a region
    `id` CHAR(20) NOT NULL,

    -- Cluster Name, Unique within a workspace
    -- The max length of use set is 128. The system will be auto rename to <name>.<id> when deleted.
    -- Thus the VARCHAR should be define as 149 (128 + 20 + 1)
    `name` VARCHAR(149) NOT NULL,

    -- The flink version.
    `version` VARCHAR(63) NOT NULL,

    -- The cluster status. 1 => "deleted" 2 => "running" 3 => "stopped" 4 => "starting" 5 => "exception" 6 => "Arrears"
    `status` TINYINT(1) UNSIGNED NOT NULL,

    -- Flink task number for TaskManager. Is required, Min 1, Max ?
    `task_num` INT,

    -- Flink JobManager's cpu and memory. 1CU = 1C + 4GB. Is required, Min 0.5, Max 8
    `job_cu` FLOAT,

    -- Flink TaskManager's cpu and memory. 1CU = 1C + 4GB. Is required, Min 0.5, Max 8
    `task_cu` FLOAT,

    -- Network config.
    `network_id` CHAR(20) NOT NULL,

    -- Config of host aliases
    `host_aliases` JSON,

    -- Flink config.
    `config` JSON,

    -- The user-id of created this flink cluster.
    `created_by` VARCHAR(128) NOT NULL,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    -- Timestamp of update time, Update when some changed, default value should be same as "created"
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`),
    UNIQUE KEY unique_flink_cluster_name(`space_id`, `name`)

) ENGINE=InnoDB COMMENT='The flink cluster info.';

