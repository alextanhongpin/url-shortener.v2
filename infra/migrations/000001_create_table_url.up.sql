CREATE TABLE IF NOT EXISTS url (
	code TEXT PRIMARY KEY,
	long_url TEXT,
	expire_at TIMESTAMP WITH TIME ZONE
);
