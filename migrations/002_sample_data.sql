-- Sample Data for GA Ticketing System
-- This file contains sample data for testing and development

-- Insert sample users
-- Password for all users: password123
INSERT INTO users (employee_id, name, email, department, role, password_hash) VALUES
('EMP001', 'John Doe', 'john.doe@company.com', 'Finance', 'requester', crypt('password123', gen_salt('bf'))),
('EMP002', 'Jane Smith', 'jane.smith@company.com', 'HR', 'requester', crypt('password123', gen_salt('bf'))),
('EMP003', 'Michael Brown', 'michael.brown@company.com', 'IT', 'admin', crypt('password123', gen_salt('bf'))),
('EMP004', 'Sarah Wilson', 'sarah.wilson@company.com', 'Operations', 'approver', crypt('password123', gen_salt('bf'))),
('EMP005', 'David Lee', 'david.lee@company.com', 'Marketing', 'requester', crypt('password123', gen_salt('bf'))),
('EMP006', 'Emily Chen', 'emily.chen@company.com', 'Finance', 'approver', crypt('password123', gen_salt('bf'))),
('EMP007', 'Robert Taylor', 'robert.taylor@company.com', 'Facilities', 'admin', crypt('password123', gen_salt('bf'))),
('EMP008', 'Lisa Anderson', 'lisa.anderson@company.com', 'Sales', 'requester', crypt('password123', gen_salt('bf')))
ON CONFLICT (email) DO NOTHING;

-- Insert sample assets
INSERT INTO assets (asset_code, name, description, category, quantity, available_quantity, location, condition, unit_cost) VALUES
('OFF-DESK-001', 'Office Desk', 'Standard office desk with drawers', 'office_furniture', 20, 18, 'Warehouse A', 'good', 1500000),
('OFF-CHAIR-001', 'Ergonomic Chair', 'Ergonomic office chair with lumbar support', 'office_furniture', 30, 25, 'Warehouse A', 'good', 2000000),
('SUP-PEN-001', 'Ballpoint Pen (Box of 12)', 'Blue ballpoint pens', 'office_supplies', 50, 45, 'Supply Room B', 'good', 50000),
('SUP-PAPER-001', 'A4 Paper (Ream)', 'White A4 paper 80gsm', 'office_supplies', 100, 85, 'Supply Room B', 'good', 40000),
('PAN-COFFEE-001', 'Coffee Beans (1kg)', 'Premium coffee beans', 'pantry_supplies', 10, 8, 'Pantry Storage', 'good', 150000),
('PAN-WATER-001', 'Drinking Water (Gallon)', '19L drinking water', 'pantry_supplies', 30, 25, 'Pantry Storage', 'good', 20000),
('MTG-PROJ-001', 'LCD Projector', 'Full HD projector for meetings', 'meeting_room_equipment', 5, 4, 'Meeting Room 1', 'good', 8000000),
('MTG-BOARD-001', 'Whiteboard', 'Large whiteboard with markers', 'meeting_room_equipment', 8, 8, 'Meeting Rooms', 'good', 500000),
('FAC-AC-001', 'Air Conditioner Unit', '1.5HP split air conditioner', 'facility_equipment', 15, 12, 'Building A', 'needs_maintenance', 5000000),
('FAC-VACUUM-001', 'Vacuum Cleaner', 'Industrial vacuum cleaner', 'cleaning_supplies', 3, 3, 'Janitor Room', 'good', 2500000),
('OFF-MONITOR-001', '24" LED Monitor', 'Full HD LED monitor', 'office_furniture', 40, 35, 'Warehouse A', 'good', 1800000),
('SUP-STAPLER-001', 'Heavy Duty Stapler', 'Stapler with 1000 staples', 'office_supplies', 25, 20, 'Supply Room B', 'good', 75000);

