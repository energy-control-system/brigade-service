insert into brigades (status)
values (1)
returning id, status, created_at, updated_at;
