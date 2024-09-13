create table "sessions" (
    "id" bigserial primary key,
    "user_id" bigint not null,
    "access_token" text not null,
    "refresh_token" text not null,
    "access_token_expires_at" timestamp with time zone not null,
    "refresh_token_expires_at" timestamp with time zone not null,
    "is_active" bool default true,
    "ip" varchar(100),
    "user_agent" varchar(255),
    "logged_out" timestamp with time zone,
    "created_at" timestamp with time zone not null default current_timestamp,
    "updated_at" timestamp with time zone
);

create index idx_users on "sessions" ("user_id");
create index idx_session_metadata on "sessions" ("ip", "user_agent");