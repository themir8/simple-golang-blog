create table simpleGolangBlog (
  id serial primary key,
  title varchar not null,
  anons text not null,
  full_text text not null
);