create table if not exists data_source (
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- DataSource ID
    `id` CHAR(20) NOT NULL,

    -- unique in a workspace.
    -- The max length of use set is 64. The system will be auto rename to <name>.<id> when deleted.
    -- Thus the VARCHAR should be define as 85 (64 + 20 + 1)
    `name` VARCHAR(85) NOT NULL,

    -- DataSource description.
    `desc` varchar(256),

    -- Type, 1->MySQL 2->PostgreSQL 3->Kafka 4->S3 5->ClickHouse 6->Hbase 7->Ftp 8->HDFS
    `type` TINYINT(1) UNSIGNED DEFAULT 0 NOT NULL,

    -- URL of data source settings..
    `url` JSON,

    -- Status, 1 => "Delete", 2 => "enabled", 3 => "disabled"
    `status` TINYINT(1) UNSIGNED DEFAULT 0 NOT NULL,

    -- User ID of created this data source.
    `created_by` VARCHAR(65) NOT NULL,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    -- Timestamp of update time, Update when some changed, default value should be same as "created"
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`),
    UNIQUE KEY unique_source_name (`space_id`, `name`)
) ENGINE=InnoDB COMMENT='The data source info';

create table if not exists data_source_connection (
    -- Only used to query sort by
    `id` BIGINT(20) NOT NULL AUTO_INCREMENT,

    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- DataSource ID
    `source_id` CHAR(20) NOT NULL,

    -- Network ID
    `network_id` CHAR(20) NOT NULL,

    -- Status, 1 => "Delete", 2 => "Enabled"
    `status` TINYINT(1) UNSIGNED DEFAULT 0 NOT NULL,

    -- result, 1-=> success 2 => failed
    `result` TINYINT(1) UNSIGNED DEFAULT 0 NOT NULL,

    -- Message is the reason when connection failure.
    `message` VARCHAR(1024) NOT NULL,

    -- Use time. unit in ms.
    `elapse` INT,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`),
    INDEX mul_query_with_space_id(`space_id`),
    INDEX mul_query_with_source_id(`source_id`)

) ENGINE=InnoDB COMMENT='The connection info for datasource';
