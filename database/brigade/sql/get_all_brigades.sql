select id, status, created_at, updated_at
from brigades
order by id
limit $1 offset $2;