-- Get user IDs for sample tickets (we'll use these in variables)
DO $$
DECLARE
    requester1_id UUID;
    requester2_id UUID;
    requester3_id UUID;
    admin1_id UUID;
    admin2_id UUID;
    approver1_id UUID;
    approver2_id UUID;

    ticket1_id UUID;
    ticket2_id UUID;
    ticket3_id UUID;
    ticket4_id UUID;
    ticket5_id UUID;
    ticket6_id UUID;
    ticket7_id UUID;
    ticket8_id UUID;
BEGIN
    -- Get user IDs
    SELECT id INTO requester1_id FROM users WHERE employee_id = 'EMP001';
    SELECT id INTO requester2_id FROM users WHERE employee_id = 'EMP002';
    SELECT id INTO requester3_id FROM users WHERE employee_id = 'EMP005';
    SELECT id INTO admin1_id FROM users WHERE employee_id = 'EMP003';
    SELECT id INTO admin2_id FROM users WHERE employee_id = 'EMP007';
    SELECT id INTO approver1_id FROM users WHERE employee_id = 'EMP004';
    SELECT id INTO approver2_id FROM users WHERE employee_id = 'EMP006';

    -- Insert sample tickets
    INSERT INTO tickets (ticket_number, title, description, category, priority, status, requester_id, assigned_admin_id, estimated_cost, requires_approval)
    VALUES
        ('', 'Need new office desk', 'Requesting a new office desk for the new employee in Finance department. The current temporary desk is not suitable for long-term use.', 'office_furniture', 'medium', 'pending', requester1_id, NULL, 1500000, true)
        RETURNING id INTO ticket1_id;

    INSERT INTO tickets (ticket_number, title, description, category, priority, status, requester_id, assigned_admin_id, estimated_cost, requires_approval)
    VALUES
        ('', 'Office supplies for new team', 'Need pens, papers, and staplers for 5 new team members joining next week.', 'office_supplies', 'high', 'in_progress', requester2_id, admin1_id, 500000, false)
        RETURNING id INTO ticket2_id;

    INSERT INTO tickets (ticket_number, title, description, category, priority, status, requester_id, assigned_admin_id, estimated_cost, actual_cost, requires_approval)
    VALUES
        ('', 'Pantry coffee restock', 'Coffee beans running low. Need to restock for the month.', 'pantry_supplies', 'low', 'completed', requester3_id, admin2_id, 300000, 300000, false)
        RETURNING id INTO ticket3_id;

    INSERT INTO tickets (ticket_number, title, description, category, priority, status, requester_id, assigned_admin_id, estimated_cost, requires_approval)
    VALUES
        ('', 'AC maintenance in Meeting Room 2', 'Air conditioner in Meeting Room 2 is not cooling properly. Needs immediate maintenance.', 'facility_maintenance', 'high', 'waiting_approval', requester1_id, admin2_id, 2000000, true)
        RETURNING id INTO ticket4_id;

    INSERT INTO tickets (ticket_number, title, description, category, priority, status, requester_id, assigned_admin_id, estimated_cost, requires_approval)
    VALUES
        ('', 'Book Meeting Room 1', 'Need to book Meeting Room 1 for client presentation on Friday, 2PM-4PM.', 'meeting_room', 'medium', 'approved', requester2_id, admin1_id, 0, true)
        RETURNING id INTO ticket5_id;

    INSERT INTO tickets (ticket_number, title, description, category, priority, status, requester_id, assigned_admin_id, estimated_cost, requires_approval)
    VALUES
        ('', 'Install new monitors', 'Need 10 new monitors for the Marketing department workspace upgrade.', 'office_furniture', 'medium', 'pending', requester3_id, NULL, 18000000, true)
        RETURNING id INTO ticket6_id;

    INSERT INTO tickets (ticket_number, title, description, category, priority, status, requester_id, assigned_admin_id, estimated_cost, requires_approval)
    VALUES
        ('', 'Drinking water refill', 'Pantry drinking water running low, need refill by tomorrow.', 'pantry_supplies', 'medium', 'in_progress', requester1_id, admin2_id, 200000, false)
        RETURNING id INTO ticket7_id;

    INSERT INTO tickets (ticket_number, title, description, category, priority, status, requester_id, assigned_admin_id, estimated_cost, requires_approval)
    VALUES
        ('', 'Broken ergonomic chair replacement', 'My ergonomic chair is broken (hydraulic pump not working). Need replacement urgently.', 'office_furniture', 'high', 'rejected', requester2_id, admin1_id, 2000000, true)
        RETURNING id INTO ticket8_id;

    -- Insert comments for some tickets
    INSERT INTO comments (ticket_id, user_id, content) VALUES
        (ticket2_id, admin1_id, 'I have allocated the supplies from our inventory. Will deliver by end of day.'),
        (ticket2_id, requester2_id, 'Thank you! Much appreciated.'),
        (ticket3_id, admin2_id, 'Coffee beans delivered to pantry. Ticket closed.'),
        (ticket3_id, requester3_id, 'Received. Thanks!'),
        (ticket4_id, requester1_id, 'This is urgent as we have important client meetings scheduled this week.'),
        (ticket7_id, admin2_id, 'Water supplier contacted. Delivery scheduled for tomorrow morning.'),
        (ticket8_id, requester2_id, 'This is affecting my productivity. Can we expedite?');

    -- Insert approvals
    INSERT INTO approvals (ticket_id, approver_id, status, comments) VALUES
        (ticket4_id, approver1_id, 'pending', NULL),
        (ticket5_id, approver2_id, 'approved', 'Approved for Friday 2-4 PM. Please confirm with facilities team.'),
        (ticket6_id, approver1_id, 'pending', NULL),
        (ticket8_id, approver2_id, 'rejected', 'Budget constraints this quarter. Please use existing chair from storage room.');

    -- Insert status history
    INSERT INTO status_history (ticket_id, from_status, to_status, changed_by, comments) VALUES
        (ticket2_id, NULL, 'pending', requester2_id, 'Ticket created'),
        (ticket2_id, 'pending', 'in_progress', admin1_id, 'Started processing the request'),
        (ticket3_id, NULL, 'pending', requester3_id, 'Ticket created'),
        (ticket3_id, 'pending', 'in_progress', admin2_id, 'Ordering coffee beans'),
        (ticket3_id, 'in_progress', 'completed', admin2_id, 'Coffee delivered'),
        (ticket4_id, NULL, 'pending', requester1_id, 'Ticket created'),
        (ticket4_id, 'pending', 'waiting_approval', admin2_id, 'Sent for budget approval'),
        (ticket5_id, NULL, 'pending', requester2_id, 'Ticket created'),
        (ticket5_id, 'pending', 'waiting_approval', admin1_id, 'Checking room availability'),
        (ticket5_id, 'waiting_approval', 'approved', approver2_id, 'Meeting room approved'),
        (ticket7_id, NULL, 'pending', requester1_id, 'Ticket created'),
        (ticket7_id, 'pending', 'in_progress', admin2_id, 'Processing water order'),
        (ticket8_id, NULL, 'pending', requester2_id, 'Ticket created'),
        (ticket8_id, 'pending', 'waiting_approval', admin1_id, 'Sent for approval'),
        (ticket8_id, 'waiting_approval', 'rejected', approver2_id, 'Rejected due to budget');

    -- Insert inventory logs
    INSERT INTO inventory_logs (asset_id, change_type, quantity, reason, created_by) VALUES
        ((SELECT id FROM assets WHERE asset_code = 'OFF-DESK-001'), 'remove', 2, 'Allocated for new employees in Finance', admin1_id),
        ((SELECT id FROM assets WHERE asset_code = 'SUP-PEN-001'), 'remove', 5, 'Issued for office supplies request', admin1_id),
        ((SELECT id FROM assets WHERE asset_code = 'SUP-PAPER-001'), 'remove', 15, 'Issued for office supplies request', admin1_id),
        ((SELECT id FROM assets WHERE asset_code = 'PAN-COFFEE-001'), 'remove', 2, 'Pantry restocking', admin2_id),
        ((SELECT id FROM assets WHERE asset_code = 'PAN-WATER-001'), 'remove', 5, 'Pantry water refill', admin2_id),
        ((SELECT id FROM assets WHERE asset_code = 'OFF-CHAIR-001'), 'remove', 5, 'Allocated to Marketing department', admin1_id),
        ((SELECT id FROM assets WHERE asset_code = 'MTG-PROJ-001'), 'remove', 1, 'Maintenance required', admin2_id),
        ((SELECT id FROM assets WHERE asset_code = 'SUP-STAPLER-001'), 'remove', 5, 'Office supplies distribution', admin1_id),
        ((SELECT id FROM assets WHERE asset_code = 'OFF-MONITOR-001'), 'remove', 5, 'IT department upgrade', admin1_id),
        ((SELECT id FROM assets WHERE asset_code = 'FAC-AC-001'), 'adjust', 3, 'Marked for maintenance', admin2_id);

END $$;

-- Update asset available quantities based on inventory logs
UPDATE assets SET available_quantity = 18 WHERE asset_code = 'OFF-DESK-001';
UPDATE assets SET available_quantity = 40 WHERE asset_code = 'SUP-PEN-001';
UPDATE assets SET available_quantity = 85 WHERE asset_code = 'SUP-PAPER-001';
UPDATE assets SET available_quantity = 8 WHERE asset_code = 'PAN-COFFEE-001';
UPDATE assets SET available_quantity = 25 WHERE asset_code = 'PAN-WATER-001';
UPDATE assets SET available_quantity = 25 WHERE asset_code = 'OFF-CHAIR-001';
UPDATE assets SET available_quantity = 4 WHERE asset_code = 'MTG-PROJ-001';
UPDATE assets SET available_quantity = 20 WHERE asset_code = 'SUP-STAPLER-001';
UPDATE assets SET available_quantity = 35 WHERE asset_code = 'OFF-MONITOR-001';
UPDATE assets SET available_quantity = 12 WHERE asset_code = 'FAC-AC-001';
