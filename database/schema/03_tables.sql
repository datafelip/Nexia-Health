-- ===========================================================
-- Nexia Health
-- 03_tables.sql
-- Criação das tabelas do sistema
-- ===========================================================


-- ###########################################################
-- TABELAS INDEPENDENTES
-- Estas tabelas não possuem chaves estrangeiras e devem ser
-- criadas primeiro.
-- ###########################################################


-- ===========================================================
-- Tabela: Especialidade
-- ===========================================================

CREATE TABLE especialidade (
    id_especialidade UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    nome VARCHAR(100) NOT NULL UNIQUE,

    descricao TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);


-- ===========================================================
-- Tabela: Setor
-- ===========================================================

CREATE TABLE setor (
    id_setor UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    nome VARCHAR(100) NOT NULL UNIQUE,

    descricao TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);


-- ===========================================================
-- Tabela: Usuario
-- ===========================================================

CREATE TABLE usuario (
    id_usuario UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    nome VARCHAR(150) NOT NULL,

    email VARCHAR(255) NOT NULL UNIQUE,

    senha_hash TEXT NOT NULL,

    perfil perfil_usuario_enum NOT NULL,

    status status_usuario_enum NOT NULL DEFAULT 'Ativo',

    ultimo_login TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CHECK (char_length(nome) >= 3)
);


-- ===========================================================
-- Tabela: Paciente
-- ===========================================================

CREATE TABLE paciente (
    id_paciente UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    nome VARCHAR(150) NOT NULL,

    cpf CHAR(11) NOT NULL UNIQUE,

    data_nascimento DATE NOT NULL,

    sexo sexo_enum NOT NULL,

    telefone VARCHAR(15),

    endereco TEXT,

    cartao_sus CHAR(15) UNIQUE,

    data_cadastro TIMESTAMP NOT NULL DEFAULT NOW(),

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CHECK (char_length(nome) >= 3)
);


-- ===========================================================
-- Tabela: Exame
-- ===========================================================

CREATE TABLE exame (
    id_exame UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    nome VARCHAR(100) NOT NULL UNIQUE,

    descricao TEXT,

    preparo TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);



-- ###########################################################
-- TABELAS DEPENDENTES
-- Estas tabelas possuem chaves estrangeiras e dependem das
-- tabelas acima.
-- ###########################################################


-- ===========================================================
-- Tabela: Profissional
-- Depende de:
-- • Usuario
-- • Especialidade
-- • Setor
-- ===========================================================

CREATE TABLE profissional (
    id_profissional UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    id_usuario UUID NOT NULL,

    id_especialidade UUID NOT NULL,

    id_setor UUID NOT NULL,

    cpf CHAR(11) NOT NULL UNIQUE,

    conselho_profissional VARCHAR(50) NOT NULL,

    numero_registro VARCHAR(50) NOT NULL UNIQUE,

    telefone VARCHAR(15),

    status status_profissional_enum NOT NULL DEFAULT 'Ativo',

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_profissional_usuario
        FOREIGN KEY (id_usuario)
        REFERENCES usuario(id_usuario)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_profissional_especialidade
        FOREIGN KEY (id_especialidade)
        REFERENCES especialidade(id_especialidade)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_profissional_setor
        FOREIGN KEY (id_setor)
        REFERENCES setor(id_setor)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);


-- ===========================================================
-- Tabela: Triagem
-- Depende de:
-- • Paciente
-- • Profissional
-- ===========================================================

CREATE TABLE triagem (
    id_triagem UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    id_paciente UUID NOT NULL,

    id_profissional UUID NOT NULL,

    classificacao_risco classificacao_risco_enum NOT NULL,

    pressao_arterial VARCHAR(10),

    temperatura NUMERIC(4,1),

    frequencia_cardiaca SMALLINT,

    observacoes TEXT,

    data_triagem TIMESTAMP NOT NULL DEFAULT NOW(),

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CHECK (temperatura IS NULL OR temperatura BETWEEN 30.0 AND 45.0),

    CHECK (
        frequencia_cardiaca IS NULL
        OR frequencia_cardiaca BETWEEN 20 AND 250
    ),

    CONSTRAINT fk_triagem_paciente
        FOREIGN KEY (id_paciente)
        REFERENCES paciente(id_paciente)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_triagem_profissional
        FOREIGN KEY (id_profissional)
        REFERENCES profissional(id_profissional)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);


