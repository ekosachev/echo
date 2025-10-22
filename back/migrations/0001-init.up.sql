CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    timezone TEXT DEFAULT 'UTC',
    created_at TIMESTAMP DEFAULT now()
);
CREATE INDEX idx_users_email ON users(email);


CREATE TABLE calendars (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    color TEXT,
    created_at TIMESTAMP DEFAULT now()
);
CREATE INDEX idx_calendars_owner ON calendars(owner_id);


CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    calendar_id UUID NOT NULL REFERENCES calendars(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    location TEXT,
    start_at TIMESTAMP NOT NULL,
    end_at TIMESTAMP NOT NULL,
    timezone TEXT DEFAULT 'UTC',
    all_day BOOLEAN DEFAULT FALSE,
    priority INT DEFAULT 0,
    color TEXT,
    meta JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
CREATE INDEX idx_events_calendar ON events(calendar_id);
CREATE INDEX idx_events_start_end ON events(start_at, end_at);


CREATE TABLE event_rrule (
    event_id UUID PRIMARY KEY REFERENCES events(id) ON DELETE CASCADE,
    rrule TEXT,
    rdate TIMESTAMP[],
    exdate TIMESTAMP[],
    until TIMESTAMP,
    count INT
);


CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    calendar_id UUID NOT NULL REFERENCES calendars(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    due_at TIMESTAMP,
    timezone TEXT DEFAULT 'UTC',
    status TEXT CHECK (status IN ('open', 'in_progress', 'done')) DEFAULT 'open',
    priority INT DEFAULT 0,
    checklist JSONB DEFAULT '[]'::jsonb,
    tags TEXT[],
    created_at TIMESTAMP DEFAULT now()
);
CREATE INDEX idx_tasks_calendar ON tasks(calendar_id);
CREATE INDEX idx_tasks_due_at ON tasks(due_at);


CREATE TABLE reminders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID REFERENCES events(id) ON DELETE CASCADE,
    task_id UUID REFERENCES tasks(id) ON DELETE CASCADE,
    type TEXT CHECK (type IN ('relative', 'absolute')),
    offset_minutes INT,
    fire_at TIMESTAMP,
    channels TEXT[] DEFAULT ARRAY['web-push'],
    policy_id UUID
);
CREATE INDEX idx_reminders_event ON reminders(event_id);
CREATE INDEX idx_reminders_task ON reminders(task_id);


CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reminder_id UUID REFERENCES reminders(id) ON DELETE CASCADE,
    channel TEXT CHECK (channel IN ('web-push', 'email', 'telegram')) NOT NULL,
    status TEXT CHECK (status IN ('queued', 'sent', 'failed')) DEFAULT 'queued',
    attempt_count INT DEFAULT 0,
    last_attempt_at TIMESTAMP,
    error TEXT,
    created_at TIMESTAMP DEFAULT now()
);
CREATE INDEX idx_notifications_user ON notifications(user_id);
CREATE INDEX idx_notifications_status ON notifications(status);


CREATE TABLE bot_links (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    telegram_chat_id BIGINT NOT NULL,
    verified BOOLEAN DEFAULT FALSE,
    link_code TEXT,
    expires_at TIMESTAMP,
    PRIMARY KEY (user_id, telegram_chat_id)
);
CREATE INDEX idx_bot_links_code ON bot_links(link_code);
