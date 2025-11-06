-- GA Ticketing System Schema Rollback
-- Migration: 001_initial_schema
-- Description: Rollback initial database schema for GA ticketing system

-- Drop triggers
DROP TRIGGER IF EXISTS generate_ticket_number_trigger ON tickets;
DROP TRIGGER IF EXISTS validate_ticket_status_transition_trigger ON tickets;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_tickets_updated_at ON tickets;
DROP TRIGGER IF EXISTS update_assets_updated_at ON assets;

-- Drop functions
DROP FUNCTION IF EXISTS generate_ticket_number();
DROP FUNCTION IF EXISTS validate_ticket_status_transition();
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop constraints
ALTER TABLE assets DROP CONSTRAINT IF EXISTS check_available_quantity;

-- Drop tables in reverse order of creation
DROP TABLE IF EXISTS inventory_logs;
DROP TABLE IF EXISTS status_history;
DROP TABLE IF EXISTS approvals;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS assets;
DROP TABLE IF EXISTS tickets;
DROP TABLE IF EXISTS users;

-- Drop custom types
DROP TYPE IF EXISTS change_type;
DROP TYPE IF EXISTS approval_status;
DROP TYPE IF EXISTS asset_condition;
DROP TYPE IF EXISTS asset_category;
DROP TYPE IF EXISTS ticket_status;
DROP TYPE IF EXISTS ticket_priority;
DROP TYPE IF EXISTS ticket_category;
DROP TYPE IF EXISTS user_role;

-- Drop extensions
DROP EXTENSION IF EXISTS "pgcrypto";
DROP EXTENSION IF EXISTS "uuid-ossp";