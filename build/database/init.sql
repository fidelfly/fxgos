CREATE SCHEMA fxgosdb;

CREATE USER fxgos IDENTIFIED BY 'fxgos@lyismydg';

GRANT ALL PRIVILEGES ON *.* TO 'fxgos'@'%' with grant option;