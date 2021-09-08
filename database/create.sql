CREATE DATABASE provisa
    WITH 
    OWNER = postgres
    ENCODING = 'UTF8'
    CONNECTION LIMIT = -1;


CREATE SCHEMA provisa
    AUTHORIZATION postgres;


CREATE TABLE provisa.role
(
    id integer NOT NULL,
    name character varying(20) NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE provisa.role
    OWNER to postgres;


CREATE TABLE provisa."user"
(
    id integer NOT NULL,
    email character varying(50) NOT NULL,
    password character varying(100) NOT NULL,
    created timestamp with time zone NOT NULL,
    deleted timestamp with time zone,
    deleted_by integer,
    last_login timestamp with time zone,
    first_name character varying(20) NOT NULL,
    last_name character varying(20) NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE provisa."user"
    OWNER to postgres;

CREATE UNIQUE INDEX idx_user_email
    ON provisa."user" USING btree
    (email COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;
    
CREATE TABLE provisa.user_role
(
    user_id integer NOT NULL,
    role_id integer NOT NULL,
    PRIMARY KEY (user_id, role_id),
    CONSTRAINT "user" FOREIGN KEY (user_id)
        REFERENCES provisa."user" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT role FOREIGN KEY (role_id)
        REFERENCES provisa.role (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

ALTER TABLE provisa.user_role
    OWNER to postgres;

CREATE TABLE provisa.client
(
    id integer NOT NULL,
    full_name character varying(20) NOT NULL,
    created timestamp with time zone NOT NULL,
    deleted timestamp with time zone,
    deleted_by integer,
    last_login time with time zone,
    PRIMARY KEY (id)
);

ALTER TABLE provisa.client
    OWNER to postgres;


insert into provisa."role" (id,name) values (1, 'admin');
insert into provisa."role" (id,name) values (2, 'manager');
insert into provisa."role" (id,name) values (3, 'customer');
commit;

insert into provisa."user" (id,email,password,created,first_name, last_name) values (1, 'novik.g.v@gmail.com', '1', CURRENT_TIMESTAMP, 'Gennady', 'Novik');
insert into provisa."user" (id,email,password,created,first_name, last_name) values (2, 'ngorov@gmail.com', '1', CURRENT_TIMESTAMP, 'Nikolai', 'Gorov');
commit;

insert into provisa."user_role" (user_id,role_id) values (1, 1);
insert into provisa."user_role" (user_id, role_id) values (2,1);
commit;	