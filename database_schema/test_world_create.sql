SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

DROP SCHEMA IF EXISTS `test_world` ;
CREATE SCHEMA IF NOT EXISTS `test_world` DEFAULT CHARACTER SET latin1 ;
USE `test_world` ;

-- -----------------------------------------------------
-- Table `test_world`.`esriHeader`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `test_world`.`esriHeader` ;

CREATE TABLE IF NOT EXISTS `test_world`.`esriHeader` (
  `id_esriHeader` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  `nCols` BIGINT(20) NOT NULL,
  `nRows` BIGINT(20) NOT NULL,
  `xllCorner` FLOAT NOT NULL,
  `yllCorner` FLOAT NOT NULL,
  `cellSize` FLOAT NOT NULL,
  `noDataValue` FLOAT NULL DEFAULT NULL,
  PRIMARY KEY (`id_esriHeader`),
  UNIQUE INDEX `id_esriHeader_UNIQUE` (`id_esriHeader` ASC))
ENGINE = InnoDB
AUTO_INCREMENT = 20
DEFAULT CHARACTER SET = latin1;


-- -----------------------------------------------------
-- Table `test_world`.`esriData`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `test_world`.`esriData` ;

CREATE TABLE IF NOT EXISTS `test_world`.`esriData` (
  `id_esriData` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `lat` FLOAT NOT NULL,
  `lon` FLOAT NOT NULL,
  `value` FLOAT NULL DEFAULT NULL,
  `esriHeader_id_esriHeader` BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (`id_esriData`),
  INDEX `fk_esriData_esriHeader_idx` (`esriHeader_id_esriHeader` ASC),
  UNIQUE INDEX `id_esriData_UNIQUE` (`id_esriData` ASC),
  CONSTRAINT `fk_esriData_esriHeader`
    FOREIGN KEY (`esriHeader_id_esriHeader`)
    REFERENCES `test_world`.`esriHeader` (`id_esriHeader`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = latin1;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
