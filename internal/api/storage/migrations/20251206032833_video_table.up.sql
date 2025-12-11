CREATE TABLE course_videos (
    video_id SERIAL PRIMARY KEY,
    course_id INT REFERENCES courses(course_id) ON DELETE CASCADE,
    video_title TEXT NOT NULL,
    video_filename TEXT NOT NULL, 
    uploaded_at TIMESTAMP DEFAULT NOW()
);
