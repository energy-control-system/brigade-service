select brigade_id, inspector_id, assigned_at
from brigade_members
where brigade_id = any ($1);
