drop database if exists mangalist;
create database mangalist;
use mangalist;

CREATE TABLE manga (
  manga_id         INT AUTO_INCREMENT NOT NULL,
  title      VARCHAR(128) NOT NULL,
  description     VARCHAR(255) NOT NULL,
  PRIMARY KEY (manga_id, title)
);

CREATE TABLE review (
  manga_id         INT AUTO_INCREMENT NOT NULL,
  user_id    int, 
  title      VARCHAR(128) NOT NULL,
  description     VARCHAR(255) NOT NULL,
  PRIMARY KEY (manga_id, title)
);

create table author
	(A_ID			INT AUTO_INCREMENT NOT NULL,
	 a_name			varchar(20) not null, 
	 primary key (A_ID)
	);

create table user
(
user_id INT AUTO_INCREMENT NOT NULL, 
name varchar(20), 
email varchar(120), 
password varchar(60), 
primary key(user_id)
);

INSERT INTO manga
  (title, description)
VALUES
  ('Bad Anime', 'This anime is bad'),
  ('Good anime', 'This anime is good'),
  ('Whatever anime', 'This anime is good'),
  ('Jeru', 'Gerry Mulligan');

insert into author values (10101, 'Srinivasan');
insert into author values (12121, 'Wu');
insert into author values (15151, 'Mozart');
insert into author values (22222, 'Einstein');
insert into author values (32343, 'El Said');
insert into author values (33456, 'Gold');
insert into author values (45565, 'Katz');
insert into author values (58583, 'Califieri');
insert into author values (76543, 'Singh');
insert into author values (76766, 'Crick');
insert into author values (83821, 'Brandt');
insert into author values (98345, 'Kim');

insert into user (name, email, password) values ('john', 'john_doe@email.com', 'passwd');
