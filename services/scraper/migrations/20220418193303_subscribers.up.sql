CREATE TABLE IF NOT EXISTS subscribers
(
    email      string            not null,
    name       string            not null,
    subscribed integer default 1 not null
);

