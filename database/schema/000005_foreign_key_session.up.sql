alter table "sessions" add foreign key ("user_id") references "users" ("id");