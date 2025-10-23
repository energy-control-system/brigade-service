-- +goose Up
create table if not exists brigade_statuses
(
    id   int primary key generated always as identity,
    name text not null
);

insert into brigade_statuses (name)
values ('Idle'),
       ('OnTask'),
       ('Archived');

create table if not exists brigades
(
    id         int primary key generated always as identity,
    status     int         not null references brigade_statuses (id) on delete restrict,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

-- Участники бригады
create table if not exists brigade_members
(
    brigade_id   int         not null references brigades (id) on delete cascade,
    inspector_id int         not null,
    assigned_at  timestamptz not null default now(),
    primary key (brigade_id, inspector_id)
);

create index if not exists idx_brigades_status on brigades (status);

-- +goose StatementBegin
create or replace function update_updated_at_column()
    returns trigger as
$$
begin
    new.updated_at = now();
    return new;
end;
$$ language plpgsql;
-- +goose StatementEnd

create trigger trg_brigades_updated_at
    before update
    on brigades
    for each row
execute function update_updated_at_column();

-- +goose Down
drop trigger if exists trg_brigades_updated_at on brigades;
drop function if exists update_updated_at_column();
drop table if exists brigade_members;
drop table if exists brigades;
drop table if exists brigade_statuses;
