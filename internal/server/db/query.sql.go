// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: query.sql

package db

import (
	"context"
)

const deleteRepository = `-- name: DeleteRepository :exec
delete
from repositories
where path = ?
`

func (q *Queries) DeleteRepository(ctx context.Context, path string) error {
	_, err := q.db.ExecContext(ctx, deleteRepository, path)
	return err
}

const getByName = `-- name: GetByName :one
SELECT path, name, git_dirty, git_branch, git_action, git_count_left, git_count_right, usage_count, update_time
FROM repositories
where name = ?
`

func (q *Queries) GetByName(ctx context.Context, name string) (Repository, error) {
	row := q.db.QueryRowContext(ctx, getByName, name)
	var i Repository
	err := row.Scan(
		&i.Path,
		&i.Name,
		&i.GitDirty,
		&i.GitBranch,
		&i.GitAction,
		&i.GitCountLeft,
		&i.GitCountRight,
		&i.UsageCount,
		&i.UpdateTime,
	)
	return i, err
}

const incRepoUsageCount = `-- name: IncRepoUsageCount :exec
update repositories
set usage_count = usage_count + 1,
    update_time = CURRENT_TIMESTAMP
where path = ?
`

func (q *Queries) IncRepoUsageCount(ctx context.Context, path string) error {
	_, err := q.db.ExecContext(ctx, incRepoUsageCount, path)
	return err
}

const listRepositories = `-- name: ListRepositories :many
SELECT path, name, git_dirty, git_branch, git_action, git_count_left, git_count_right, usage_count, update_time
FROM repositories
order by usage_count, lower(name) desc
`

func (q *Queries) ListRepositories(ctx context.Context) ([]Repository, error) {
	rows, err := q.db.QueryContext(ctx, listRepositories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Repository
	for rows.Next() {
		var i Repository
		if err := rows.Scan(
			&i.Path,
			&i.Name,
			&i.GitDirty,
			&i.GitBranch,
			&i.GitAction,
			&i.GitCountLeft,
			&i.GitCountRight,
			&i.UsageCount,
			&i.UpdateTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const upsertRepository = `-- name: UpsertRepository :exec
insert into repositories (path, name, git_branch, git_dirty, git_action, git_count_left, git_count_right)
values (@path, @name, @git_branch, @git_dirty, @git_action, @git_count_left, @git_count_right)
on conflict (path) do update set name            = @name,
                                 git_branch      = @git_branch,
                                 git_dirty       = @git_dirty,
                                 git_action      = @git_action,
                                 git_count_left  = @git_count_left,
                                 git_count_right = @git_count_right,
                                 update_time     = CURRENT_TIMESTAMP
`

type UpsertRepositoryParams struct {
	Path          string `json:"path"`
	Name          string `json:"name"`
	GitBranch     string `json:"gitBranch"`
	GitDirty      bool   `json:"gitDirty"`
	GitAction     string `json:"gitAction"`
	GitCountLeft  int64  `json:"gitCountLeft"`
	GitCountRight int64  `json:"gitCountRight"`
}

func (q *Queries) UpsertRepository(ctx context.Context, arg UpsertRepositoryParams) error {
	_, err := q.db.ExecContext(ctx, upsertRepository,
		arg.Path,
		arg.Name,
		arg.GitBranch,
		arg.GitDirty,
		arg.GitAction,
		arg.GitCountLeft,
		arg.GitCountRight,
	)
	return err
}
