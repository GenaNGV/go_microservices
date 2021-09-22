CREATE TABLE provisa.job_status
(
    id   integer               NOT NULL,
    name character varying(20) NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE provisa.job_status
    OWNER to postgres;


CREATE TABLE provisa.job_info
(
    id            SERIAL                   NOT NULL,
    created       timestamp with time zone NOT NULL,
    file_name     character varying(50)    NOT NULL,
    created_by    integer,
    job_status_id integer                  NOT NULL,
    finished      timestamp with time zone,
    PRIMARY KEY (id),
    CONSTRAINT "job_status_fk1" FOREIGN KEY (job_status_id)
        REFERENCES provisa."job_status" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

ALTER TABLE provisa.job_info
    OWNER to postgres;


CREATE TABLE provisa.job_statistic
(
    id          SERIAL                NOT NULL,
    job_info_id integer               NOT NULL,
    count       integer               NOT NULL,
    term        character varying(20) NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT "job_statistic_fk1" FOREIGN KEY (job_info_id)
        REFERENCES provisa."job_info" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID

);

ALTER TABLE provisa.job_statistic
    OWNER to postgres;


insert into provisa."job_status" (id, name)
values (1, 'assigned');
insert into provisa."job_status" (id, name)
values (2, 'processing');
insert into provisa."job_status" (id, name)
values (3, 'finished');
insert into provisa."job_status" (id, name)
values (4, 'error');
commit;
