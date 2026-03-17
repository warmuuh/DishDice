# User Roles and Registration Approval System - Implementation Notes

## Implementation Summary

Successfully implemented a complete user roles and registration approval system for DishDice. The system includes:

1. **Database Schema Changes**
   - Added `role` column (user/admin) with CHECK constraint
   - Added `status` column (pending/approved/rejected) with CHECK constraint
   - Added indexes for efficient admin queries
   - Grandfathered existing users to 'approved' status

2. **Backend Implementation**
   - Updated User model with role and status fields
   - Modified user repository to auto-promote first user to admin
   - Enhanced JWT tokens to include role and status claims
   - Added status checks in login flow
   - Created AdminService for user management operations
   - Added RequireAdmin and RequireApproved middleware
   - Created admin endpoints for user management

3. **Frontend Implementation**
   - Updated type definitions with RegisterResponse and AdminUser
   - Modified AuthContext to handle pending registrations
   - Added isAdmin computed property to context
   - Created AdminRoute guard component
   - Updated Header with admin link for admin users
   - Modified Register page to redirect pending users to waiting page
   - Enhanced Login page with specific error messages
   - Created WaitingApproval page for pending users
   - Created AdminPanel page with user management UI

## Key Features

### First User Bootstrap
The first user to register automatically becomes an admin with approved status, solving the "chicken and egg" problem of needing an admin to approve admins.

### User Status Flow
- **New Users**: Register → Pending status → No token → Waiting page
- **First User**: Register → Admin + Approved → Token → Dashboard
- **Login Attempts**: Blocked for pending/rejected users with clear messages
- **Admin Approval**: Changes status to 'approved', user can then login

### Security
- All protected routes require approved status
- Admin routes require both approved status AND admin role
- JWT tokens include role and status for efficient authorization
- Middleware chains ensure proper access control

## Manual Steps Required

### 1. Run Migrations
Migrations run automatically when the backend starts. Verify with:
```bash
psql $DATABASE_URL -c "\d users"
```

Expected output should show `role` and `status` columns.

### 2. Promote Specific User to Admin
To promote the user with ID `5450ee66-7833-4e7d-a967-7d8f8b9e064d`:

```bash
# Connect to database
psql $DATABASE_URL

# Run the SQL script
\i backend/scripts/promote_admin.sql

# Or manually:
UPDATE users
SET role = 'admin', status = 'approved'
WHERE id = '5450ee66-7833-4e7d-a967-7d8f8b9e064d';
```

Verify:
```sql
SELECT id, email, role, status FROM users WHERE role = 'admin';
```

### 3. Test Flows

**Test New User Registration:**
1. Register a new user
2. Should see "Waiting for Approval" page
3. Check database: `SELECT * FROM users ORDER BY created_at DESC LIMIT 1;`
4. Status should be 'pending', role should be 'user'

**Test Pending Login:**
1. Try to login with pending user
2. Should see error: "Your account is pending approval. Please wait."

**Test Admin Access:**
1. Login as admin user
2. Should see "Admin" link in header
3. Navigate to /admin
4. Should see User Management panel

**Test Approval Flow:**
1. As admin, view pending users
2. Click "Approve" on a user
3. User status changes to 'approved'
4. User can now login successfully

**Test Rejected Login:**
1. As admin, reject a user
2. Try to login as rejected user
3. Should see error: "Your account has been rejected"

## Database Schema

### Users Table (Updated)
```sql
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  preferences TEXT,
  language VARCHAR(10) NOT NULL DEFAULT 'en',
  role VARCHAR(20) NOT NULL DEFAULT 'user',
  status VARCHAR(20) NOT NULL DEFAULT 'pending',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT check_role CHECK (role IN ('user', 'admin')),
  CONSTRAINT check_status CHECK (status IN ('pending', 'approved', 'rejected'))
);

CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_role ON users(role);
```

## API Endpoints

### New Admin Endpoints
- `GET /api/admin/users` - List all users (admin only)
- `GET /api/admin/users/pending` - List pending users (admin only)
- `PUT /api/admin/users/:id/approve` - Approve user (admin only)
- `PUT /api/admin/users/:id/reject` - Reject user (admin only)

### Modified Endpoints
- `POST /api/auth/register` - Returns RegisterResponse with status field
- `POST /api/auth/login` - Checks user status before allowing login
- All protected routes now require 'approved' status

## JWT Token Structure

Tokens now include:
```json
{
  "user_id": "uuid",
  "role": "admin|user",
  "status": "pending|approved|rejected",
  "exp": 1234567890
}
```

## Frontend Routes

### New Routes
- `/waiting-approval` - Public page for pending users
- `/admin` - Admin panel (admin only)

### Protected Routes
All existing protected routes now require approved status:
- `/dashboard`
- `/proposals/*`
- `/shopping-list`
- `/preferences`

## Files Created/Modified

### Backend Files Created
- `backend/migrations/004_add_roles_and_approval.up.sql`
- `backend/migrations/004_add_roles_and_approval.down.sql`
- `backend/scripts/promote_admin.sql`
- `backend/internal/services/admin_service.go`
- `backend/internal/handlers/admin_handler.go`
- `backend/internal/middleware/auth.go` (replaced)

### Backend Files Modified
- `backend/internal/models/user.go`
- `backend/internal/repository/user_repository.go`
- `backend/internal/services/auth_service.go`
- `backend/internal/handlers/auth_handler.go`
- `backend/cmd/api/main.go`

### Frontend Files Created
- `frontend/src/components/AdminRoute.tsx`
- `frontend/src/pages/WaitingApproval.tsx`
- `frontend/src/pages/AdminPanel.tsx`
- `frontend/src/services/adminService.ts`

### Frontend Files Modified
- `frontend/src/types/index.ts`
- `frontend/src/context/AuthContext.tsx`
- `frontend/src/services/authService.ts`
- `frontend/src/components/Header.tsx`
- `frontend/src/pages/Register.tsx`
- `frontend/src/pages/Login.tsx`
- `frontend/src/App.tsx`

## Build Status

✅ Backend compiles successfully
✅ Frontend builds successfully
✅ All TypeScript types validated
✅ No compilation errors

## Next Steps

1. Start the backend server
2. Verify migrations ran successfully
3. Promote the specific user to admin using the SQL script
4. Test all user flows (registration, approval, login, admin panel)
5. Deploy to production if all tests pass

## Notes

- The system prevents deleted/rejected users from re-registering (email remains in database)
- Admin users can see their own admin panel link in the header
- Non-admin users attempting to access /admin are redirected to dashboard
- Pending users see a friendly waiting page instead of being stuck at login
- Clear, specific error messages guide users through each status scenario
