# Tasks: GA Ticketing System

**Input**: Design documents from `/specs/001-ga-ticketing/`
**Prerequisites**: plan.md (completed), spec.md (completed), research.md, data-model.md (completed), contracts/api.yaml (completed)

**Tests**: Included - TDD approach with unit, integration, and contract tests required by constitution compliance

**Organization**: Tasks grouped by user story to enable independent implementation and testing of each story

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Single backend service**: `src/`, `tests/` at repository root
- **Clean architecture**: Domain, Application, Infrastructure, Interface layers
- Paths shown follow plan.md structure

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [ ] T001 Create project structure per clean architecture plan
- [ ] T002 Initialize Go 1.22+ module with required dependencies (chi, pgx, testify, etc.)
- [ ] T003 [P] Configure golangci-lint and gofmt tools
- [ ] T004 [P] Setup automated quality gates (coverage, complexity, static analysis)
- [ ] T005 [P] Configure Makefile for build, test, and development tasks
- [ ] T006 Create .env.example with all required configuration variables
- [ ] T007 Setup Docker Compose for local development (PostgreSQL, Redis)
- [ ] T008 [P] Configure gitignore and pre-commit hooks

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Database Setup
- [ ] T009 Setup PostgreSQL connection pool and configuration in src/infrastructure/database/
- [ ] T010 Create database migration framework using migrate library
- [ ] T011 Implement all table schemas (users, tickets, assets, comments, approvals, status_history, inventory_logs) in migrations/
- [ ] T012 Create recommended database indexes for performance
- [ ] T013 Setup Redis connection for caching in src/infrastructure/cache/

### Authentication & Authorization
- [ ] T014 [P] Implement JWT token generation and validation in src/infrastructure/auth/
- [ ] T015 [P] Create OAuth2/OIDC integration framework
- [ ] T016 Implement role-based access control (RBAC) middleware
- [ ] T017 Create user authentication service with password hashing

### Architecture Foundation
- [ ] T018 Setup clean architecture layers (Domain, Application, Infrastructure, Interface)
- [ ] T019 Configure dependency injection container using wire or similar
- [ ] T020 Implement error handling and structured logging infrastructure
- [ ] T021 Create base HTTP server configuration with Chi router
- [ ] T022 Setup environment configuration management in src/config/
- [ ] T023 Create base repository interfaces in src/domain/repositories/
- [ ] T024 [P] Setup CORS, rate limiting, and security middleware

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Employee Service Request Submission (Priority: P1) 🎯 MVP

**Goal**: Employees can submit service requests and track their status

**Independent Test**: An employee can create a service request, receive confirmation, and track the request through completion without needing other system features.

### Tests for User Story 1 (MANDATORY - TDD Constitution Requirement) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T025 [P] [US1] Unit tests for Ticket entity in tests/unit/domain/test_ticket.go
- [ ] T026 [P] [US1] Unit tests for User entity in tests/unit/domain/test_user.go
- [ ] T027 [P] [US1] Unit tests for TicketService domain logic in tests/unit/domain/services/test_ticket_service.go
- [ ] T028 [P] [US1] Contract test for POST /v1/tickets in tests/contract/test_ticket_creation.go
- [ ] T029 [P] [US1] Contract test for GET /v1/tickets in tests/contract/test_ticket_list.go
- [ ] T030 [P] [US1] Integration test for ticket submission workflow in tests/integration/test_ticket_workflow.go
- [ ] T031 [P] [US1] Performance test for ticket creation under load in tests/performance/test_ticket_creation.go

### Implementation for User Story 1 (Clean Architecture & DDD)

**Domain Layer** (innermost - no external dependencies):
- [ ] T032 [P] [US1] Create Ticket aggregate root in src/domain/entities/ticket.go
- [ ] T033 [P] [US1] Create User aggregate root in src/domain/entities/user.go
- [ ] T034 [P] [US1] Create Comment entity in src/domain/entities/comment.go
- [ ] T035 [P] [US1] Create StatusHistory entity in src/domain/entities/status_history.go
- [ ] T036 [P] [US1] Create value objects (Money, Email) in src/domain/valueobjects/
- [ ] T037 [US1] Implement TicketService domain logic in src/domain/services/ticket_service.go (depends on T032-T036)
- [ ] T038 [P] [US1] Define repository interfaces in src/domain/repositories/ticket_repository.go
- [ ] T039 [P] [US1] Define repository interfaces in src/domain/repositories/user_repository.go

**Application Layer** (use cases orchestration):
- [ ] T040 [US1] Implement CreateTicket use case in src/application/usecases/create_ticket.go
- [ ] T041 [US1] Implement GetTickets use case in src/application/usecases/get_tickets.go
- [ ] T042 [US1] Implement GetTicket use case in src/application/usecases/get_ticket.go
- [ ] T043 [P] [US1] Create input/output DTOs in src/application/dto/ticket_dto.go

