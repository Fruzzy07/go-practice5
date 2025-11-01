CREATE TABLE jobs(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    company VARCHAR(255) NOT NULL,
    salary INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);


INSERT INTO jobs (title, company, salary) VALUES
   ('Go Backend Developer', 'Kolesa', 600000),
   ('Frontend React Dev', 'Yandex', 550000),
   ('Mobile Kotlin Developer', 'Kaspi', 580000),
   ('QA Automation Engineer', 'Chocofamily', 500000),
   ('DevOps Engineer', 'Bi Group', 650000),
   ('Data Analyst', 'Halyk Bank', 530000),
   ('Product Manager', 'Kolesa', 800000),
   ('Python Backend Developer', 'Freedom Finance', 620000),
   ('SRE Engineer', 'Kazakhtelecom', 700000),
   ('Software Architect', 'Astana Hub', 900000);

