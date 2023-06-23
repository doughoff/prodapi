begin;

alter table production_orders add column output bigint;

commit;