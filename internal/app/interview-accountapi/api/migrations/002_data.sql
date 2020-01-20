-- +migrate Up

INSERT INTO "Account"
(id,organisation_id,version,is_deleted,is_locked,created_on,modified_on,record)
VALUES (
    lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6))),
    lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6))),
    0,
    false,
    false,
    datetime('now', 'localtime'),
    datetime('now', 'localtime'),
    '{
           "country": "GB",
           "base_currency": "GBP",
           "account_number": "41426819",
           "bank_id": "400300",
           "bank_id_code": "GBDSC",
           "bic": "NWBKGB22",
           "iban": "GB11NWBK40030041426819",
           "title": "Ms",
           "first_name": "Samantha",
           "bank_account_name": "Samantha Holder",
           "alternative_bank_account_names": [
             "Sam Holder"
           ],
           "account_classification": "Personal",
           "joint_account": false,
           "account_matching_opt_out": false,
           "secondary_identification": "A1B2C3D4"
     }'),(
    lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6))),
    lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6))),
    0,
    false,
    false,
    datetime('now', 'localtime'),
    datetime('now', 'localtime'),
    '{
           "country": "GB",
           "base_currency": "GBP",
           "account_number": "51426819",
           "bank_id": "400300",
           "bank_id_code": "GBDSC",
           "bic": "NWBKGB22",
           "iban": "GB11NWBK40030041426819",
           "title": "Mr",
           "first_name": "Barry",
           "bank_account_name": "White",
           "alternative_bank_account_names": [
             "Baz White"
           ],
           "account_classification": "Personal",
           "joint_account": false,
           "account_matching_opt_out": false,
           "secondary_identification": "JJZDEDE"
     }');

-- +migrate Down
DELETE FROM "Accounts";