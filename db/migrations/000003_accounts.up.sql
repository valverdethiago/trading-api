CREATE TABLE IF NOT EXISTS account
(
  account_uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
  username TEXT NOT NULL,
  email TEXT NOT NULL,
  created_date TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
  updated_date TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
  created_by TEXT,
  updated_by TEXT,
  PRIMARY KEY(account_uuid),
  UNIQUE(username)
);