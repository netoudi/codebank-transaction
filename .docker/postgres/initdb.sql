CREATE TABLE IF NOT EXISTS credit_cards
(
	id               uuid      NOT NULL,
	name             VARCHAR   NOT NULL,
	number           VARCHAR   NOT NULL,
	expiration_month VARCHAR   NOT NULL,
	expiration_year  VARCHAR,
	cvv              VARCHAR   NOT NULL,
	balance          float     not null,
	balance_Limit    float     not null,
	created_at       timestamp not null,
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS transactions
(
	id             uuid      NOT NULL,
	credit_card_id uuid      NOT NULL references credit_cards (id),
	amount         float     NOT NULL,
	status         VARCHAR   NOT NULL,
	description    VARCHAR,
	store          VARCHAR   NOT NULL,
	created_at     timestamp not null,
	PRIMARY KEY (id)
);
