drop table if exists sessions;
alter table if exists sessions
    drop constraint if exists sessions_fk;
drop table if exists verify_emails cascade;
