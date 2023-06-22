begin;

alter table stock_movements
    add column document_number text;

commit;