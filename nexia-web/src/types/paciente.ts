
export type Sexo = 'Masculino' | 'Feminino' | 'Outro'

export interface Paciente {
  id: string
  nome: string
  cpf: string
  data_nascimento: string
  sexo: Sexo
  telefone?: string | null
  endereco?: string | null
  cartao_sus?: string | null
  data_cadastro: string
  created_at: string
  updated_at: string
}

export interface PacientePayload {
  nome: string
  cpf: string
  data_nascimento: string
  sexo: Sexo
  telefone?: string | null
  endereco?: string | null
  cartao_sus?: string | null
}
