CREATE TABLE base_mints
(
    transaction_hash varchar(255) NOT NULL,
    evt_index int4 NOT NULL,
    operator varchar(100),
    from_address varchar(100) NOT NULL,
    to_address varchar(100) NOT NULL,
    token_id varchar(255) NOT NULL,
    value varchar(100),
    symbol varchar(100),
    token_standard varchar(100),
    name varchar(255),
    contract_address varchar(100),
    mint_type varchar(100),
    number_of_nft_minted int4,
    block_number int8 NOT NULL,
    block_time timestamp NOT NULL,
    block_date date NOT NULL,
    PRIMARY KEY (transaction_hash,evt_index,token_id)
);