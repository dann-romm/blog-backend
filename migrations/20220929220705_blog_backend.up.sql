-- migration up file for blog_backend database

-- create users table
create table users
(
    id                       serial primary key,
    name                     varchar(255)               not null,
    username                 varchar(255)               not null unique,
    password                 varchar(255)               not null,
    email                    varchar(255)               not null,
    created_at               timestamp    default now() not null,
    updated_at               timestamp    default now() not null,
    role                     varchar(255)               not null,
    description              varchar(255) default ''    not null,
    articles_count           int          default 0     not null,
    comments_count           int          default 0     not null,
    favorites_articles_count int          default 0     not null,
    favorites_comments_count int          default 0     not null,
    followers_count          int          default 0     not null,
    followings_count         int          default 0     not null
);

-- create articles table
create table articles
(
    id               serial primary key,
    author_id        int                     not null,
    title            varchar(255)            not null,
    description      varchar(255)            not null,
    content          text                    not null,
    created_at       timestamp default now() not null,
    updated_at       timestamp default now() not null,
    views_count      int       default 0     not null,
    comments_count   int       default 0     not null,
    favorites_count  int       default 0     not null,
    votes_up_count   int       default 0     not null,
    votes_down_count int       default 0     not null,
    foreign key (author_id) references users (id)
);

-- create comments table
create table comments
(
    id               serial primary key,
    author_id        int                     not null,
    article_id       int                     not null,
    parent_id        int       default null,
    content          text                    not null,
    created_at       timestamp default now() not null,
    updated_at       timestamp default now() not null,
    votes_up_count   int       default 0     not null,
    votes_down_count int       default 0     not null,
    foreign key (author_id) references users (id),
    foreign key (article_id) references articles (id)
);

-- create users_followers table
create table users_followers
(
    id           serial primary key,
    follower_id  int not null,
    following_id int not null,
    foreign key (follower_id) references users (id),
    foreign key (following_id) references users (id)
);

-- create users_articles_favorites table
create table users_articles_favorites
(
    id         serial primary key,
    user_id    int not null,
    article_id int not null,
    foreign key (user_id) references users (id),
    foreign key (article_id) references articles (id)
);

-- create users_comments_favorites table
create table users_comments_favorites
(
    id         serial primary key,
    user_id    int not null,
    comment_id int not null,
    foreign key (user_id) references users (id),
    foreign key (comment_id) references comments (id)
);

-- create tags table
create table tags
(
    id          serial primary key,
    description varchar(255) not null
);

-- create articles_tags table
create table articles_tags
(
    id         serial primary key,
    article_id int not null,
    tag_id     int not null,
    foreign key (article_id) references articles (id),
    foreign key (tag_id) references tags (id)
);

-- create votes_articles_up table
create table votes_articles_up
(
    id         serial primary key,
    user_id    int not null,
    article_id int not null,
    foreign key (user_id) references users (id),
    foreign key (article_id) references articles (id)
);

-- create votes_articles_down table
create table votes_articles_down
(
    id         serial primary key,
    user_id    int not null,
    article_id int not null,
    foreign key (user_id) references users (id),
    foreign key (article_id) references articles (id)
);

-- create votes_comments_up table
create table votes_comments_up
(
    id         serial primary key,
    user_id    int not null,
    comment_id int not null,
    foreign key (user_id) references users (id),
    foreign key (comment_id) references comments (id)
);

-- create votes_comments_down table
create table votes_comments_down
(
    id         serial primary key,
    user_id    int not null,
    comment_id int not null,
    foreign key (user_id) references users (id),
    foreign key (comment_id) references comments (id)
);

-- create articles_views table
create table articles_views
(
    id         serial primary key,
    user_id    int not null,
    article_id int not null,
    foreign key (user_id) references users (id),
    foreign key (article_id) references articles (id)
);
