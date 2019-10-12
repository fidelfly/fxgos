CREATE SCHEMA fxgosdb CHARACTER SET utf8 COLLATE utf8_bin;

CREATE USER fxgos IDENTIFIED BY 'fxgos@lyismydg';

GRANT ALL PRIVILEGES ON fxgosdb.* TO 'fxgos'@'%' with grant option;