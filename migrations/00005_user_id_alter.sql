-- +goose Up
-- +goose StatementBegin
alter table workouts
    add column user_id bigint not null references users (id) on delete cascade;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table workouts
    drop column user_id;
-- +goose StatementEnd