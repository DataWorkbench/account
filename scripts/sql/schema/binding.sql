-- Table for for describes dependencies between modules.
CREATE TABLE IF NOT EXISTS `binding` (
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- module_id represents which resources are bound to this module.
    `module_id` CHAR(20) NOT NULL,

    -- module_version is the version of module.
    -- This filed maybe empty.
    `module_version` CHAR(16) NOT NULL,

    -- resource_id represents the module bound resources.
    `resource_id` CHAR(20) NOT NULL,

    -- resource_version is the version of resource.
    -- Notice: Reserved field, unused on present.
    `resource_version` CHAR(16) NOT NULL,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`space_id`, `module_id`, `module_version`, `resource_id`, `resource_version`)
    #     INDEX mul_index_with_module_id(`module_id`, `module_version`),
    #     INDEX mul_index_with_resource_id(`resource_id`, `resource_version`)

) ENGINE=InnoDB COMMENT='describes dependencies between modules.';
