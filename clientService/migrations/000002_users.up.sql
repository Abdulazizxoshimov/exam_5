CREATE TABLE IF NOT EXISTS jobs (
    id VARCHAR(64),
    name VARCHAR(64),
    company_name VARCHAR(64),
    start_date VARCHAR(64),
  location VARCHAR(200),
    end_date VARCHAR(64),
    status BOOLEAN,
    client_id VARCHAR(64),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

INSERT INTO jobs (id, name, company_name, start_date, end_date, status, client_id, location)
VALUES 
    ('a6dc102f-697f-46fa-8b00-cf64b3b7a679', 'Software Developer', 'ABC Inc.', '2024-04-23', '2024-04-23', true, '07b95949-e60b-4aec-9b61-b3ce480cd3cb', 'Toshkent'),
    ('b75d5c50-6496-4b88-8d04-cf9252f9156c', 'Data Analyst', 'XYZ Corp.', '2024-04-24', '2024-04-24', true, 'f1aa54cf-370d-4552-9c25-400a4c2c3feb', 'London'),
    ('8c066bb3-e50f-45d5-8b5d-2d82d422b8a2', 'Project Manager', '123 Company', '2024-04-25', '2024-04-25', false, '2ba4672f-36be-40c5-98a1-41c2d8689a74', 'Tokyo'),
    ('2af3a916-3640-48eb-b5b3-ea04256b1fb1', 'Marketing Specialist', '456 Corp.', '2024-04-26', '2024-04-26', true, 'fd85d050-6fde-4740-b44e-2017bc3ad77f', 'Washington'),
    ('3927b736-d2bf-4045-aa25-1e9ce10a8c0c', 'Financial Analyst', '789 Ltd.', '2024-04-27', '2024-04-27', false, '251f1fa2-ddd7-48d8-8662-05432bd610d8', 'Berlin'),
    ('52208aa1-7492-42a5-8bd0-53ce03dd6934', 'HR Coordinator', 'DEF Inc.', '2024-04-28', '2024-04-28', true, 'b4f3a915-cba7-400e-9ef6-e17aef650ef7', 'Paris'),
    ('6c4fa269-27ff-4050-b0c8-7c20f54a8e3c', 'Sales Representative', 'GHI Corp.', '2024-04-29', '2024-04-29', true, '8f0ba599-f9b2-4a86-bff6-50323b4dd698', 'Moskva'),
    ('18a2b53f-5806-49ce-a0c7-fbfb48dd3f8d', 'Customer Service Specialist', 'JKL Ltd.', '2024-04-30', '2024-04-30', false, '6c96cfb3-f169-40a4-b33f-16126e4b3aef', 'Nyu-York'),
    ('39f99512-0cc7-4c96-83bd-c96391369942', 'Accountant', 'MNO Inc.', '2024-05-01', '2024-05-01', true, '7261e9ee-9dae-4550-a358-d592d70e7108', 'Madrid'),
    ('5f28afbb-fdd1-4e03-ac9b-a1c4d04e661a', 'Graphic Designer', 'PQR Corp.', '2024-05-02', '2024-05-02', false, 'b8ac1071-b04d-4af0-88a7-9f8e332af12e', 'Rim');
