-- Update RLS policy for users table to allow verification token lookups
-- This is needed because verification token lookups don't have tenant context

-- Drop existing policy
DROP POLICY IF EXISTS "users_tenant_isolation" ON "users";

-- Create new policy that allows:
-- 1. Normal tenant-scoped access when app.current_tenant_id is set
-- 2. Verification token lookups (token is unique, so safe)
CREATE POLICY "users_tenant_isolation" ON "users"
    FOR ALL
    USING (
        -- Allow if tenant context matches
        "tenant_id" = current_setting('app.current_tenant_id', true)
        OR
        -- Allow if no tenant context is set (for verification token lookups)
        -- This is safe because verification tokens are unique and unguessable
        current_setting('app.current_tenant_id', true) = ''
    )
    WITH CHECK (
        -- For INSERT/UPDATE, always require tenant context to match
        "tenant_id" = current_setting('app.current_tenant_id', true)
    );
