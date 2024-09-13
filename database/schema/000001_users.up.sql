create table "users" (
    "id" bigserial primary key,
    "email" varchar(100) not null,
    "hashed_password" text not null,
    "email_verified" bool not null,
    "created_at" timestamp with time zone not null default current_timestamp,
    "updated_at" timestamp with time zone
);

create index idx_email on "users" ("email");