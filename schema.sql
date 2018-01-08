CREATE DATABASE IF NOT EXISTS `family`;

USE `family`;

CREATE TABLE `cities` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `region_id` int(11) DEFAULT NULL,
  `name` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `countries` (
  `code` char(2) NOT NULL,
  `name` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `people` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `first_name` varchar(100) NOT NULL,
  `middle_name` varchar(100) NOT NULL,
  `last_name` varchar(100) NOT NULL,
  `nick_name` varchar(100) NOT NULL,
  `mother_id` int(11) DEFAULT NULL,
  `father_id` int(11) DEFAULT NULL,
  `birth_date` date DEFAULT NULL,
  `is_alive` tinyint(1) NOT NULL DEFAULT '1',
  `home_city_id` int(11) DEFAULT NULL,
  `birth_city_id` int(11) DEFAULT NULL,
  `gender` enum('M','F') NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `regions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `country_code` char(2) DEFAULT NULL,
  `code` varchar(8) DEFAULT NULL,
  `name` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `countryCodeIdx` (`country_code`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `spouses` (
  `person1_id` int(11) NOT NULL,
  `person2_id` int(11) NOT NULL,
  `status` int(11) NOT NULL,
  `married_date` date DEFAULT NULL,
  PRIMARY KEY (`person1_id`,`person2_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE VIEW `city_view` AS
    SELECT
        `ci`.`id` AS `city_id`,
        `ci`.`name` AS `city_name`,
        `r`.`name` AS `region_name`,
        `r`.`code` AS `region_code`,
        `co`.`code` AS `country_code`,
        `co`.`name` AS `country_name`
    FROM
        ((`family`.`cities` `ci`
        JOIN `family`.`regions` `r` ON ((`ci`.`region_id` = `r`.`id`)))
        JOIN `family`.`countries` `co` ON ((`co`.`code` = `r`.`country_code`)));
