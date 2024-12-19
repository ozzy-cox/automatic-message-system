CREATE TABLE IF NOT EXISTS messages (
	id SERIAL PRIMARY KEY,
	content VARCHAR(1023),
	to_ VARCHAR(1023),
	is_sent  BOOLEAN DEFAULT false,
	sent_at TIMESTAMP NULL DEFAULT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

TRUNCATE TABLE messages;

INSERT INTO messages (content, to_, is_sent, sent_at)
SELECT 
    'Message content ' || gs.id AS content,
    LPAD((1000000000 + FLOOR(random() * 8999999999))::TEXT, 10, '0') AS to_,
    FALSE AS is_sent,
    NULL AS sent_at
FROM generate_series(1, 100) AS gs(id);

