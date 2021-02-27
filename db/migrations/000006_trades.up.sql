CREATE TABLE IF NOT EXISTS trade
(
  trade_uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
  account_uuid UUID NOT NULL,
  symbol TEXT NOT NULL,
  quantity NUMERIC(9) NOT NULL,
  side trade_side NOT NULL,
  price NUMERIC(11,2) NOT NULL,
  status trade_status NOT NULL DEFAULT 'SUBMITTED'::trade_status,
  created_date TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
  updated_date TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
  created_by TEXT,
  updated_by TEXT,
  PRIMARY KEY(trade_uuid),
  FOREIGN KEY (account_uuid) REFERENCES account (account_uuid)
);