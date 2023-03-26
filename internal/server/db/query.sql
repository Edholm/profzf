-- name: GetByName :one
SELECT *
FROM repositories
where name = ?;

-- name: ListRepositories :many
SELECT *
FROM repositories
order by usage_count, lower(name) desc;

-- name: UpsertRepository :exec
insert into repositories (path, name, git_branch, git_dirty, git_action)
values (@path, @name, @git_branch, @git_dirty, @git_action)
on conflict (path) do update set name        = @name,
                                 git_branch  = @git_branch,
                                 git_dirty   = @git_dirty,
                                 git_action  = @git_action,
                                 update_time = CURRENT_TIMESTAMP;

-- name: IncRepoUsageCount :exec
update repositories
set usage_count = usage_count + 1,
    update_time = CURRENT_TIMESTAMP
where path = ?;

-- name: DeleteRepository :exec
delete
from repositories
where path = ?;