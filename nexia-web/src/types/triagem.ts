
export type ClassificacaoRisco =
  | 'Vermelho'
  | 'Amarelo'
  | 'Verde'

// Representa uma triagem retornada pela API Go.
export interface Triagem {
  id_triagem: string
  id_paciente: string
  id_profissional: string
  classificacao_risco: ClassificacaoRisco

  pressao_arterial?: string | null
  temperatura?: number | null
  frequencia_cardiaca?: number | null
  observacoes?: string | null

  data_triagem: string
  created_at: string
  updated_at: string

  // Cancelamento com preservação do histórico.
  cancelada_em?: string | null
  motivo_cancelamento?: string | null
  cancelada_por?: string | null
}

// Dados enviados para cadastrar uma triagem.
export interface TriagemPayload {
  id_paciente: string
  id_profissional: string
  classificacao_risco: ClassificacaoRisco

  pressao_arterial?: string | null
  temperatura?: number | null
  frequencia_cardiaca?: number | null
  observacoes?: string | null
}

// Dados enviados para reclassificar o risco.
export interface ReclassificarTriagemPayload {
  classificacao_risco: ClassificacaoRisco
}

// Dados enviados para cancelar uma triagem.
export interface CancelarTriagemPayload {
  motivo: string
}
