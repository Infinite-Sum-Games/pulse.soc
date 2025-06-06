-- Seed data that should be cleaned from production

-- +goose Up
-- +goose StatementBegin
INSERT INTO user_account (first_name, middle_name, last_name, email, ghUsername, status, bounty) VALUES
('Ritesh', NULL, 'Koushik', 'cb.en.u4cse22038@cb.students.amrita.edu', 'IAmRiteshKoushik', true, 0),
('Ashwin', 'Narayanan', 'S', 'cb.en.u4cse21004', 'Ashrockzzz2003', true, 0);
-- +goose StatementEnd
