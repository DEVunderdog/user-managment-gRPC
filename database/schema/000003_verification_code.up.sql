create table "verification_codes" (
    "id" bigserial primary key,
    "user_id" bigint not null,
    "code" varchar(50) not null,
    "expires_at" timestamp with time zone not null,
    "is_used" bool not null default false,
    "created_at" timestamp with time zone not null default current_timestamp,
    "updated_at" timestamp with time zone
);

alter table "verification_codes" add foreign key ("user_id") references "users" ("id");
create index idx_user_is_used on "verification_codes" ("user_id", "is_used");