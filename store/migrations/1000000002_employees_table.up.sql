CREATE TABLE "public"."employees" (
    "id" varchar(36) NOT NULL,
    "name" varchar(255) NOT NULL,
    "currency" varchar(20) NOT NULL,
    "salary" int NOT NULL,
    "department" varchar(50) NOT NULL,
    "sub_department" varchar(50) NOT NULL,
    "on_contract" boolean NOT NULL,
    "status" varchar(20) NOT NULL,
    "created_at" timestamptz NOT NULL,
    "created_by" varchar(36) NOT NULL,
    "updated_at" timestamptz NOT NULL,
    "updated_by" varchar(36) NOT NULL,
    CONSTRAINT "employees_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE INDEX "employees_department_idx"
ON "public"."employees" USING hash ("department");

CREATE INDEX "employees_on_contract_idx"
ON "public"."employees" USING hash ("on_contract");
