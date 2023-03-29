DROP TABLE IF EXISTS repositories;
CREATE TABLE repositories
(
    path            text PRIMARY KEY NOT NULL,
    name            text             NOT NULL,
    git_dirty       bool             NOT NULL,
    git_branch      text             NOT NULL,
    git_action      text             NOT NULL,
    git_count_left  INTEGER          NOT NULL DEFAULT 0,
    git_count_right INTEGER          NOT NULL DEFAULT 0,
    usage_count     INTEGER          NOT NULL DEFAULT 0,
    update_time     timestamp        NOT NULL DEFAULT CURRENT_TIMESTAMP
);
