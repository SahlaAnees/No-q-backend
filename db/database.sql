CREATE DATABASE IF NOT EXISTS db;

use db;

CREATE TABLE IF NOT EXISTS category (
    category varchar(120) NOT NULL primary key,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO category (category) VALUES ("Health");
INSERT INTO category (category) VALUES ("Beauty");
INSERT INTO category (category) VALUES ("Food");
INSERT INTO category (category) VALUES ("Home & Deco");
INSERT INTO category (category) VALUES ("Animal Care");
INSERT INTO category (category) VALUES ("Theater");

CREATE TABLE IF NOT EXISTS merchant (
    id int unsigned NOT NULL auto_increment primary key,
    category varchar(120) NOT NULL,
    name varchar(120) NOT NULL,
    email varchar(120),
    password varchar(1024),
    facebook varchar(255) NOT NULL DEFAULT '',
    instagram varchar(255) NOT NULL DEFAULT '',
    website varchar(255) NOT NULL DEFAULT '',
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT merchant_category_fk FOREIGN KEY (category) REFERENCES category (category) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS token (
    token_id int unsigned NOT NULL auto_increment primary key,
    merchant_id int unsigned NOT NULL,
    auth_token varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT token_merchant_fk FOREIGN KEY (merchant_id) REFERENCES merchant (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS queue (
    id int unsigned NOT NULL auto_increment primary key,
    merchant_id int unsigned NOT NULL,
    name varchar(120) NOT NULL,
    intervals int unsigned NOT NULL,
    start_time timestamp NOT NULL,
    end_time timestamp NOT NULL,
    is_available tinyint(1) NOT NULL DEFAULT "0",
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY merchant_queue_name (merchant_id, name),
    CONSTRAINT queue_merchant_fk FOREIGN KEY (merchant_id) REFERENCES merchant (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS unavailable (
    id int unsigned NOT NULL auto_increment primary key,
    queue_id int unsigned NOT NULL,
    date timestamp NOT NULL,
    CONSTRAINT unavailable_queue_fk FOREIGN KEY (queue_id) REFERENCES queue (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user (
    id int unsigned NOT NULL auto_increment primary key,
    name varchar(120) NOT NULL,
    phone varchar(10) NOT NULL,
    email varchar(120),
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS reserved_slots (
    token_no int unsigned NOT NULL auto_increment primary key,
    queue_id int unsigned NOT NULL,
    start_time timestamp NOT NULL,
    end_time timestamp NOT NULL,
    reserved_by int unsigned NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT slot_user_fk FOREIGN KEY (reserved_by) REFERENCES user (id) ON DELETE CASCADE
);