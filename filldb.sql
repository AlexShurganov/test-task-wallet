INSERT INTO wallets (id, balance)
SELECT 
    gen_random_uuid(), 
    round((random() * 10000)::numeric, 2)
FROM generate_series(1, 1000);
