/*create tables*/

CREATE TABLE _user (
    username varChar(60) PRIMARY KEY NOT NULL,
    user_given_name varChar(70) NOT NULL,
    user_family_name varChar(70),
    user_email varChar(80) NOT NULL,
    user_password varChar(80) NOT NULL

);

CREATE TABLE _plan (
    plan_id serial PRIMARY KEY NOT NULL,
    plan_title varChar(100) NOT NULL,
    plan_completed boolean NOT NULL,
    plan_owner varChar(60)
);

CREATE TABLE _item (
    item_id serial PRIMARY KEY NOT NULL,
    item_title varChar(100) NOT NULL,
    item_description varChar(200),
    item_completed boolean NOT NULL,
    plan_id int
);

CREATE TABLE _privileges (
    privilege_id serial PRIMARY KEY NOT NULL,
    plan_id int NOT NULL,
    username varChar(60) NOT NULL,
    write boolean NOT NULL

);

ALTER TABLE _plan ADD
CONSTRAINT FK_username FOREIGN KEY (plan_owner) REFERENCES _user(username) ON DELETE CASCADE;

ALTER TABLE _item ADD
CONSTRAINT FK_plan_id FOREIGN KEY (plan_id) REFERENCES _plan(plan_id) ON DELETE CASCADE;

ALTER TABLE _privileges ADD
CONSTRAINT FK_Username FOREIGN KEY (username) REFERENCES _user(username) ON DELETE CASCADE;

ALTER TABLE _privileges ADD
CONSTRAINT FK_plan_id FOREIGN KEY (plan_id) REFERENCES _plan(plan_id) ON DELETE CASCADE;
