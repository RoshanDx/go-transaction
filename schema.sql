CREATE TABLE IF NOT EXISTS users
(
    id         BIGSERIAL PRIMARY KEY,
    username   TEXT                           NOT NULL UNIQUE,
    firstname  TEXT                           NULL,
    activated  BOOL                           NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS roles
(
    id   BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS user_role
(
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (role_id) REFERENCES roles (id),
    UNIQUE (user_id, role_id)
);

INSERT INTO users (username, firstname, activated) VALUES ('goku', 'san', true);

INSERT INTO roles (name) VALUES ('user');
INSERT INTO roles (name) VALUES ('admin');

INSERT INTO user_role (user_id, role_id) VALUES (1, 1);