-- Seed data which will not be moved to production
-- +goose Up

-- +goose StatementBegin
INSERT INTO issues (id, title, repoId, url, tags, difficulty, resolved) VALUES
(
  gen_random_uuid(),
  'Add support for group chats',
  (SELECT id FROM repository WHERE name = 'Gochaat: A Chat App in Go'),
  'https://github.com/aadit-n3rdy/gochaat/issues/1',
  '{"Go", "WebSockets"}',
  'Medium',
  false
),
(
  gen_random_uuid(), 
  'Implement basic authentication',
  (SELECT id FROM repository WHERE name = 'Gochaat: A Chat App in Go'),
  'https://github.com/aadit-n3rdy/gochaat/issues/2',
  '{"Go", "Auth"}',
  'Easy',
  false
),
(
  gen_random_uuid(),
  'Add Bloom filters for efficient key lookups',
  (SELECT id FROM repository WHERE name = 'Whodis: Keys, Values and Concurrency'),
  'https://github.com/aadit-n3rdy/whodis-server/issues/1', 
  '{"Java", "Data Structures"}',
  'Hard',
  false
),
(
  gen_random_uuid(),
  'Write Java client library',
  (SELECT id FROM repository WHERE name = 'Whodis: Keys, Values and Concurrency'),
  'https://github.com/aadit-n3rdy/whodis-server/issues/2',
  '{"Java"}',
  'Medium', 
  false
),
(
  gen_random_uuid(),
  'Add support for TCP protocol testing',
  (SELECT id FROM repository WHERE name = 'Load-Pulse: Load Testing Tool'),
  'https://github.com/Naganathan05/Load-Pulse/issues/1',
  '{"Go", "TCP"}',
  'Medium',
  false
),
(
  gen_random_uuid(),
  'Implement Kubernetes deployment',
  (SELECT id FROM repository WHERE name = 'Load-Pulse: Load Testing Tool'), 
  'https://github.com/Naganathan05/Load-Pulse/issues/2',
  '{"Kubernetes", "DevOps"}',
  'Hard',
  false
),
(
  gen_random_uuid(),
  'Add support for persistent storage',
  (SELECT id FROM repository WHERE name = 'NodeGainsDB: Ligthweight In-Memory GraphDB'),
  'https://github.com/Astrasv/NodeGainsDB/issues/1',
  '{"Python", "Databases"}',
  'Medium',
  false
),
(
  gen_random_uuid(),
  'Implement text-to-speech for multiple languages',
  (SELECT id FROM repository WHERE name = 'Notivos-AI'),
  'https://github.com/adithya-menon-r/Notivos-AI/issues/1',
  '{"JavaScript", "TTS"}',
  'Easy',
  false
),
(
  gen_random_uuid(),
  'Add support for collaborative note editing',
  (SELECT id FROM repository WHERE name = 'Notivos-AI'),
  'https://github.com/adithya-menon-r/Notivos-AI/issues/2', 
  '{"JavaScript", "WebSockets"}',
  'Hard',
  false
);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO issue_claims (id, ghUsername, issue_id, claimed_on, elapsed_on) VALUES
(
  gen_random_uuid(),
  'IAmRiteshKoushik',
  (SELECT id FROM issues WHERE url = 'https://github.com/aadit-n3rdy/gochaat/issues/1'),
  NOW(),
  NOW() + INTERVAL '7 days'
),
(
  gen_random_uuid(), 
  'Ashrockzzz2003',
  (SELECT id FROM issues WHERE url = 'https://github.com/aadit-n3rdy/gochaat/issues/1'),
  NOW(),
  NOW() + INTERVAL '7 days'
),
(
  gen_random_uuid(),
  'IAmRiteshKoushik',
  (SELECT id FROM issues WHERE url = 'https://github.com/Astrasv/NodeGainsDB/issues/1'),
  NOW() - INTERVAL '2 days',
  NOW() + INTERVAL '5 days'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE user_account CASCADE;
TRUNCATE TABLE issues CASCADE;
-- +goose StatementEnd