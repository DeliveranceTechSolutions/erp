-- Version: 1.1
-- Description: Create table users
CREATE TABLE users (
	user_id       UUID,
	name          TEXT,
	email         TEXT UNIQUE,
	roles         TEXT[],
	password_hash TEXT,
	date_created  TIMESTAMP,
	date_updated  TIMESTAMP,
	auctions_won  UUID[],

	PRIMARY KEY (user_id)
);

-- Version: 1.2
-- Description: Create table products
CREATE TABLE products (
	product_id   UUID,
	name         TEXT,
	quantity     INT,
	owner_id      UUID,
	date_created TIMESTAMP,
	date_updated TIMESTAMP,
	min_retainer INT,
	starter_bid	 INT,
	bids		 UUID[],

	PRIMARY KEY (product_id),
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Version: 1.3
-- Description: Create table bids
CREATE TABLE bids (
	bid_id       UUID,
	user_id      UUID,
	product_id   UUID,
	quantity     INT,
	bid_price    INT,
	date_created TIMESTAMP,

	PRIMARY KEY (sale_id),
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
	FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE
);

-- Version: 1.0
-- Description: Create vendor group that can be used to reference bids, users watching, and history of products sold.
CREATE TABLE vendors (
	vendor_id	UUID,
	watchers	UUID[],
	products	UUID[],
	history		UUID[],

	PRIMARY KEY (vendor_id)
);