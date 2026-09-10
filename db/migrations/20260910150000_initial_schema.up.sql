CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    username varchar NOT NULL,
    email varchar NOT NULL,
    password varchar NOT NULL,
    status varchar NOT NULL DEFAULT 'active'
);

CREATE UNIQUE INDEX IF NOT EXISTS users_email_key ON users (email);

CREATE TABLE IF NOT EXISTS projects (
    id uuid PRIMARY KEY,
    name varchar NOT NULL,
    slug varchar NOT NULL,
    goal varchar,
    shouldbe varchar,
    color varchar,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    deadline timestamptz,
    starting_at timestamptz,
    started_at timestamptz,
    finished_at timestamptz,
    sort_order bigint NOT NULL,
    user_id uuid NOT NULL REFERENCES users (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS projects_slug_key ON projects (slug);

CREATE TABLE IF NOT EXISTS tasks (
    id uuid PRIMARY KEY,
    title varchar NOT NULL,
    description varchar NOT NULL DEFAULT '',
    status varchar NOT NULL DEFAULT 'scheduled',
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    finished_at timestamptz,
    started_at timestamptz,
    archived_at timestamptz,
    starting_at timestamptz,
    deadline timestamptz,
    project_id uuid NOT NULL REFERENCES projects (id)
);
