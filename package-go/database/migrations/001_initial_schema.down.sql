-- Drop all foreign key constraints first
ALTER TABLE IF EXISTS public."gcasUser" DROP CONSTRAINT IF EXISTS "gcasUser_organizationCategoryId_fkey";
ALTER TABLE IF EXISTS public."gcasGroupUserRelation" DROP CONSTRAINT IF EXISTS "gcasGroupUserRelation_gcasUserId_fkey";
ALTER TABLE IF EXISTS public."gcasGroupUserRelation" DROP CONSTRAINT IF EXISTS "gcasGroupUserRelation_groupId_fkey";
ALTER TABLE IF EXISTS public."gcasGroupUserRelation" DROP CONSTRAINT IF EXISTS "gcasGroupUserRelation_userRoleId_fkey";
ALTER TABLE IF EXISTS public."gcasGroupSystemRelation" DROP CONSTRAINT IF EXISTS "gcasGroupSystemRelation_systemId_fkey";
ALTER TABLE IF EXISTS public."gcasGroupSystemRelation" DROP CONSTRAINT IF EXISTS "gcasGroupSystemRelation_groupId_fkey";
ALTER TABLE IF EXISTS public.project DROP CONSTRAINT IF EXISTS "project_localGovernmentId_fkey";
ALTER TABLE IF EXISTS public."projectCost" DROP CONSTRAINT IF EXISTS "projectCost_projectId_fkey";
ALTER TABLE IF EXISTS public."projectSystemRelation" DROP CONSTRAINT IF EXISTS "projectSystemRelation_projectId_fkey";
ALTER TABLE IF EXISTS public."projectSystemRelation" DROP CONSTRAINT IF EXISTS "projectSystemRelation_systemId_fkey";
ALTER TABLE IF EXISTS public.system DROP CONSTRAINT IF EXISTS "system_localGovernmentId_fkey";
ALTER TABLE IF EXISTS public."systemBasicInformation" DROP CONSTRAINT IF EXISTS "systemBasicInformation_projectId_fkey";

-- Drop indexes
DROP INDEX IF EXISTS public."gcasUser_mailAddress_unique";
DROP INDEX IF EXISTS public."gcasGroup_groupName_unique";
DROP INDEX IF EXISTS public."system_systemName_unique";

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS public."systemBasicInformation";
DROP TABLE IF EXISTS public."projectSystemRelation";
DROP TABLE IF EXISTS public."projectCost";
DROP TABLE IF EXISTS public.project;
DROP TABLE IF EXISTS public."gcasGroupSystemRelation";
DROP TABLE IF EXISTS public."gcasGroupUserRelation";
DROP TABLE IF EXISTS public."gcasGroup";
DROP TABLE IF EXISTS public."gcasUser";
DROP TABLE IF EXISTS public.system;

-- Drop master tables
DROP TABLE IF EXISTS public."m_userRole";
DROP TABLE IF EXISTS public."m_organizationCategory";
DROP TABLE IF EXISTS public."m_localGovernment";

-- Drop sequences
DROP SEQUENCE IF EXISTS public."m_userRole_id_seq";
DROP SEQUENCE IF EXISTS public."m_organizationCategory_id_seq";

-- Drop drizzle migration tracking table and schema
DROP TABLE IF EXISTS drizzle.__drizzle_migrations;
DROP SEQUENCE IF EXISTS drizzle.__drizzle_migrations_id_seq;
DROP SCHEMA IF EXISTS drizzle;

