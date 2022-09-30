ALTER TABLE `log` ADD `analog_input_1` FLOAT NOT NULL AFTER `speed`;
ALTER TABLE `log` ADD `payload` VARCHAR(5000) NULL AFTER `azimuth`;