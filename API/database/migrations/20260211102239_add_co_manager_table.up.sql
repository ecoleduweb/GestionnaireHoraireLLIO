-- Migration Up
CREATE TABLE IF NOT EXISTS co_managers (
   project_id BIGINT(20) UNSIGNED NOT NULL,
   user_id    BIGINT(20) UNSIGNED NOT NULL,

   PRIMARY KEY (project_id, user_id),

    CONSTRAINT fk_co_managers_project
    FOREIGN KEY (project_id) REFERENCES projects(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,

    CONSTRAINT fk_co_managers_user
    FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
)