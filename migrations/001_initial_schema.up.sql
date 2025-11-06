-- GA Ticketing System Initial Schema
-- Migration: 001_initial_schema
-- Created: 2025-11-06
-- Description: Create initial database schema for GA ticketing system

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create custom types
CREATE TYPE user_role AS ENUM ('requester', 'approver', 'admin');
CREATE TYPE ticket_category AS ENUM ('office_supplies', 'facility_maintenance', 'pantry_supplies', 'meeting_room', 'office_furniture', 'general_service');
CREATE TYPE ticket_priority AS ENUM ('low', 'medium', 'high');
CREATE TYPE ticket_status AS ENUM ('pending', 'waiting_approval', 'approved', 'rejected', 'in_progress', 'completed', 'closed');
CREATE TYPE asset_category AS ENUM ('office_furniture', 'office_supplies', 'pantry_supplies', 'facility_equipment', 'meeting_room_equipment', 'cleaning_supplies');
CREATE TYPE asset_condition AS ENUM ('good', 'needs_maintenance', 'broken');
CREATE TYPE approval_status AS ENUM ('pending', 'approved', 'rejected');
CREATE TYPE change_type AS ENUM ('add', 'remove', 'adjust');

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    department VARCHAR(255),
    role user_role NOT NULL DEFAULT 'requester',
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tickets table
CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_number VARCHAR(20) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    category ticket_category NOT NULL,
    priority ticket_priority NOT NULL DEFAULT 'medium',
    status ticket_status NOT NULL DEFAULT 'pending',
    requester_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    assigned_admin_id UUID REFERENCES users(id) ON DELETE SET NULL,
    estimated_cost BIGINT NOT NULL DEFAULT 0 CHECK (estimated_cost >= 0),
    actual_cost BIGINT CHECK (actual_cost >= 0),
    requires_approval BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE
);

-- Assets table
CREATE TABLE assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    asset_code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category asset_category NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1 CHECK (quantity >= 0),
    available_quantity INTEGER NOT NULL DEFAULT 1 CHECK (available_quantity >= 0),
    location VARCHAR(255),
    condition asset_condition NOT NULL DEFAULT 'good',
    unit_cost BIGINT NOT NULL DEFAULT 0 CHECK (unit_cost >= 0),
    last_maintenance_at TIMESTAMP WITH TIME ZONE,
    next_maintenance_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Comments table
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Approvals table
CREATE TABLE approvals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    approver_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status approval_status NOT NULL DEFAULT 'pending',
    comments TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Status history table
CREATE TABLE status_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    from_status ticket_status,
    to_status ticket_status NOT NULL,
    changed_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    comments TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Inventory logs table
