select id, status, created_at, updated_at
from brigades
where id = $1
limit 1;
