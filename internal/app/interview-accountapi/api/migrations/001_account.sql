-- +migrate Up
CREATE TABLE IF NOT EXISTS "Account"
(
  id              UUID PRIMARY KEY NOT NULL,
  organisation_id UUID             NOT NULL,
  version         INT              NOT NULL,
  is_deleted      BOOLEAN          NOT NULL,
  is_locked       BOOLEAN          NOT NULL,
  created_on      TIMESTAMP,
  modified_on     TIMESTAMP,
  record          JSONB,
  pagination_id   SERIAL
);

CREATE UNIQUE INDEX ON "Account" (id);
CREATE UNIQUE INDEX Account_paginationid ON "Account" (pagination_id);

-- +migrate Down
DROP TABLE IF EXISTS "Account";

