begin;

drop index if exists recipe_group_id_revision_uniq;
drop index if exists recipe_group_id_is_current_uniq;

commit;