
import type { Paciente } from '../types/paciente'

const NOMES_DE_TESTE = new Set([
  'paciente teste',
  'paciente teste 2',
])

export function ehPacienteDeTeste(paciente: Paciente): boolean {
  const nomeNormalizado = paciente.nome
    .trim()
    .toLocaleLowerCase('pt-BR')
    .replace(/\s+/g, ' ')

  return NOMES_DE_TESTE.has(nomeNormalizado)
}

export function filtrarPacientesVisiveis(
  pacientes: Paciente[]
): Paciente[] {
  return pacientes.filter(
    (paciente) => !ehPacienteDeTeste(paciente)
  )
}
