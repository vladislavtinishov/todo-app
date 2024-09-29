create table users
(
    id int auto_increment primary key,
    name varchar(255) not null,
    username varchar(255) not null,
    password_hash varchar(255) not null
);

create table todo_lists
(
    id int auto_increment primary key,
    title varchar(255) not null,
    description varchar(255)
);

create table users_lists
(
    id int auto_increment primary key,
    user_id int not null,
    list_id int not null,
    foreign key (user_id) references users(id),
    foreign key (list_id) references todo_lists(id)
);

create table todo_items
(
    id int auto_increment primary key,
    title varchar(255) not null,
    description varchar(255)
);

create table lists_items
(
    id int auto_increment primary key,
    list_id int not null,
    item_id int not null,
    foreign key (list_id) references todo_lists(id),
    foreign key (item_id) references todo_items(id)
);