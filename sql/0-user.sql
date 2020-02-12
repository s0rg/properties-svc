CREATE DATABASE `usersdb`;

USE `usersdb`;


CREATE TABLE `user_settings`(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    user_id    INT NOT NULL,
    settings   VARBINARY(8192) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT NOW(),
    expires_at DATETIME
);

CREATE UNIQUE INDEX `user_settings_idx` 
    ON `user_settings`(user_id, created_at, expires_at);


CREATE USER `usr-us` IDENTIFIED BY 'usr-pw';
GRANT SELECT, INSERT, UPDATE, DELETE ON `usersdb`.* TO `usr-us`;

