USE transaction_db;

SET FOREIGN_KEY_CHECKS=0;
DROP TABLE IF EXISTS account_transaction;
DROP TABLE IF EXISTS transaction_type;
SET FOREIGN_KEY_CHECKS=1;

CREATE TABLE account_transaction (
    id_transaction VARCHAR(255) NOT NULL PRIMARY KEY,
    sender_id VARCHAR(255) NOT NULL,
    recipient_id VARCHAR(255) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    t_date DATE NOT NUll,
    fk_t_type INT UNSIGNED NOT NULL
);

CREATE TABLE transaction_type (
    id_transaction_type INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    t_type VARCHAR(255) NOT NULL
);

ALTER TABLE account_transaction
ADD CONSTRAINT fkc_transaction_type_account_transaction
FOREIGN KEY (fk_t_type)
REFERENCES transaction_type(id_transaction_type)
ON UPDATE CASCADE
ON DELETE CASCADE;
