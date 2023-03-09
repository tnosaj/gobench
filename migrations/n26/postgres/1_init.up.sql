CREATE TYPE account_status AS ENUM ('OPEN', 'CLOSED');

CREATE TYPE card_type AS ENUM ('MASTERCARD', 'MAESTRO');

CREATE TYPE address_type AS ENUM ('SHIPPING', 'WORKING','PASSPORT','EMERGENCY');

CREATE  TABLE users ( 
	id                   uuid  PRIMARY KEY  ,
	first_name           varchar(100)  NOT NULL  ,
	last_name            varchar(100)  NOT NULL  ,
	birth_date           date  NOT NULL  ,
	birth_place          varchar(100)    ,
	gender               boolean  NOT NULL  ,
	nationality          char(3)  NOT NULL  ,
	email                varchar(320)  NOT NULL  ,
	mobile_phone_number  varchar(15)  NOT NULL  ,
	created              bigint DEFAULT extract(epoch from CURRENT_TIMESTAMP)*1000 NOT NULL  ,
	updated              bigint DEFAULT extract(epoch from CURRENT_TIMESTAMP)*1000 NOT NULL
 );

CREATE  TABLE addresses ( 
	user_id              uuid REFERENCES users( id ) NOT NULL,
  address_type         address_type NOT NULL,
	street               varchar(100)  NOT NULL  ,
	house_number_block   varchar(100)  NOT NULL  ,
	zip_code             varchar(10)  NOT NULL  ,
	city                 varchar(100)  NOT NULL  ,
	country              varchar(100)  NOT NULL  ,
	created              bigint DEFAULT extract(epoch from CURRENT_TIMESTAMP)*1000 NOT NULL  ,
	updated              bigint DEFAULT extract(epoch from CURRENT_TIMESTAMP)*1000 NOT NULL  ,
	PRIMARY KEY ( user_id, address_type )
 );


CREATE  TABLE accounts ( 
	id                   uuid  PRIMARY KEY  ,
	user_id              uuid REFERENCES users( id ) NOT NULL,
	currency             char(3)  NOT NULL  ,
	status               account_status DEFAULT 'OPEN' NOT NULL  ,
	iban                 char(34)    ,
	balance              numeric NOT NULL  ,
	created              bigint DEFAULT extract(epoch from CURRENT_TIMESTAMP)*1000 NOT NULL  ,
	updated              bigint DEFAULT extract(epoch from CURRENT_TIMESTAMP)*1000 NOT NULL
 );


CREATE  TABLE cards ( 
	id                   uuid  PRIMARY KEY  ,
	account_id           uuid REFERENCES accounts( id ) NOT NULL,
	type                 card_type NOT NULL  ,
	expiration_date      date  NOT NULL  ,
	masked_card_number   char(19)  NOT NULL  ,
	created              bigint DEFAULT extract(epoch from CURRENT_TIMESTAMP)*1000 NOT NULL  ,
	updated              bigint DEFAULT extract(epoch from CURRENT_TIMESTAMP)*1000 NOT NULL
 );


CREATE  TABLE transactions ( 
	id                   uuid  PRIMARY KEY  ,
	card_id              uuid  REFERENCES cards( id ) NOT NULL,
	account_id           uuid  REFERENCES accounts( id ) NOT NULL,
	amount               numeric  NOT NULL  ,
	created              bigint DEFAULT extract(epoch from CURRENT_TIMESTAMP)*1000 NOT NULL
 );

CREATE OR REPLACE FUNCTION transactions_insert_function()
RETURNS TRIGGER AS $$
DECLARE
	partition_date TEXT;
	partition_name TEXT;
	start_of_month TEXT;
	end_of_next_month TEXT;
  transaction_date DATE;
BEGIN
  transaction_date := to_timestamp(NEW.created::double precision /1000);
	partition_date := to_char(transaction_date,'YYYY_MM');
 	partition_name := 'transactions_' || partition_date;
	start_of_month := to_char(transaction_date,'YYYY-MM') || '-01';
	end_of_next_month := to_char((transaction_date + interval '1 month'),'YYYY-MM') || '-01';
IF NOT EXISTS
	(SELECT 1
   	 FROM   information_schema.tables 
   	 WHERE  table_name = partition_name) 
THEN
	RAISE NOTICE 'A partition has been created %', partition_name;
	EXECUTE format(E'CREATE TABLE %I (CHECK ( date_trunc(\'day\', to_timestamp(created::double precision /1000)) >= ''%s'' AND date_trunc(\'day\', to_timestamp(created::double precision /1000)) < ''%s'')) INHERITS (transactions)', partition_name, start_of_month,end_of_next_month);
	-- EXECUTE format('GRANT SELECT,INSERT,DROP ON TABLE %I TO leslie_lamport', partition_name);
END IF;
EXECUTE format('INSERT INTO %I (id,card_id,account_id,amount,created) VALUES($1,$2,$3,$4,$5)', partition_name) using NEW.id, NEW.card_id, NEW.account_id, NEW.amount, NEW.created;
RETURN NULL;
END
$$
LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER insert_transactions_trigger
    BEFORE INSERT ON transactions
    FOR EACH ROW EXECUTE PROCEDURE transactions_insert_function();

CREATE INDEX idx_fk_addresses_user_id ON addresses (user_id);

CREATE INDEX idx_fk_accounts_user_id ON accounts (user_id);

CREATE INDEX idx_fk_cards_account_id ON cards (account_id);

CREATE INDEX idx_fk_transactions_account_id ON transactions (account_id);

CREATE INDEX idx_fk_transactions_card_id ON transactions (card_id);

COMMENT ON COLUMN users.gender IS 'currently only 1 or 0 in the sample data';

COMMENT ON COLUMN users.nationality IS 'Alpha-3 code';

COMMENT ON COLUMN users.email IS '320 characters limit on emails according to RFC';

COMMENT ON COLUMN users.mobile_phone_number IS 'E164 designates max 15 numbers';

COMMENT ON COLUMN users.created IS 'GMT timestamp 1microsecond resolution';

COMMENT ON COLUMN users.updated IS 'GMT timestamp 1microsecond resolution';

COMMENT ON COLUMN addresses.zip_code IS 'incl US 9 digit+hyphen code';

COMMENT ON COLUMN accounts.currency IS 'ISO 3-Letter Currency Code';

COMMENT ON COLUMN accounts.iban IS 'Up to 34 characters long, an IBAN is a combination of letters and numbers.';

COMMENT ON COLUMN cards.account_id IS 'redundant imho';

COMMENT ON COLUMN cards.created IS 'GMT timestamp 1microsecond resolution';

COMMENT ON COLUMN cards.updated IS 'GMT timestamp 1microsecond resolution';