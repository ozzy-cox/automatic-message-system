INSERT INTO messages (content, to_, is_sent, sent_at)
SELECT 
    'Message content ' || gs.id AS content,
    LPAD((1000000000 + FLOOR(random() * 8999999999))::TEXT, 10, '0') AS to_,
    FALSE AS is_sent,
    NULL AS sent_at
FROM generate_series(1, 10) AS gs(id);
