CREATE TABLE sbtest (
  `id` VARCHAR(255) NOT NULL,
  `driver_id` VARCHAR(255) NOT NULL DEFAULT '',
  `connected` BOOLEAN NOT NULL DEFAULT 0,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  `loc` POINT NOT NULL,
  `bearing` FLOAT NOT NULL,
  `altitude` FLOAT NOT NULL,

  `region_id` VARCHAR(255) NOT NULL,
  `queue_zone_id` VARCHAR(255) DEFAULT '',
  `queue_zone_left_id` VARCHAR(255) DEFAULT '',
  `queue_zone_left_at` DATETIME(6),
  `queue_zone_hit_at` DATETIME(6),
  `queue_zone_loc` POINT,

  `loc_ts` DATETIME(6) NOT NULL,

  `speed` FLOAT NOT NULL DEFAULT 0.0,
  `alu_serial` bigint(20) NOT NULL,

  `snoozed_at` DATETIME(6),
  
  PRIMARY KEY (`id`),
  INDEX `region_queuezone_id` USING HASH (`region_id`, `queue_zone_id`)
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;