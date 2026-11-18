
-- ===========================================================
-- Nexia Health
-- Migração 001 - Cancelamento de Triagens
-- Preservação do histórico clínico
-- ===========================================================

BEGIN;

-- Data e horário do cancelamento.
ALTER TABLE triagem
ADD COLUMN IF NOT EXISTS cancelada_em TIMESTAMP;

-- Justificativa do cancelamento.
ALTER TABLE triagem
ADD COLUMN IF NOT EXISTS motivo_cancelamento TEXT;

-- Responsável pelo cancelamento.
-- Será preenchido futuramente por autenticação confiável.
ALTER TABLE triagem
ADD COLUMN IF NOT EXISTS cancelada_por UUID;

-- Garante consistência entre a data e o motivo.
-- O responsável poderá ser nulo enquanto não houver autenticação.
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'chk_triagem_cancelamento'
          AND conrelid = 'triagem'::regclass
    ) THEN
        ALTER TABLE triagem
        ADD CONSTRAINT chk_triagem_cancelamento
        CHECK (
            (
                cancelada_em IS NULL
                AND motivo_cancelamento IS NULL
                AND cancelada_por IS NULL
            )
            OR
            (
                cancelada_em IS NOT NULL
                AND NULLIF(TRIM(motivo_cancelamento), '') IS NOT NULL
            )
        );
    END IF;
END $$;

-- Referência opcional ao usuário responsável.
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_triagem_cancelada_por'
          AND conrelid = 'triagem'::regclass
    ) THEN
        ALTER TABLE triagem
        ADD CONSTRAINT fk_triagem_cancelada_por
        FOREIGN KEY (cancelada_por)
        REFERENCES usuario(id_usuario)
        ON DELETE RESTRICT;
    END IF;
END $$;

-- Índice para consultas de triagens canceladas.
CREATE INDEX IF NOT EXISTS idx_triagem_cancelada_em
ON triagem (cancelada_em)
WHERE cancelada_em IS NOT NULL;

COMMIT;
