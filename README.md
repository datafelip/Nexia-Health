# Nexia Health

**Sistema de gestão de pacientes, triagens, filas e atendimentos, com componentes de inteligência artificial e previsão de demanda.**

O Nexia Health é um projeto de software para apoiar a organização do fluxo de atendimento em ambientes de saúde. A solução reúne uma API, uma interface web, um serviço de IA e um módulo de treinamento de modelos, com persistência em PostgreSQL.

> **Status:** projeto em desenvolvimento. A disponibilidade e a integração de cada funcionalidade devem ser verificadas no ambiente de execução.

## Visão geral

O projeto busca centralizar operações relacionadas ao atendimento, permitindo organizar informações e acompanhar etapas do fluxo clínico. Também explora o uso de dados sintéticos e aprendizado de máquina para análise de demanda.

### Principais módulos

- **Pacientes:** cadastro, consulta, atualização e exclusão lógica, com validações para impedir a exclusão enquanto houver fila ativa ou atendimento clínico em aberto.
- **Triagens:** registro e acompanhamento das informações de triagem e classificação de risco.
- **Filas:** gerenciamento de estados e transições do fluxo de atendimento.
- **Atendimentos:** integração do atendimento clínico ao fluxo operacional.
- **Gestão de apoio:** componentes para profissionais, setores, especialidades, exames, solicitações e agendamentos.
- **Inteligência artificial:** serviço Python e módulo de machine learning para experimentos de previsão de demanda.

> A lista descreve os módulos presentes na estrutura do projeto; não representa garantia de que todas as telas, rotas e integrações estejam finalizadas.

## Tecnologias

| Camada | Tecnologias |
| --- | --- |
| Backend | Go |
| Frontend | React, TypeScript, Vite |
| Banco de dados | PostgreSQL, SQL |
| IA e ML | Python, modelo serializado em Joblib |
| Documentação | Markdown, Draw.io, PDF, DOCX e ODS |
| Versionamento | Git e GitHub |

## Arquitetura

```text
                 +-------------------------+
                 | Frontend React / Vite   |
                 |        nexia-web        |
                 +------------+------------+
                              |
                              | HTTP / API
                              v
                 +-------------------------+
                 |      Backend Go         |
                 |        nexia-api        |
                 +------------+------------+
                              |
                              v
                 +-------------------------+
                 |       PostgreSQL        |
                 |         database        |
                 +-------------------------+

                 +-------------------------+
                 |   Serviço Python / IA   |
                 |         nexia-ai        |
                 +------------+------------+
                              |
                              v
                 +-------------------------+
                 | Treinamento e artefatos |
                 |         nexia-ml        |
                 +-------------------------+
```

O desenho é uma **visão conceitual** dos componentes, não um diagrama de implantação ou de conexões de rede verificadas.

## Estrutura do repositório

```text
Nexia-Health/
├── database/
│   ├── migrations/       # Evoluções do banco
│   ├── schema/           # Extensões, tipos, tabelas e índices
│   └── seeds/            # Dados iniciais
├── docs/
│   ├── Diagramas/
│   ├── Documentacao_Tecnica/
│   └── Projeto_Pesquisa/
├── nexia-ai/
│   ├── app/
│   └── requirements.txt
├── nexia-api/
│   ├── cmd/server/
│   ├── configs/
│   ├── internal/
│   │   ├── database/
│   │   ├── handlers/
│   │   ├── models/
│   │   ├── repositories/
│   │   ├── routes/
│   │   └── services/
│   ├── go.mod
│   └── go.sum
├── nexia-ml/
│   ├── datasets/
│   ├── models/
│   ├── gerar_dataset.py
│   └── treinar_modelo.py
├── nexia-web/
│   ├── public/
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── services/
│   │   ├── types/
│   │   └── utils/
│   └── package.json
└── README.md
```

## Pré-requisitos

Para executar os componentes localmente, instale as ferramentas correspondentes:

- Git;
- Go, em versão compatível com `nexia-api/go.mod`;
- Node.js e npm, em versões compatíveis com `nexia-web/package.json`;
- Python e pip, em versões compatíveis com as dependências de `nexia-ai/requirements.txt`;
- PostgreSQL, configurado para o projeto.

