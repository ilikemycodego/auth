




-- ======================================
-- AUTH-SERVICE
-- ======================================

-- Таблица одноразовых кодов для входа по email
CREATE TABLE IF NOT EXISTS auth_codes (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    code VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    used BOOLEAN DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_auth_codes_email ON auth_codes(email);

-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL DEFAULT 'amigo',
    role VARCHAR(50) NOT NULL DEFAULT 'boss', -- по умолчанию мастер-одиночка/начальник
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, blocked, pending
    device_id VARCHAR(255), -- опционально для привязки к устройству
    created_at TIMESTAMP DEFAULT NOW()
);

-- Таблица сессий
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id VARCHAR(255),  -- опционально привязка к устройству
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);

-- ======================================
-- COMPANY-SERVICE
-- ======================================

-- Таблица компаний
CREATE TABLE IF NOT EXISTS companies (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL REFERENCES users(id),  -- создатель/директор
    name VARCHAR(255) NOT NULL,
    plan VARCHAR(50) NOT NULL DEFAULT 'free',   -- тариф компании
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, blocked, trial
    created_at TIMESTAMP DEFAULT NOW()
);

-- Таблица сотрудников компании
CREATE TABLE IF NOT EXISTS company_members (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'master',  -- boss, master, admin
    created_at TIMESTAMP DEFAULT NOW()
);

-- Таблица приглашений сотрудников (optional)
CREATE TABLE IF NOT EXISTS company_invitations (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'master',
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, accepted, declined
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_company_invite_email ON company_invitations(email);

-- ======================================
-- Примечания:
-- 1) Auth-service хранит пользователей, их роли по умолчанию и сессии.
-- 2) Company-service хранит бизнес-логику: компании, сотрудников, роли, тарифы и статусы.
-- 3) device_id в users и sessions позволяет ограничивать триал и контролировать устройства.
-- 4) Приглашения нужны для безопасного добавления сотрудников без отдельной регистрации.