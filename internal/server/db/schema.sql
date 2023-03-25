CREATE TABLE IF NOT EXISTS repositories
(
    path        text PRIMARY KEY NOT NULL,
    name        text             NOT NULL,
    git_dirty   bool             NOT NULL,
    git_branch  text             NOT NULL,
    git_action  text             NOT NULL,
    usage_count INTEGER          NOT NULL DEFAULT 0,
    update_time timestamp        NOT NULL DEFAULT CURRENT_TIMESTAMP
);