-- ===========================================================
-- Tabela: Fila
-- Depende de:
-- • Paciente
-- • Triagem
-- ===========================================================

CREATE TABLE fila (
    id_fila UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    id_paciente UUID NOT NULL,

    id_triagem UUID NOT NULL,

    posicao INTEGER NOT NULL,

    status status_fila_enum NOT NULL DEFAULT 'Aguardando',

    data_entrada TIMESTAMP NOT NULL DEFAULT NOW(),

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CHECK (posicao > 0),

    CONSTRAINT fk_fila_paciente
        FOREIGN KEY (id_paciente)
        REFERENCES paciente(id_paciente)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_fila_triagem
        FOREIGN KEY (id_triagem)
        REFERENCES triagem(id_triagem)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);


-- ===========================================================
-- Tabela: Agendamento
-- Depende de:
-- • Paciente
-- ===========================================================

CREATE TABLE agendamento (
    id_agendamento UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    id_paciente UUID NOT NULL,

    tipo tipo_agendamento_enum NOT NULL,

    data_agendada DATE NOT NULL,

    horario TIME NOT NULL,

    status status_agendamento_enum NOT NULL DEFAULT 'Agendado',

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_agendamento_paciente
        FOREIGN KEY (id_paciente)
        REFERENCES paciente(id_paciente)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);


-- ===========================================================
-- Tabela: Atendimento
-- Depende de:
-- • Paciente
-- • Profissional
-- • Triagem
-- • Agendamento
-- • Especialidade
-- • Setor
-- ===========================================================

CREATE TABLE atendimento (
    id_atendimento UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    id_paciente UUID NOT NULL,

    id_profissional UUID NOT NULL,

    id_triagem UUID NOT NULL,

    id_agendamento UUID,

    id_especialidade UUID NOT NULL,

    id_setor UUID NOT NULL,

    data_inicio TIMESTAMP NOT NULL,

    data_fim TIMESTAMP,

    diagnostico TEXT,

    observacoes TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CHECK (
        data_fim IS NULL
        OR data_fim >= data_inicio
    ),

    CONSTRAINT fk_atendimento_paciente
        FOREIGN KEY (id_paciente)
        REFERENCES paciente(id_paciente)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_atendimento_profissional
        FOREIGN KEY (id_profissional)
        REFERENCES profissional(id_profissional)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_atendimento_triagem
        FOREIGN KEY (id_triagem)
        REFERENCES triagem(id_triagem)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_atendimento_agendamento
        FOREIGN KEY (id_agendamento)
        REFERENCES agendamento(id_agendamento)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_atendimento_especialidade
        FOREIGN KEY (id_especialidade)
        REFERENCES especialidade(id_especialidade)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_atendimento_setor
        FOREIGN KEY (id_setor)
        REFERENCES setor(id_setor)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);


-- ===========================================================
-- Tabela: Solicitacao_Exame
-- Depende de:
-- • Atendimento
-- • Exame
-- ===========================================================

CREATE TABLE solicitacao_exame (
    id_solicitacao UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    id_atendimento UUID NOT NULL,

    id_exame UUID NOT NULL,

    data_solicitacao TIMESTAMP NOT NULL DEFAULT NOW(),

    status status_exame_enum NOT NULL DEFAULT 'Solicitado',

    resultado TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_solicitacao_atendimento
        FOREIGN KEY (id_atendimento)
        REFERENCES atendimento(id_atendimento)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_solicitacao_exame
        FOREIGN KEY (id_exame)
        REFERENCES exame(id_exame)
        ON DELETE CASCADE
        ON UPDATE cascade
);