
CREATE TABLE IF NOT EXISTS `network` (
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- Network ID, unique within a region
    `id` CHAR(20) NOT NULL,

    -- Network Name, Unique within a workspace
    -- The max length of use set is 128. The system will be auto rename to <name>.<id> when deleted.
    -- Thus the VARCHAR should be define as 149 (128 + 20 + 1)
    `name` VARCHAR(149) NOT NULL,

    -- VPC's route_id.
    `router_id` VARCHAR(32) NOT NULL,

    -- VPC's vxnet_id.
    `vxnet_id` VARCHAR(32) NOT NULL,

    -- The user-id of created this network.
    `created_by` VARCHAR(128) NOT NULL,

    -- The cluster status. 1 => "deleted" 2 => "Enabled"
    `status` TINYINT(1) UNSIGNED NOT NULL DEFAULT 1,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    -- Timestamp of update time, Update when some changed, default value should be same as "created"
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`),
    UNIQUE KEY unique_network_name (`space_id`, `name`)

) ENGINE=InnoDB COMMENT='The network info.';

