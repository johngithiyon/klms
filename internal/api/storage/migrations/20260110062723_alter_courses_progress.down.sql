-- 1. Drop foreign key constraint
ALTER TABLE course_progress
DROP CONSTRAINT IF EXISTS fk_course;

-- 2. Drop the course_id column
ALTER TABLE course_progress
DROP COLUMN IF EXISTS course_id;
