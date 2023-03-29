-- name: GetByName :one
SELECT *
FROM repositories
where name = ?;

-- name: ListRepositories :many
SELECT *
FROM repositories
order by usage_count, lower(name) desc;

-- name: UpsertRepository :exec
insert into repositories (path, name, git_branch, git_dirty, git_action, git_count_left, git_count_right)
values (@path, @name, @git_branch, @git_dirty, @git_action, @git_count_left, @git_count_right)
on conflict (path) do update set name            = @name,
                                 git_branch      = @git_branch,
                                 git_dirty       = @git_dirty,
                                 git_action      = @git_action,
                                 git_count_left  = @git_count_left,
                                 git_count_right = @git_count_right,
                                 update_time     = CURRENT_TIMESTAMP;

-- name: IncRepoUsageCount :exec
update repositories
set usage_count = usage_count + 1,
    update_time = CURRENT_TIMESTAMP
where path = ?;

-- name: DeleteRepository :exec
delete
from repositories
where path = ?;