-- Таблица для комментариев к проблемам
CREATE TABLE IF NOT EXISTS problem_comments (
    id SERIAL PRIMARY KEY,
    problem_id INT NOT NULL REFERENCES problem_shipments(id) ON DELETE CASCADE,
    admin_id INT NOT NULL REFERENCES users(id),
    comment TEXT NOT NULL,
    comment_type VARCHAR(50) DEFAULT 'comment', -- 'comment', 'status_change', 'assignment', 'resolution'
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Индексы для быстрого поиска
CREATE INDEX idx_problem_comments_problem_id ON problem_comments(problem_id);
CREATE INDEX idx_problem_comments_admin_id ON problem_comments(admin_id);
CREATE INDEX idx_problem_comments_created_at ON problem_comments(created_at);

-- Таблица для истории изменения статусов проблем
CREATE TABLE IF NOT EXISTS problem_status_history (
    id SERIAL PRIMARY KEY,
    problem_id INT NOT NULL REFERENCES problem_shipments(id) ON DELETE CASCADE,
    admin_id INT REFERENCES users(id),
    old_status VARCHAR(50),
    new_status VARCHAR(50) NOT NULL,
    old_assigned_to INT REFERENCES users(id),
    new_assigned_to INT REFERENCES users(id),
    comment TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Индексы для истории статусов
CREATE INDEX idx_problem_status_history_problem_id ON problem_status_history(problem_id);
CREATE INDEX idx_problem_status_history_created_at ON problem_status_history(created_at);

-- Триггер для автоматического создания записи в истории при изменении статуса или назначения
CREATE OR REPLACE FUNCTION track_problem_status_changes()
RETURNS TRIGGER AS $$
BEGIN
    -- Записываем изменения статуса или назначения
    IF (OLD.status IS DISTINCT FROM NEW.status) OR (OLD.assigned_to IS DISTINCT FROM NEW.assigned_to) THEN
        INSERT INTO problem_status_history (
            problem_id, 
            old_status, 
            new_status,
            old_assigned_to,
            new_assigned_to,
            metadata
        ) VALUES (
            NEW.id,
            OLD.status,
            NEW.status,
            OLD.assigned_to,
            NEW.assigned_to,
            jsonb_build_object(
                'updated_by', 'system',
                'trigger_fired', NOW()
            )
        );
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Привязываем триггер к таблице
DROP TRIGGER IF EXISTS problem_status_change_trigger ON problem_shipments;
CREATE TRIGGER problem_status_change_trigger
    AFTER UPDATE ON problem_shipments
    FOR EACH ROW
    EXECUTE FUNCTION track_problem_status_changes();

-- Добавляем несколько тестовых комментариев
INSERT INTO problem_comments (problem_id, admin_id, comment, comment_type, metadata)
SELECT 
    1, -- ID проблемы (если существует)
    1, -- ID админа
    'Начато расследование проблемы. Связались с курьерской службой.',
    'comment',
    '{"priority": "high", "contact_method": "phone"}'
WHERE EXISTS (SELECT 1 FROM problem_shipments WHERE id = 1);

INSERT INTO problem_comments (problem_id, admin_id, comment, comment_type, metadata)
SELECT 
    1,
    1,
    'Статус изменен на "В работе". Назначен ответственный администратор.',
    'status_change',
    '{"old_status": "open", "new_status": "investigating"}'
WHERE EXISTS (SELECT 1 FROM problem_shipments WHERE id = 1);