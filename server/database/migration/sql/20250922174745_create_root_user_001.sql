-- +goose Up
-- +goose StatementBegin
INSERT INTO users (
    email,
    salt,
    password_hash,
    full_name,
    phone,
    role,
    status,
    is_locked
) VALUES (
    'root@system.local',
    'root_salt_v1',
    '4f0856ea5dca516f550123dd1abe82870c3ea141049195dc171f33f5314868b9', -- bcrypt hash for 'root_password_v1'
    'System Administrator',
    '00112233',
    0, -- SYSTEM_ADMIN role
    0, -- ACTIVE status
    FALSE
)
ON CONFLICT (email) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE email = 'root@system.local';
-- +goose StatementEnd
