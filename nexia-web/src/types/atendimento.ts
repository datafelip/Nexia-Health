
export interface Atendimento {
  id_atendimento: string

  id_paciente: string
  id_profissional: string
  id_triagem: string

  id_agendamento?: string | null

  id_especialidade: string
  id_setor: string

  data_inicio: string
  data_fim?: string | null

  diagnostico?: string | null
  observacoes?: string | null

  created_at: string
  updated_at: string
}

// Dados utilizados para iniciar um atendimento.
// Os campos obrigatórios correspondem aos vínculos
// presentes no modelo do backend.
export interface AtendimentoPayload {
  id_paciente: string
  id_profissional: string
  id_triagem: string

  id_agendamento?: string | null

  id_especialidade: string
  id_setor: string

  diagnostico?: string | null
  observacoes?: string | null
}

// Representação de um atendimento enriquecido
// com informações de outros módulos.
export interface AtendimentoDetalhado {
  atendimento: Atendimento

  nomePaciente: string
  nomeProfissional: string

  classificacaoRisco:
    | 'Vermelho'
    | 'Amarelo'
    | 'Verde'
    | null

  statusFila:
    | 'Aguardando'
    | 'Chamado'
    | 'Em Atendimento'
    | 'Finalizado'
    | 'Cancelado'
    | null
}
