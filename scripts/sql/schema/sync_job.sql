
-- The table of sync job.
CREATE TABLE IF NOT EXISTS `sync_job` (
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- Job ID, unique within a region
    `id` CHAR(20) NOT NULL,

    -- The workflow version id
    `version` CHAR(16) NOT NULL,

    -- PID is the parent id(directory). pid is "" means root(`/`)
    `pid` CHAR(20) NOT NULL,

    -- IsDirectory represents this job whether a directory.
    `is_directory` BOOL,

    -- Job Name, Unique within a workspace.
    -- The max length of use set is 128. The system will be auto rename to <name>.<id> when deleted.
    -- Thus the VARCHAR should be define as 149 (128 + 20 + 1)
    `name` VARCHAR(149) NOT NULL,

    -- Job description
    `desc` VARCHAR(1024) CHARACTER SET utf8mb4 DEFAULT '' NOT NULL,

    -- Job type, 0 = "NoType", 1 => "StreamOperator" 2 => "StreamSQL" 3 => "StreamJAR" 4 => "StreamPython" 5 => "StreamScala"
    `type` TINYINT(1) UNSIGNED NOT NULL,

    -- Workspace status, 1 => "deleted", 2 => "enabled"
    `status` TINYINT(1) UNSIGNED DEFAULT 1 NOT NULL,

    -- User ID of created this job.
    `created_by` VARCHAR(65) NOT NULL,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,
    -- Timestamp of update time, Update when some changed, default value should be same as "created"
    `updated` BIGINT(20) UNSIGNED NOT NULL,
    -- sync source type
    `source_type` TINYINT(2) UNSIGNED DEFAULT 1 NOT NULL,
    -- sync target description
    `target_type` TINYINT(2) UNSIGNED DEFAULT 1 NOT NULL,
    PRIMARY KEY (`id`, `version`),
    UNIQUE KEY unique_job_name (`space_id`, `version`, `name`)

) ENGINE=InnoDB COMMENT='The sync job info.';

-- The table of sync job property.
CREATE TABLE IF NOT EXISTS `sync_job_property` (
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- Job ID it belongs to
    `id` CHAR(20) NOT NULL,

    -- Release version, unique
    `version` CHAR(16) NOT NULL,

    -- The environment conf that format with JSON.
    `conf` JSON,

    -- The schedule property that format with JSON.
    `schedule` JSON,

    PRIMARY KEY (`id`, `version`)

) ENGINE=InnoDB COMMENT='The meta of sync workflow.';

-- The table of sync job release.
CREATE TABLE IF NOT EXISTS `sync_job_release` (
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- Job ID it belongs to
    `id` CHAR(20) NOT NULL,

    -- The release version
    `version` CHAR(16) NOT NULL,

    -- Job Name, Unique within a workspace
    -- The max length of use set is 128. The system will be auto rename to <name>.<id> when deleted.
    -- Thus the VARCHAR should be define as 149 (128 + 20 + 1)
    `name` VARCHAR(149) NOT NULL,

    -- Job type, 1 => "StreamOperator" 2 => "StreamSQL" 3 => "StreamJAR" 4 => "StreamPython" 5 => "StreamScala"
    `type` TINYINT(1) UNSIGNED NOT NULL,

    -- Release status, 1 => "Deleted", 2 => "Inline", 3 => "Offline",
    `status` TINYINT(1) UNSIGNED DEFAULT 1 NOT NULL,

    -- Job release description
    `desc` VARCHAR(1024) CHARACTER SET utf8mb4 DEFAULT '' NOT NULL,

    -- User ID of release this job.
    `created_by` VARCHAR(65) NOT NULL,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    -- Timestamp of update time, Update when some changed, default value should be same as "created"
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`),
    INDEX mul_list_record_by_space_id(`space_id`)

) ENGINE=InnoDB COMMENT='The sync job release latest info';

-- The table of sync job versions.
-- create table table_name_new like table_name_old;
CREATE TABLE IF NOT EXISTS `sync_job_version` like `sync_job`;

-- The table of sync job meta version.
CREATE TABLE IF NOT EXISTS `sync_job_property_version` like `sync_job_property`;

create table if not exists `sync_job_connection` (
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- The id of sync jod.
    `job_id` CHAR(20) NOT NULL,

    -- The if of flink cluster.
    `cluster_id` CHAR(20) NOT NULL,

    -- The datasource id that as a sync source.
    `source_id` CHAR(20) NOT NULL,

    -- The datasource id that as a sync target.
    `target_id` CHAR(20) NOT NULL,

    -- Status, 1 => "Deleted", 2 => "Enabled"
    `status` TINYINT(1) UNSIGNED DEFAULT 0 NOT NULL,

    -- result, 1-=> success 2 => failed
    `result` TINYINT(1) UNSIGNED DEFAULT 0 NOT NULL,

    -- Message is the reason when connection failure.
    `message` VARCHAR(1024) NOT NULL,

    -- Use time. unit in ms.
    `elapse` INT,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`job_id`, `cluster_id`, `source_id`, `target_id`)

) ENGINE=InnoDB COMMENT='The connection info for sync job';