As versões exatas e as variáveis obrigatórias devem ser confirmadas nos arquivos de configuração.

## Configuração e execução local

### 1. Clonar o repositório

```bash
git clone https://github.com/datafelip/Nexia-Health.git
cd Nexia-Health
```

### 2. Preparar o banco de dados

Crie um banco PostgreSQL para o projeto e configure as credenciais **somente no ambiente local**. Os scripts estão organizados em `database/schema/`, `database/migrations/` e `database/seeds/`.

**Importante:** confira as dependências e a ordem de execução dos scripts antes de aplicá-los. Não execute seeds em um banco com dados reais sem revisar o conteúdo.

### 3. Iniciar a API Go

```bash
cd nexia-api
go mod download
# Configure as variáveis de ambiente exigidas pela aplicação.
go run ./cmd/server
```

O backend utiliza configuração de banco por variáveis de ambiente, incluindo `DB_PASSWORD`. Consulte `nexia-api/configs/config.go` para identificar todas as variáveis exigidas. Não versione arquivos `.env` reais.

### 4. Iniciar o frontend

Em outro terminal, a partir da raiz do projeto:

```bash
cd nexia-web
npm ci
npm run dev
```

Verifique em `nexia-web/src/services/` como a URL da API é configurada. A URL e a porta dependem do ambiente local.

### 5. Preparar o serviço de IA

Em outro terminal, a partir da raiz do projeto:

```bash
cd nexia-ai
python3 -m venv .venv
source .venv/bin/activate
python -m pip install -r requirements.txt
```

O ponto de entrada do serviço está em `nexia-ai/app/main.py`. Confirme no código o comando de inicialização, a porta e as variáveis de ambiente antes de executá-lo.

### 6. Machine learning

O diretório `nexia-ml/` contém scripts para geração de dados sintéticos e treinamento, além de métricas e um modelo serializado. Antes de executar os scripts, confira as dependências e os caminhos esperados em `gerar_dataset.py` e `treinar_modelo.py`.

**Segurança:** carregue arquivos `.joblib` apenas de fontes confiáveis. O carregamento de artefatos serializados pode executar código.

## Testes

Os testes do backend Go podem ser executados com:

```bash
cd nexia-api
go test ./... -v -count=1
```

Há testes de repositório para operações de pacientes e triagens. Alguns testes podem depender de um PostgreSQL de teste corretamente configurado. Não execute testes de integração sobre o banco de produção.

Para verificar os scripts disponíveis no frontend:

```bash
cd nexia-web
npm run
```

## Regras de negócio em destaque

### Exclusão lógica de pacientes

A exclusão de pacientes preserva o registro e marca sua exclusão por meio de `excluido_em`. A operação é bloqueada quando há vínculo com fila ativa ou atendimento clínico em aberto.

### Fluxo de filas

As transições de status seguem regras de negócio e validações no banco de dados. A lógica evita mudanças inválidas de estado e ajuda a manter a consistência do fluxo de atendimento.

## Documentação

Materiais complementares estão em [`docs/`](docs/), incluindo:

- Diagramas de entidades e relacionamentos, classes e casos de uso;
- Documentação técnica e arquitetura;
- Dicionário de dados;
- Documentos de pesquisa do projeto.

## Segurança e dados

- Não envie senhas, tokens, chaves privadas ou arquivos `.env` ao repositório.
- Utilize contas e bancos separados para desenvolvimento, testes e produção.
- Utilize dados sintéticos nos experimentos e evite publicar informações pessoais ou clínicas reais.
- Revise scripts de seed, documentos e artefatos de modelo antes de publicá-los.

## Estado e evolução

O Nexia Health está em evolução contínua. Melhorias de interface, integração entre serviços, testes, documentação e validações operacionais podem ser incorporadas conforme o desenvolvimento avançar.

---

**Nexia Health** — projeto de gestão do fluxo de atendimento em saúde com integração de tecnologias web, banco de dados e inteligência artificial.