**Infrastructure Layer** (external concerns):
- [ ] T044 [US1] Implement TicketRepository in src/infrastructure/database/repositories/ticket_repository.go
- [ ] T045 [US1] Implement UserRepository in src/infrastructure/database/repositories/user_repository.go
- [ ] T046 [P] [US1] Setup password hashing service in src/infrastructure/auth/password.go

**Interface Layer** (outermost):
- [ ] T047 [US1] Implement ticket HTTP handlers in src/interface/http/handlers/ticket_handler.go
- [ ] T048 [US1] Create ticket validation middleware in src/interface/http/middleware/validation.go
- [ ] T049 [US1] Add ticket routes to Chi router in src/interface/http/router/routes.go
- [ ] T050 [US1] Create authentication middleware in src/interface/http/middleware/auth.go
- [ ] T051 [US1] Add structured logging with correlation IDs for ticket operations

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Admin Ticket Management and Processing (Priority: P1)

**Goal**: GA administrators can view, assign, and process all tickets in the system

**Independent Test**: An admin can view the complete ticket queue, assign tickets, update inventory, and process requests to completion without requiring employee interaction.

### Tests for User Story 2 (MANDATORY - TDD Constitution Requirement) ⚠️

- [ ] T052 [P] [US2] Unit tests for Asset entity in tests/unit/domain/test_asset.go
- [ ] T053 [P] [US2] Unit tests for AssetService domain logic in tests/unit/domain/services/test_asset_service.go
- [ ] T054 [P] [US2] Contract test for POST /v1/tickets/{id}/assign in tests/contract/test_ticket_assignment.go
- [ ] T055 [P] [US2] Contract test for PUT /v1/tickets/{id} in tests/contract/test_ticket_update.go
- [ ] T056 [P] [US2] Contract test for GET /v1/assets in tests/contract/test_asset_management.go
- [ ] T057 [P] [US2] Integration test for admin ticket processing workflow in tests/integration/test_admin_workflow.go

### Implementation for User Story 2

**Domain Layer**:
- [ ] T058 [P] [US2] Create Asset aggregate root in src/domain/entities/asset.go
- [ ] T059 [P] [US2] Create InventoryLog entity in src/domain/entities/inventory_log.go
- [ ] T060 [US2] Implement AssetService domain logic in src/domain/services/asset_service.go
- [ ] T061 [P] [US2] Extend TicketService with assignment logic
- [ ] T062 [P] [US2] Define asset repository interfaces in src/domain/repositories/asset_repository.go

**Application Layer**:
- [ ] T063 [US2] Implement AssignTicket use case in src/application/usecases/assign_ticket.go
- [ ] T064 [US2] Implement UpdateTicket use case in src/application/usecases/update_ticket.go
- [ ] T065 [US2] Implement CreateAsset use case in src/application/usecases/create_asset.go
- [ ] T066 [US2] Implement GetAssets use case in src/application/usecases/get_assets.go
- [ ] T067 [P] [US2] Create asset DTOs in src/application/dto/asset_dto.go

**Infrastructure Layer**:
- [ ] T068 [US2] Implement AssetRepository in src/infrastructure/database/repositories/asset_repository.go
- [ ] T069 [US2] Add admin role checks to authentication middleware

**Interface Layer**:
- [ ] T070 [US2] Extend ticket handlers with assignment capabilities
- [ ] T071 [US2] Implement asset HTTP handlers in src/interface/http/handlers/asset_handler.go
- [ ] T072 [US2] Add asset routes to Chi router
- [ ] T073 [US2] Implement admin authorization middleware

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Approval Workflow Management (Priority: P1)

**Goal**: Approvers can review, approve, and reject requests requiring budget approval

**Independent Test**: An approver can review all pending approval requests, make decisions with rationale, and track approval history without needing to interact with other user roles.

### Tests for User Story 3 (MANDATORY - TDD Constitution Requirement) ⚠️

- [ ] T074 [P] [US3] Unit tests for Approval entity in tests/unit/domain/test_approval.go
- [ ] T075 [P] [US3] Unit tests for approval workflow logic in tests/unit/domain/services/test_approval_service.go
- [ ] T076 [P] [US3] Contract test for POST /v1/tickets/{id}/approve in tests/contract/test_approval.go
- [ ] T077 [P] [US3] Contract test for POST /v1/tickets/{id}/reject in tests/contract/test_rejection.go
- [ ] T078 [P] [US3] Integration test for approval workflow with optimistic locking in tests/integration/test_approval_workflow.go
- [ ] T079 [P] [US3] Performance test for concurrent approval scenarios in tests/performance/test_concurrent_approval.go

### Implementation for User Story 3

