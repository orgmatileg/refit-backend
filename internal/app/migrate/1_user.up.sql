CREATE TABLE `user`
( 
  `id`                    INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `email`                 VARCHAR(30) NOT NULL,
  `password`              VARCHAR(25) NOT NULL,
  `created_at`            DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at`            DATETIME ON UPDATE CURRENT_TIMESTAMP, 
  CONSTRAINT id  PRIMARY KEY (id)
);