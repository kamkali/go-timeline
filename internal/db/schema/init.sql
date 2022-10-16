-- Simulate IF NOT EXIST on CREATE DATABASE
SELECT 'CREATE DATABASE timeline'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'timeline')\gexec