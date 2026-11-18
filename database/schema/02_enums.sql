-- ===========================================================
-- Nexia Health
-- 02_enums.sql
-- Definição dos tipos ENUM utilizados no sistema
-- ===========================================================

-- Sexo do paciente
create type sexo_enum as ENUM(
	'Masculino',
	'Feminino',
	'Outro'
);

-- Perfil do usuário
create type perfil_usuario_enum as ENUM(
	'Administrador',
	'Recepcionista',
	'Enfermeiro',
	'Médico',
	'Gestor'
);

-- Status do usuário
create type status_usuario_enum as ENUM(
	'Ativo',
	'Inativo',
	'Bloqueado'
);

-- Status do profissional
create type status_profissional_enum as ENUM(
	'Ativo',
	'Inativo',
	'Bloqueado'
);


-- Classificação de risco
CREATE TYPE classificacao_risco_enum AS ENUM (
    'Vermelho',
    'Amarelo',
    'Verde'
);

-- Status da fila
    CREATE TYPE status_fila_enum AS ENUM (
        'Aguardando',
        'Chamado',
        'Em Atendimento',
        'Finalizado',
        'Cancelado'
    );

-- Tipo de agendamento
CREATE TYPE tipo_agendamento_enum AS ENUM (
    'Consulta',
    'Exame'
);

-- Status do agendamento
CREATE TYPE status_agendamento_enum AS ENUM (
    'Agendado',
    'Confirmado',
    'Cancelado',
    'Realizado'
);

-- Status da solicitação de exame
CREATE TYPE status_exame_enum AS ENUM (
    'Solicitado',
    'Agendado',
    'Realizado',
    'Cancelado'
);