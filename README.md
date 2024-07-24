# Morphic

Morphic is a simple summaries service, written in Golang.

# Table of Contents

1. [Features](#features)
2. [Requirements](#requirements)
3. [Setup](#setup)
4. [Install & Run](#install--run)
5. [Test](#test)

# Features

1. Add and delete employees.
2. List SS summaries for all, departments and sub-departments.
3. Loin user and get auth.

# Requirements

1. Go 1.20.x or later
2. Postgres (or can use test database credentials)

# Setup

1. Create a file `config.yaml` in root directory
2. Copy the content from `config/config_template.yaml` to `config.yaml`
3. Update the values in `config.yaml` as per your environment like postgres connection details, port etc.

# Install & Run

1. Run `make all` to start both the postgres and the server
2. Run the following command `psql 'postgres://postgres:wDnfWovh4uf3@localhost:5432/morphic?sslmode=disable' -f scripts/dummy_data.sql` or take the dummy data from `/scripts/dummy_data.sql` and run it by login to the database `postgres://postgres:wDnfWovh4uf3@localhost:5432/morphic?sslmode=disable`
3. Go to `http://localhost:8060/swagger/index.html` to access Swagger UI
4. Login with "dummyuser", "dummy_password" and get the access_token from header
5. Add the token as "Bearer <token>" in the swagger.
6. Now you can test all the APIs.

# Test

1. Run `make test` to run all the tests
2. Run `make test-coverage` to check the test coverage
