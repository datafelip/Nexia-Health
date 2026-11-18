-- ===========================================================
-- Nexia Health
-- 04_indexes.sql
-- Índices adicionais para otimização das consultas
-- ===========================================================


-- ###########################################################
-- TABELA: PROFISSIONAL
-- ###########################################################

-- Índice para consultas por usuário
CREATE INDEX idx_profissional_usuario
ON profissional (id_usuario);

-- Índice para consultas por especialidade
CREATE INDEX idx_profissional_especialidade
ON profissional (id_especialidade);

-- Índice para consultas por setor
CREATE INDEX idx_profissional_setor
ON profissional (id_setor);

-- Índice para filtragem por status
CREATE INDEX idx_profissional_status
ON profissional (status);


-- ###########################################################
-- TABELA: TRIAGEM
-- ###########################################################

-- Índice para consultas de triagens de um paciente
CREATE INDEX idx_triagem_paciente
ON triagem (id_paciente);

-- Índice para consultas por profissional responsável
CREATE INDEX idx_triagem_profissional
ON triagem (id_profissional);

-- Índice para classificação de risco
CREATE INDEX idx_triagem_classificacao_risco
ON triagem (classificacao_risco);

-- Índice para consultas cronológicas
CREATE INDEX idx_triagem_data
ON triagem (data_triagem);


-- ###########################################################
-- TABELA: FILA
-- ###########################################################

-- Índice para localizar a fila de um paciente
CREATE INDEX idx_fila_paciente
ON fila (id_paciente);

-- Índice para relacionamentos com a triagem
CREATE INDEX idx_fila_triagem
ON fila (id_triagem);

-- Índice para consultas por status
CREATE INDEX idx_fila_status
ON fila (status);

-- Índice para ordenação da fila
CREATE INDEX idx_fila_posicao
ON fila (posicao);

-- Índice para consultas pela entrada na fila
CREATE INDEX idx_fila_data_entrada
ON fila (data_entrada);


-- ###########################################################
-- TABELA: AGENDAMENTO
-- ###########################################################

-- Índice para localizar agendamentos de um paciente
CREATE INDEX idx_agendamento_paciente
ON agendamento (id_paciente);

-- Índice para consultas por data
CREATE INDEX idx_agendamento_data
ON agendamento (data_agendada);

-- Índice para consultas por status
CREATE INDEX idx_agendamento_status
ON agendamento (status);

-- Índice para consultas por tipo
CREATE INDEX idx_agendamento_tipo
ON agendamento (tipo);

-- Índice para consultas por data e horário
CREATE INDEX idx_agendamento_data_horario
ON agendamento (data_agendada, horario);


-- ###########################################################
-- TABELA: ATENDIMENTO
-- ###########################################################

-- Índice para consultas de atendimentos de um paciente
CREATE INDEX idx_atendimento_paciente
ON atendimento (id_paciente);

-- Índice para consultas por profissional
CREATE INDEX idx_atendimento_profissional
ON atendimento (id_profissional);

-- Índice para consultas por triagem
CREATE INDEX idx_atendimento_triagem
ON atendimento (id_triagem);

-- Índice para consultas por agendamento
CREATE INDEX idx_atendimento_agendamento
ON atendimento (id_agendamento);

-- Índice para consultas por especialidade
CREATE INDEX idx_atendimento_especialidade
ON atendimento (id_especialidade);

-- Índice para consultas por setor
CREATE INDEX idx_atendimento_setor
ON atendimento (id_setor);

-- Índice para consultas cronológicas
CREATE INDEX idx_atendimento_data_inicio
ON atendimento (data_inicio);


-- ###########################################################
-- TABELA: SOLICITACAO_EXAME
-- ###########################################################

-- Índice para localizar solicitações de um atendimento
CREATE INDEX idx_solicitacao_exame_atendimento
ON solicitacao_exame (id_atendimento);

-- Índice para localizar solicitações de um exame
CREATE INDEX idx_solicitacao_exame_exame
ON solicitacao_exame (id_exame);

-- Índice para consultas por status
CREATE INDEX idx_solicitacao_exame_status
ON solicitacao_exame (status);

-- Índice para consultas cronológicas
CREATE INDEX idx_solicitacao_exame_data
ON solicitacao_exame (data_solicitacao);