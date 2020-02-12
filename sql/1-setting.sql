CREATE DATABASE `settingsdb`;

USE `settingsdb`;

-- settings

CREATE TABLE `settings`(
    id        INT AUTO_INCREMENT PRIMARY KEY,
    name      VARCHAR(255) NOT NULL UNIQUE,
    notify    BIT(1) NOT NULL
);

-- settings_values

CREATE TABLE `settings_values`(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    setting_id INT NOT NULL,
    name       VARCHAR(255) NOT NULL,
    value      VARCHAR(255) NOT NULL
);

CREATE INDEX `settings_values_setting_id`
    ON `settings_values`(setting_id);

-- bundles

CREATE TABLE `bundles`(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    parent_id  INT NOT NULL DEFAULT 0,
    tag        VARCHAR(255),
    name       VARCHAR(255) NOT NULL
);

CREATE UNIQUE INDEX `bundles_enabled_idx`
    ON `bundles`(tag, name);

-- bundles_values

CREATE TABLE `bundles_values`(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    bundle_id  INT NOT NULL,
    value_id   INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT NOW(),
    expired_at DATETIME
);

CREATE INDEX `bundles_values_idx`
    ON `bundles_values`(bundle_id, created_at, expired_at);

    
CREATE USER `set-us` IDENTIFIED BY 'set-pw';
GRANT SELECT, INSERT, UPDATE, DELETE ON `settingsdb`.* TO `set-us`;

