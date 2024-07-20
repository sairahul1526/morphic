CREATE TABLE "public"."users" (
    "id" varchar(36) NOT NULL,
    "username" varchar(255) NOT NULL,
    "password" text NOT NULL,
    "status" varchar(20) NOT NULL,
    "created_at" timestamptz NOT NULL,
    "created_by" varchar(36) NOT NULL,
    "updated_at" timestamptz NOT NULL,
    "updated_by" varchar(36) NOT NULL,
    CONSTRAINT "users_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "users_username_unique" UNIQUE ("username")
) WITH (oids = false);