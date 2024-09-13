create table "jwtkeys" (
    "id" bigserial primary key,
    "public_key" text not null,
    "private_key" text not null,
    "algorithm" text not null,
    "is_active" bool default true,
    "expires_at" timestamp with time zone,
    "created_at" timestamp with time zone not null default current_timestamp,
    "updated_at" timestamp with time zone
);

create index idx_active_jwtkeys on "jwtkeys" ("is_active");