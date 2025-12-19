CREATE TABLE course_progress (
    course_name TEXT NOT NULL,
    student_name TEXT NOT NULL,
    no_of_videos INT NOT NULL,
    progress INT DEFAULT 0,
    status TEXT DEFAULT 'not completed',
    PRIMARY KEY (course_name, student_name)
);