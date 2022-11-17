drop database if exists mangalist;
create database mangalist;
use mangalist;

CREATE TABLE manga (
  manga_id         INT AUTO_INCREMENT NOT NULL,
  title      VARCHAR(128) NOT NULL,
  description     VARCHAR(255) NOT NULL,
  PRIMARY KEY (manga_id)
);


create table reviewer
(
reviewer_id INT AUTO_INCREMENT NOT NULL, 
user_id int,
name varchar(20), 
primary key(reviewer_id)
);

create table belong(
 manga_id int,
 G_ID int, 
 primary key (manga_id, G_ID),
 foreign key (manga_id) references manga(manga_id)
 on delete cascade
);

create table request(
  request_id INT AUTO_INCREMENT NOT NULL,
  reviewer_id int not null,
  title      VARCHAR(128) NOT NULL,
  primary key (request_id)
);

create table favorite(
 reviewer_id int,
 manga_ID int, 
 primary key (reviewer_id, manga_ID)
);

CREATE TABLE review(
  review_id INT AUTO_INCREMENT NOT NULL,
  reviewer_id    INT NOT NULL, 
  manga_id int NOT NULL, 
  title      VARCHAR(128) NOT NULL,
  description     VARCHAR(255) NOT NULL,
  PRIMARY KEY (review_id)
);

CREATE TABLE genres(
  G_ID INT AUTO_INCREMENT NOT NULL, 
  G_NAME varchar(20), 
  primary key(G_ID)
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
primary key(user_id, name)
);


INSERT INTO manga
  (title, description)
VALUES
  ('Bad Manga', 'This Manga is bad'),
  ('Grand Blue', 'Acquired Taste'),
  ('Good Manga', 'This Manga is good'),
  ('Whatever Manga', 'This Manga is good'),
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

insert into genres(g_name) values ('comedy'), ('alcoholism');
insert into belong values (2, 1);
insert into belong values (2, 2);