**Domain Layer**:
- [ ] T080 [P] [US3] Create Approval entity in src/domain/entities/approval.go
- [ ] T081 [US3] Implement ApprovalService domain logic with optimistic locking in src/domain/services/approval_service.go
- [ ] T082 [P] [US3] Extend TicketService with approval rules validation
- [ ] T083 [P] [US3] Define approval repository interfaces in src/domain/repositories/approval_repository.go

**Application Layer**:
- [ ] T084 [US3] Implement ApproveTicket use case in src/application/usecases/approve_ticket.go
- [ ] T085 [US3] Implement RejectTicket use case in src/application/usecases/reject_ticket.go
- [ ] T086 [P] [US3] Create approval DTOs in src/application/dto/approval_dto.go

**Infrastructure Layer**:
- [ ] T087 [US3] Implement ApprovalRepository with optimistic locking in src/infrastructure/database/repositories/approval_repository.go

**Interface Layer**:
- [ ] T088 [US3] Implement approval HTTP handlers in src/interface/http/handlers/approval_handler.go
- [ ] T089 [US3] Add approval routes to Chi router
- [ ] T090 [US3] Implement approver role authorization middleware

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: User Story 4 - Ticket Communication and History (Priority: P2)

**Goal**: All users involved in a ticket can add comments and track complete history

**Independent Test**: Users can add comments to tickets they have access to and view complete comment history without needing other features to function.

### Tests for User Story 4

- [ ] T091 [P] [US4] Unit tests for comment creation and retrieval in tests/unit/domain/test_comment_service.go
- [ ] T092 [P] [US4] Contract test for POST /v1/tickets/{id}/comments in tests/contract/test_comment_creation.go
- [ ] T093 [P] [US4] Contract test for GET /v1/tickets/{id}/comments in tests/contract/test_comment_list.go
- [ ] T094 [P] [US4] Integration test for ticket communication workflow in tests/integration/test_ticket_communication.go

### Implementation for User Story 4

**Domain Layer**:
- [ ] T095 [P] [US4] Extend Comment entity with access control logic
- [ ] T096 [US4] Implement CommentService domain logic in src/domain/services/comment_service.go

**Application Layer**:
- [ ] T097 [US4] Implement AddComment use case in src/application/usecases/add_comment.go
- [ ] T098 [US4] Implement GetComments use case in src/application/usecases/get_comments.go
- [ ] T099 [P] [US4] Create comment DTOs in src/application/dto/comment_dto.go

**Infrastructure Layer**:
- [ ] T100 [US4] Implement CommentRepository in src/infrastructure/database/repositories/comment_repository.go

**Interface Layer**:
- [ ] T101 [US4] Extend ticket handlers with comment endpoints
- [ ] T102 [US4] Add comment routes to Chi router

---

## Phase 7: User Story 5 - Asset and Inventory Management (Priority: P2)

**Goal**: Administrators can maintain accurate asset records and manage inventory

**Independent Test**: An admin can perform complete asset lifecycle management (view, add, update quantities, track conditions) without affecting ticket processing functionality.

### Tests for User Story 5

- [ ] T103 [P] [US5] Unit tests for inventory management logic in tests/unit/domain/test_inventory_service.go
- [ ] T104 [P] [US5] Contract test for POST /v1/assets/{id}/inventory in tests/contract/test_inventory_update.go
- [ ] T105 [P] [US5] Contract test for PUT /v1/assets/{id} in tests/contract/test_asset_update.go
- [ ] T106 [P] [US5] Integration test for asset lifecycle management in tests/integration/test_asset_lifecycle.go

### Implementation for User Story 5

**Domain Layer**:
- [ ] T107 [P] [US5] Extend Asset entity with maintenance scheduling logic
- [ ] T108 [US5] Implement InventoryService domain logic in src/domain/services/inventory_service.go
- [ ] T109 [US5] Implement MaintenanceService domain logic in src/domain/services/maintenance_service.go

**Application Layer**:
- [ ] T110 [US5] Implement UpdateInventory use case in src/application/usecases/update_inventory.go
- [ ] T111 [US5] Implement UpdateAsset use case in src/application/usecases/update_asset.go
- [ ] T112 [US5] Implement ScheduleMaintenance use case in src/application/usecases/schedule_maintenance.go

**Infrastructure Layer**:
- [ ] T113 [US5] Extend AssetRepository with inventory management methods
- [ ] T114 [US5] Add maintenance scheduling queries and indexes

**Interface Layer**:
- [ ] T115 [US5] Extend asset handlers with inventory management endpoints
- [ ] T116 [US5] Add inventory update routes to Chi router

---

## Phase 8: Quality Assurance & Constitution Compliance

**Purpose**: Ensure all constitutional requirements are met before deployment

