begin;

create unique index recipe_group_id_revision_uniq
 on recipes (recipe_group_id, revision);

create unique index recipe_group_id_is_current_uniq
on recipes(recipe_group_id, is_current) where is_current;

commit;