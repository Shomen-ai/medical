-- Откат: очищает посеянные туркменские данные у demo-пациентов.
UPDATE users SET
    gender           = NULL,
    address          = NULL,
    id_doc_number    = NULL,
    id_doc_issued_by = NULL
WHERE id BETWEEN '33333333-0000-0000-0000-000000000001'
            AND  '33333333-0000-0000-0000-000000000010';