### Code Quality & Testing
- [ ] T117 [P] Run full test suite and verify 90%+ coverage maintained
- [ ] T118 [P] Execute golangci-lint and ensure zero violations
- [ ] T119 [P] Verify cyclomatic complexity ≤ 10 for all new code
- [ ] T120 [P] Run performance tests and validate ≤200ms API response times
- [ ] T121 [P] Execute load testing for 500+ concurrent users
- [ ] T122 [P] Validate clean architecture layer dependencies (no violations)
- [ ] T123 [P] Run integration tests with PostgreSQL and Redis
- [ ] T124 [P] Execute contract tests against API specification

### Security & Performance
- [ ] T125 [P] Security hardening review (SQL injection, XSS, CSRF protection)
- [ ] T126 [P] Validate JWT token security and OAuth2 integration
- [ ] T127 [P] Review RBAC implementation and test access controls
- [ ] T128 [P] Test optimistic locking for concurrent approval scenarios
- [ ] T129 [P] Database performance testing with recommended indexes
- [ ] T130 [P] Redis caching performance validation

### Documentation & Deployment
- [ ] T131 Update API documentation with any implementation details
- [ ] T132 Validate quickstart.md setup instructions work correctly
- [ ] T133 Create database migration rollback procedures
- [ ] T134 Setup monitoring and observability (health checks, metrics)
- [ ] T135 Create deployment configuration (Docker, environment variables)
- [ ] T136 Code cleanup and refactoring (if needed for quality gates)

### Constitution Compliance Verification
- [ ] T137 Final TDD compliance verification (all tests written first)
- [ ] T138 DDD architecture compliance review
- [ ] T139 Clean architecture layer dependency validation
- [ ] T140 Performance requirement validation (response times, concurrent users)
- [ ] T141 UX consistency review across all implemented features
- [ ] T142 Final constitution compliance sign-off

**CONSTITUTION COMPLIANCE GATE**: All TDD, DDD, Clean Architecture, Quality, and Performance requirements must be satisfied before deployment.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phases 3-7)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2)
- **Polish (Phase 8)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - May integrate with US1 but should be independently testable
- **User Story 3 (P1)**: Can start after Foundational (Phase 2) - Depends on US1 ticket structure
- **User Story 4 (P2)**: Depends on US1 and US2 (comments require tickets and users)
- **User Story 5 (P2)**: Can start after Foundational (Phase 2) - Independent of ticket workflows

### Within Each User Story

- Tests MUST be written and FAIL before implementation (TDD requirement)
- Domain entities before services
- Repository interfaces before implementations
- Application use cases before HTTP handlers
- Core implementation before integration

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- All tests for a user story marked [P] can run in parallel
- Domain entities within a story marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together (TDD approach):
Task: "Contract test for POST /v1/tickets in tests/contract/test_ticket_creation.go"
Task: "Contract test for GET /v1/tickets in tests/contract/test_ticket_list.go"
Task: "Integration test for ticket submission workflow in tests/integration/test_ticket_workflow.go"
Task: "Performance test for ticket creation under load in tests/performance/test_ticket_creation.go"

# Launch all domain entities for User Story 1 together:
Task: "Create Ticket aggregate root in src/domain/entities/ticket.go"
Task: "Create User aggregate root in src/domain/entities/user.go"
Task: "Create Comment entity in src/domain/entities/comment.go"
Task: "Create StatusHistory entity in src/domain/entities/status_history.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Add User Stories 4 & 5 → Test independently → Deploy/Demo
6. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (Core ticketing)
   - Developer B: User Story 2 (Admin management)
   - Developer C: User Story 3 (Approval workflow)
3. Stories complete and integrate independently
4. Developer D picks up User Stories 4 & 5 (P2 features)

---

## File-Based Coordination Rules

### Database Schema Files
- All migration files in `migrations/` directory
- Naming convention: `001_initial_schema.up.sql`, `001_initial_schema.down.sql`
- No duplicate table definitions across files

### Domain Entity Files
- One entity per file in `src/domain/entities/`
- File name matches entity name (e.g., `ticket.go`, `user.go`)
- No circular dependencies between entities

### Repository Implementation Files
- One repository per file in `src/infrastructure/database/repositories/`
- Interface in `src/domain/repositories/`, implementation in infrastructure
- Use dependency injection for database connections

### HTTP Handler Files
- Group handlers by resource (e.g., `ticket_handler.go`, `asset_handler.go`)
- Use consistent error response format across all handlers
- Implement proper input validation and status codes

### Test Files
- Unit tests mirror source structure in `tests/unit/`
- Integration tests in `tests/integration/`
- Contract tests in `tests/contract/`
- Performance tests in `tests/performance/`

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing (TDD requirement)
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- Follow clean architecture dependency rules strictly
- Maintain 90%+ test coverage throughout development
- All API responses must be ≤200ms for performance requirements