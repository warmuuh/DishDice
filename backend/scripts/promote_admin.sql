UPDATE users
SET role = 'admin', status = 'approved'
WHERE id = '5450ee66-7833-4e7d-a967-7d8f8b9e064d';

SELECT id, email, role, status FROM users WHERE role = 'admin';
