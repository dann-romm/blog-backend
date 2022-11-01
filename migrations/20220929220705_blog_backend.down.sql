-- migration down file for blog_backend database

drop table articles_views;

drop table votes_comments_down;

drop table votes_comments_up;

drop table votes_articles_down;

drop table votes_articles_up;

drop table articles_tags;

drop table tags;

drop table users_comments_favorites;

drop table users_articles_favorites;

drop table users_followers;

drop table comments;

drop table articles;

drop table users;

drop type role_type;
