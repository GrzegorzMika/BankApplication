drop table if exists users;
alter table if exists accounts drop constraint if exists owner_fk;
drop index if exists idx_owner_currency on accounts;
