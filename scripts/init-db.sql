-- GA Ticketing System Database Initialization
-- This script runs when PostgreSQL container starts for the first time

-- Create test database
CREATE DATABASE ga_ticketing_test;

-- Grant permissions to user for both databases
GRANT ALL PRIVILEGES ON DATABASE ga_ticketing TO ga_user;
GRANT ALL PRIVILEGES ON DATABASE ga_ticketing_test TO ga_user;

-- Connect to main database and set up extensions
\c ga_ticketing;

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Connect to test database and set up extensions
\c ga_ticketing_test;

-- Create extensions for test database
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Switch back to main database
\c ga_ticketing;