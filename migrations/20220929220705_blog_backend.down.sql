-- migration down file for blog_backend database

drop table if exists articles_views;

drop table if exists votes_comments_down;

drop table if exists votes_comments_up;

drop table if exists votes_articles_down;

drop table if exists votes_articles_up;

drop table if exists articles_tags;

drop table if exists tags;

drop table if exists users_comments_favorites;

drop table if exists users_articles_favorites;

drop table if exists users_followers;

drop table if exists comments;

drop table if exists articles;

drop table if exists users;