CREATE TABLE inventory_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    change_type change_type NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    reason TEXT NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
-- User indexes
CREATE INDEX idx_users_employee_id ON users(employee_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_active ON users(is_active);

-- Ticket indexes
CREATE INDEX idx_tickets_requester_id ON tickets(requester_id);
CREATE INDEX idx_tickets_assigned_admin_id ON tickets(assigned_admin_id);
CREATE INDEX idx_tickets_status ON tickets(status);
CREATE INDEX idx_tickets_category ON tickets(category);
CREATE INDEX idx_tickets_priority ON tickets(priority);
CREATE INDEX idx_tickets_created_at ON tickets(created_at);
CREATE INDEX idx_tickets_ticket_number ON tickets(ticket_number);
CREATE INDEX idx_tickets_requires_approval ON tickets(requires_approval);

-- Asset indexes
CREATE INDEX idx_assets_category ON assets(category);
CREATE INDEX idx_assets_condition ON assets(condition);
CREATE INDEX idx_assets_location ON assets(location);
CREATE INDEX idx_assets_asset_code ON assets(asset_code);
CREATE INDEX idx_assets_available_quantity ON assets(available_quantity);

-- Comment indexes
CREATE INDEX idx_comments_ticket_id ON comments(ticket_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_created_at ON comments(created_at);

-- Approval indexes
CREATE INDEX idx_approvals_ticket_id ON approvals(ticket_id);
CREATE INDEX idx_approvals_approver_id ON approvals(approver_id);
CREATE INDEX idx_approvals_status ON approvals(status);

-- Status history indexes
CREATE INDEX idx_status_history_ticket_id ON status_history(ticket_id);
CREATE INDEX idx_status_history_changed_by ON status_history(changed_by);
CREATE INDEX idx_status_history_created_at ON status_history(created_at);

-- Inventory log indexes
CREATE INDEX idx_inventory_logs_asset_id ON inventory_logs(asset_id);
CREATE INDEX idx_inventory_logs_created_by ON inventory_logs(created_by);
CREATE INDEX idx_inventory_logs_created_at ON inventory_logs(created_at);

-- Create trigger function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tickets_updated_at BEFORE UPDATE ON tickets FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_assets_updated_at BEFORE UPDATE ON assets FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create function to generate ticket numbers
CREATE OR REPLACE FUNCTION generate_ticket_number()
RETURNS TRIGGER AS $$
DECLARE
    year_part TEXT;
    sequence_num INTEGER;
    ticket_num TEXT;
BEGIN
    year_part := EXTRACT(year FROM CURRENT_TIMESTAMP)::TEXT;

    -- Get next sequence number for the year
    SELECT COALESCE(MAX(CAST(SUBSTRING(ticket_number, 7, 4) AS INTEGER)), 0) + 1
    INTO sequence_num
    FROM tickets
    WHERE ticket_number LIKE 'GA-' || year_part || '-%';

    -- Format as GA-YYYY-NNNN with leading zeros
    ticket_num := 'GA-' || year_part || '-' || LPAD(sequence_num::TEXT, 4, '0');

    NEW.ticket_number := ticket_num;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for automatic ticket number generation
CREATE TRIGGER generate_ticket_number_trigger
    BEFORE INSERT ON tickets
    FOR EACH ROW
    WHEN (NEW.ticket_number IS NULL OR NEW.ticket_number = '')
    EXECUTE FUNCTION generate_ticket_number();

-- Create constraint to ensure available_quantity <= quantity
ALTER TABLE assets ADD CONSTRAINT check_available_quantity
    CHECK (available_quantity <= quantity);

-- Create constraint for ticket status transitions
CREATE OR REPLACE FUNCTION validate_ticket_status_transition()
RETURNS TRIGGER AS $$
BEGIN
    -- Allow any initial status when creating
    IF TG_OP = 'INSERT' THEN
        RETURN NEW;
    END IF;

    -- Define valid status transitions
    IF OLD.status = 'pending' THEN
        IF NEW.status NOT IN ('waiting_approval', 'in_progress', 'closed') THEN
            RAISE EXCEPTION 'Invalid status transition from % to %', OLD.status, NEW.status;
        END IF;
    ELSIF OLD.status = 'waiting_approval' THEN
        IF NEW.status NOT IN ('approved', 'rejected', 'closed') THEN
            RAISE EXCEPTION 'Invalid status transition from % to %', OLD.status, NEW.status;
        END IF;
    ELSIF OLD.status = 'approved' THEN
        IF NEW.status NOT IN ('in_progress', 'closed') THEN
            RAISE EXCEPTION 'Invalid status transition from % to %', OLD.status, NEW.status;
        END IF;
    ELSIF OLD.status = 'rejected' THEN
        IF NEW.status NOT IN ('closed') THEN
            RAISE EXCEPTION 'Invalid status transition from % to %', OLD.status, NEW.status;
        END IF;
    ELSIF OLD.status = 'in_progress' THEN
        IF NEW.status NOT IN ('completed', 'closed') THEN
            RAISE EXCEPTION 'Invalid status transition from % to %', OLD.status, NEW.status;
        END IF;
    ELSIF OLD.status = 'completed' THEN
        IF NEW.status NOT IN ('closed') THEN
            RAISE EXCEPTION 'Invalid status transition from % to %', OLD.status, NEW.status;
        END IF;
    ELSIF OLD.status = 'closed' THEN
        RAISE EXCEPTION 'Cannot change status from closed';
    END IF;

    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for status transition validation
CREATE TRIGGER validate_ticket_status_transition_trigger
    BEFORE UPDATE OF status ON tickets
    FOR EACH ROW
    EXECUTE FUNCTION validate_ticket_status_transition();

-- Insert default admin user (password: admin123)
INSERT INTO users (employee_id, name, email, department, role, password_hash)
VALUES (
    'ADMIN001',
    'System Administrator',
    'admin@company.com',
    'IT',
    'admin',
    crypt('admin123', gen_salt('bf'))
) ON CONFLICT (email) DO NOTHING;