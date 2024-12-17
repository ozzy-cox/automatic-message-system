INSERT INTO messages (content, to_, sending_status, sent_at)
SELECT 
    'Message content ' || gs.id AS content,
    LPAD((1000000000 + FLOOR(random() * 8999999999))::TEXT, 10, '0') AS to_,
    FALSE AS sending_status,
    CASE 
        WHEN random() > 0.5 THEN NOW() - (random() * INTERVAL '30 days') 
        ELSE NULL 
    END AS sent_at
FROM generate_series(1, 1000) AS gs(id);
