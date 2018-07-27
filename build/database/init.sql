CREATE SCHEMA m18awsadmin;

CREATE USER m18aws IDENTIFIED BY 'm18aws@7jiaj';

GRANT ALL PRIVILEGES ON *.* TO 'm18aws'@'%' with grant option;