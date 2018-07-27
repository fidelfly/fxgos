CREATE TABLE tenant
(
    id int PRIMARY KEY AUTO_INCREMENT,
    code nvarchar(50) NOT NULL,
    name nvarchar(100) NOT NULL
    create_time datetime NOT NULL
);
CREATE UNIQUE INDEX tenant_code_uindex ON client (code);

CREATE TABLE user
(
    id int PRIMARY KEY AUTO_INCREMENT,
    tenant_id int NOT NULL,
    code nvarchar(50) NOT NULL,
    name nvarchar(100) NOT NULL,
    password varchar(50) NOT NULL,
    create_time datetime NOT NULL,
    CONSTRAINT user_tenant_id_fk FOREIGN KEY (tenant_id) REFERENCES tenant (id)
);
CREATE UNIQUE INDEX user_code_uindex ON user (code);

CREATE TABLE trace_log
(
    id int PRIMARY KEY AUTO_INCREMENT,
    log_time datetime NOT NULL,
    user_id int NOT NULL,
    user nvarchar(50) NOT NULL,
    tenant_id int NOT NULL,
    tenant nvarchar(50) NOT NULL,
    request_url nvarchar(100),
    code nvarchar(100) NOT NULL,
    type nvarchar(50) NOT NULL,
    message nvarchar(500),
    info text
);