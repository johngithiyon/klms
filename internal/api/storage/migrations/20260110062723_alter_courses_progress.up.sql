-- 1. Add course_id column
ALTER TABLE course_progress
ADD COLUMN course_id INT;

-- 2. Set course_id based on course_name and courses.title
UPDATE course_progress cp
SET course_id = c.course_id
FROM courses c
WHERE cp.course_name = c.title;

-- 3. Make course_id NOT NULL if you want to enforce it
ALTER TABLE course_progress
ALTER COLUMN course_id SET NOT NULL;

-- 4. Add foreign key constraint
ALTER TABLE course_progress
ADD CONSTRAINT fk_course
FOREIGN KEY (course_id)
REFERENCES courses(course_id)
ON DELETE CASCADE;
