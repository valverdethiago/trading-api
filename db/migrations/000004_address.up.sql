CREATE TABLE IF NOT EXISTS address
(
  address_uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  street TEXT NOT NULL,
  city TEXT NOT NULL,
  state state NOT NULL,
  zipcode TEXT NOT NULL,
  account_uuid UUID, 
  created_date TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
  updated_date TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
  created_by TEXT,
  updated_by TEXT,
  PRIMARY KEY(address_uuid),
  UNIQUE(account_uuid),
  FOREIGN KEY (account_uuid) REFERENCES account (account_uuid)
);