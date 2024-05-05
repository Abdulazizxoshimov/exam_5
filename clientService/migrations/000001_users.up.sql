CREATE TABLE IF NOT EXISTS clients (
    id VARCHAR(64),
    name VARCHAR(64),
    last_name VARCHAR(64),
    email VARCHAR(64),
    password VARCHAR(64),
    refresh_token TEXT,
    role   VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);


INSERT INTO clients (id, name, last_name, email, password, role, refresh_token)
VALUES 
    ('07b95949-e60b-4aec-9b61-b3ce480cd3cb', 'Abdulaziz', 'Xoshimov', 'abdulazizxoshimov22@gmail.com', '@Abdulaziz2004', 'user', 'refresh_token_123456'),
    ('f1aa54cf-370d-4552-9c25-400a4c2c3feb', 'Jane', 'Smith', 'abdulazizxoshimov687@gmail.com', '@Abdulaziz2004', 'admin', 'refresh_token_789012'),
    ('2ba4672f-36be-40c5-98a1-41c2d8689a74', 'Alice', 'Johnson', 'abdulazizxoshimov960@gmail.com', '@Abdulaziz2004', 'user', 'refresh_token_345678'),
    ('fd85d050-6fde-4740-b44e-2017bc3ad77f', 'Bob', 'Brown', 'golang962@gmail.com', 'passfdd@Abdulaziz2004ABC', 'user', 'refresh_token_901234'),
    ('251f1fa2-ddd7-48d8-8662-05432bd610d8', 'Charlie', 'Davis', 'charlie@example.com', 'passwordDEF', 'user', 'refresh_token_567890'),
    ('b4f3a915-cba7-400e-9ef6-e17aef650ef7', 'Emma', 'Wilson', 'emma@example.com', 'passwordGHI', 'user', 'refresh_token_234567'),
    ('8f0ba599-f9b2-4a86-bff6-50323b4dd698', 'David', 'Martinez', 'david@example.com', 'passwordJKL', 'user', 'refresh_token_890123'),
    ('6c96cfb3-f169-40a4-b33f-16126e4b3aef', 'Olivia', 'Anderson', 'olivia@example.com', 'e2sswordMNO', 'user', 'refresh_token_456789'),
    ('7261e9ee-9dae-4550-a358-d592d70e7108', 'Ethan', 'Taylor', 'ethan@example.com', 'passw34ordPQR', 'user', 'refresh_token_012345'),
    ('b8ac1071-b04d-4af0-88a7-9f8e332af12e', 'Sophia', 'Clark', 'sophia@example.com', 'passw45ordSTU', 'user', 'refresh_token_678901');
