
import { useState, type FormEvent } from 'react'
import axios from 'axios'
import {
  AlertCircle,
  CalendarDays,
  ClipboardList,
  HeartPulse,
  Loader2,
  MapPin,
  Save,
  ShieldCheck,
  UserRound,
  X,
} from 'lucide-react'

import { pacienteService } from '../../services/pacienteService'
import type {
  Paciente,
  PacientePayload,
  Sexo,
} from '../../types/paciente'

interface PacienteFormProps {
  paciente?: Paciente | null
  onSuccess: () => void
  onCancel: () => void
}

interface Formulario {
  nome: string
  cpf: string
  dataNascimento: string
  sexo: Sexo
  telefone: string
  endereco: string
  cartaoSUS: string
}

const IDADE_MAXIMA = 150

const formularioInicial: Formulario = {
  nome: '',
  cpf: '',
  dataNascimento: '',
  sexo: 'Masculino',
  telefone: '',
  endereco: '',
  cartaoSUS: '',
}

function somenteNumeros(valor: string): string {
  return valor.replace(/\D/g, '')
}

function validarCPF(cpf: string): boolean {
  const numeros = somenteNumeros(cpf)

  if (numeros.length !== 11) return false
  if (/^(\d)\1{10}$/.test(numeros)) return false

  for (let tamanho = 9; tamanho <= 10; tamanho++) {
    let soma = 0

    for (let i = 0; i < tamanho; i++) {
      soma += Number(numeros[i]) * (tamanho + 1 - i)
    }

    const resto = (soma * 10) % 11
    const digito = resto === 10 ? 0 : resto

    if (digito !== Number(numeros[tamanho])) {
      return false
    }
  }

  return true
}

function formatarDataLocal(data: Date): string {
  const ano = data.getFullYear()
  const mes = String(data.getMonth() + 1).padStart(2, '0')
  const dia = String(data.getDate()).padStart(2, '0')

  return `${ano}-${mes}-${dia}`
}

function obterLimitesNascimento() {
  const hoje = new Date()
  const dataMaxima = formatarDataLocal(hoje)

  const dataMinima = new Date(
    hoje.getFullYear() - IDADE_MAXIMA,
    hoje.getMonth(),
    hoje.getDate()
  )

  return {
    dataMinima: formatarDataLocal(dataMinima),
    dataMaxima,
  }
}

export default function PacienteForm({
  paciente,
  onSuccess,
  onCancel,
}: PacienteFormProps) {
  const editando = Boolean(paciente)

  const [form, setForm] = useState<Formulario>(() => {
    if (!paciente) {
      return { ...formularioInicial }
    }

    return {
      nome: paciente.nome,
      cpf: paciente.cpf,
      dataNascimento: paciente.data_nascimento.slice(0, 10),
      sexo: paciente.sexo,
      telefone: paciente.telefone ?? '',
      endereco: paciente.endereco ?? '',
      cartaoSUS: paciente.cartao_sus ?? '',
    }
  })

  const [salvando, setSalvando] = useState(false)
  const [erro, setErro] = useState('')

  const { dataMinima, dataMaxima } = obterLimitesNascimento()

  function atualizarCampo(
    campo: keyof Formulario,
    valor: string
  ) {
    setForm((anterior) => ({
      ...anterior,
      [campo]: valor,
    }))
  }

  async function salvarPaciente(
    event: FormEvent<HTMLFormElement>
  ) {
    event.preventDefault()
    setErro('')

    if (salvando) return

    if (form.nome.trim().length < 3) {
      setErro('Informe um nome com pelo menos 3 caracteres.')
      return
    }

    if (!validarCPF(form.cpf)) {
      setErro('Informe um CPF válido.')
      return
    }

    if (!form.dataNascimento) {
      setErro('Informe a data de nascimento.')
      return
    }

    const nascimento = new Date(
      `${form.dataNascimento}T00:00:00Z`
    )

    const dataValida =
      !Number.isNaN(nascimento.getTime()) &&
      nascimento.toISOString().slice(0, 10) ===
        form.dataNascimento

    if (
      !dataValida ||
      form.dataNascimento < dataMinima ||
      form.dataNascimento > dataMaxima
    ) {
      setErro(
        `Informe uma data de nascimento válida. ` +
        `A idade máxima permitida é de ${IDADE_MAXIMA} anos.`
      )
      return
    }

    const payload: PacientePayload = {
      nome: form.nome.trim(),
      cpf: somenteNumeros(form.cpf),
      data_nascimento: nascimento.toISOString(),
      sexo: form.sexo,
      telefone: somenteNumeros(form.telefone) || null,
      endereco: form.endereco.trim() || null,
      cartao_sus: somenteNumeros(form.cartaoSUS) || null,
    }

    try {
      setSalvando(true)

      if (paciente) {
        await pacienteService.atualizar(
          paciente.id,
          payload
        )
      } else {
        await pacienteService.cadastrar(payload)
      }

      onSuccess()
    } catch (error) {
      console.error(
        editando
          ? 'Erro ao atualizar paciente:'
          : 'Erro ao cadastrar paciente:',
        error
      )

      if (axios.isAxiosError(error)) {
        const mensagem = error.response?.data?.error

        setErro(
          typeof mensagem === 'string'
            ? mensagem
            : editando
              ? 'Não foi possível atualizar o paciente.'
              : 'Não foi possível cadastrar o paciente.'
        )
      } else {
        setErro('Ocorreu um erro inesperado.')
      }
    } finally {
      setSalvando(false)
    }
  }

  const inputClass =
    'w-full rounded-xl border border-slate-200 ' +
    'bg-slate-50/70 px-4 py-3 text-sm text-slate-800 ' +
    'outline-none transition placeholder:text-slate-400 ' +
    'focus:border-teal-500 focus:bg-white ' +
    'focus:ring-2 focus:ring-teal-500/10 ' +
    'disabled:cursor-not-allowed disabled:opacity-60'

  const labelClass =
    'mb-2 block text-sm font-semibold text-slate-700'

  return (
    <section className="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
      {/* Cabeçalho */}
      <div className="relative overflow-hidden border-b border-slate-100 bg-gradient-to-r from-white via-white to-teal-50/70 p-6 sm:p-7">
        <div className="relative flex items-start justify-between gap-4">
          <div className="flex items-start gap-4">
            <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-teal-600 text-white shadow-lg shadow-teal-100">
              <UserRound size={24} />
            </div>

            <div>
              <div className="mb-2">
                <span className="inline-flex rounded-full border border-teal-100 bg-teal-50 px-2.5 py-1 text-[10px] font-bold uppercase tracking-wider text-teal-700">
                  {editando
                    ? 'Atualização cadastral'
                    : 'Novo cadastro'}
                </span>
              </div>

              <h2 className="text-xl font-bold tracking-tight text-slate-900 sm:text-2xl">
                {editando
                  ? 'Editar paciente'
                  : 'Cadastrar paciente'}
              </h2>

              <p className="mt-2 max-w-xl text-sm leading-relaxed text-slate-500">
                {editando
                  ? 'Atualize os dados cadastrais do paciente mantendo as informações corretas e organizadas.'
                  : 'Preencha as informações abaixo para registrar um novo paciente na unidade.'}
              </p>
            </div>
          </div>

          <button
            type="button"
            onClick={onCancel}
            disabled={salvando}
            aria-label="Fechar formulário"
            className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl border border-slate-200 bg-white text-slate-500 transition hover:border-slate-300 hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-50"
          >
            <X size={18} />
          </button>
        </div>
      </div>

      <form onSubmit={salvarPaciente}>
        <div className="space-y-8 p-5 sm:p-7">
          {/* Orientação */}
          <div className="flex items-start gap-3 rounded-xl border border-blue-100 bg-blue-50/70 p-4">
            <ShieldCheck
              size={19}
              className="mt-0.5 shrink-0 text-blue-600"
            />

            <div>
              <p className="text-sm font-semibold text-blue-900">
                Informações cadastrais
              </p>

              <p className="mt-1 text-xs leading-relaxed text-blue-700">
                Os campos identificados com asterisco (*) são
                obrigatórios. Confira os dados antes de salvar.
              </p>
            </div>
          </div>

          {/* Erro */}
          {erro && (
            <div
              role="alert"
              className="flex items-start gap-3 rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-700"
            >
              <AlertCircle
                size={19}
                className="shrink-0"
              />

              <span>{erro}</span>
            </div>
          )}

          {/* Seção 1 */}
          <section className="space-y-5">
            <div className="flex items-center gap-3 border-b border-slate-100 pb-4">
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-teal-50 text-teal-600">
                <UserRound size={20} />
              </div>

              <div>
                <h3 className="text-sm font-bold text-slate-900">
                  Dados pessoais
                </h3>

                <p className="mt-0.5 text-xs text-slate-500">
                  Identificação e informações básicas
                </p>
              </div>
            </div>

            <div className="grid gap-5 md:grid-cols-2">
              <div className="md:col-span-2">
                <label
                  htmlFor="paciente-nome"
                  className={labelClass}
                >
                  Nome completo{' '}
                  <span className="text-red-500">*</span>
                </label>

                <input
                  id="paciente-nome"
                  className={inputClass}
                  value={form.nome}
                  onChange={(e) =>
                    atualizarCampo('nome', e.target.value)
                  }
                  placeholder="Nome completo do paciente"
                  autoComplete="name"
                  maxLength={150}
                  disabled={salvando}
                  required
                />
              </div>

              <div>
                <label
                  htmlFor="paciente-cpf"
                  className={labelClass}
                >
                  CPF{' '}
                  <span className="text-red-500">*</span>
                </label>

                <input
                  id="paciente-cpf"
                  className={inputClass}
                  value={form.cpf}
                  onChange={(e) =>
                    atualizarCampo(
                      'cpf',
                      somenteNumeros(e.target.value).slice(
                        0,
                        11
                      )
                    )
                  }
                  placeholder="Digite os 11 números do CPF"
                  inputMode="numeric"
                  autoComplete="off"
                  maxLength={11}
                  disabled={salvando}
                  required
                />

                <p className="mt-2 text-xs text-slate-400">
                  Informe somente números.
                </p>
              </div>

              <div>
                <label
                  htmlFor="paciente-nascimento"
                  className={labelClass}
                >
                  Data de nascimento{' '}
                  <span className="text-red-500">*</span>
                </label>

                <div className="relative">
                  <input
                    id="paciente-nascimento"
                    type="date"
                    className={inputClass}
                    value={form.dataNascimento}
                    min={dataMinima}
                    max={dataMaxima}
                    onChange={(e) =>
                      atualizarCampo(
                        'dataNascimento',
                        e.target.value
                      )
                    }
                    disabled={salvando}
                    required
                  />
                </div>

                <p className="mt-2 flex items-center gap-1.5 text-xs text-slate-400">
                  <CalendarDays size={13} />
                  Idade máxima permitida: {IDADE_MAXIMA} anos.
                </p>
              </div>

              <div>
                <label
                  htmlFor="paciente-sexo"
                  className={labelClass}
                >
                  Sexo{' '}
                  <span className="text-red-500">*</span>
                </label>

                <select
                  id="paciente-sexo"
                  className={inputClass}
                  value={form.sexo}
                  onChange={(e) =>
                    atualizarCampo(
                      'sexo',
                      e.target.value as Sexo
                    )
                  }
                  disabled={salvando}
                  required
                >
                  <option value="Masculino">
                    Masculino
                  </option>

                  <option value="Feminino">
                    Feminino
                  </option>

                  <option value="Outro">
                    Outro
                  </option>
                </select>
              </div>
            </div>
          </section>

          {/* Seção 2 */}
          <section className="space-y-5">
            <div className="flex items-center gap-3 border-b border-slate-100 pb-4">
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-blue-50 text-blue-600">
                <MapPin size={20} />
              </div>

              <div>
                <h3 className="text-sm font-bold text-slate-900">
                  Contato e endereço
                </h3>

                <p className="mt-0.5 text-xs text-slate-500">
                  Informações para comunicação e localização
                </p>
              </div>
            </div>

            <div className="grid gap-5 md:grid-cols-2">
              <div>
                <label
                  htmlFor="paciente-telefone"
                  className={labelClass}
                >
                  Telefone
                </label>

                <input
                  id="paciente-telefone"
                  className={inputClass}
                  value={form.telefone}
                  onChange={(e) =>
                    atualizarCampo(
                      'telefone',
                      somenteNumeros(
                        e.target.value
                      ).slice(0, 11)
                    )
                  }
                  placeholder="DDD + número"
                  inputMode="tel"
                  autoComplete="tel"
                  maxLength={11}
                  disabled={salvando}
                />

                <p className="mt-2 text-xs text-slate-400">
                  Campo opcional.
                </p>
              </div>

              <div className="md:col-span-2">
                <label
                  htmlFor="paciente-endereco"
                  className={labelClass}
                >
                  Endereço
                </label>

                <input
                  id="paciente-endereco"
                  className={inputClass}
                  value={form.endereco}
                  onChange={(e) =>
                    atualizarCampo(
                      'endereco',
                      e.target.value
                    )
                  }
                  placeholder="Rua, número, bairro e cidade"
                  autoComplete="street-address"
                  disabled={salvando}
                />
              </div>
            </div>
          </section>

          {/* Seção 3 */}
          <section className="space-y-5">
            <div className="flex items-center gap-3 border-b border-slate-100 pb-4">
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-violet-50 text-violet-600">
                <HeartPulse size={20} />
              </div>

              <div>
                <h3 className="text-sm font-bold text-slate-900">
                  Identificação no SUS
                </h3>

                <p className="mt-0.5 text-xs text-slate-500">
                  Informações complementares de saúde
                </p>
              </div>
            </div>

            <div className="grid gap-5 md:grid-cols-2">
              <div>
                <label
                  htmlFor="paciente-sus"
                  className={labelClass}
                >
                  Cartão Nacional de Saúde (CNS)
                </label>

                <input
                  id="paciente-sus"
                  className={inputClass}
                  value={form.cartaoSUS}
                  onChange={(e) =>
                    atualizarCampo(
                      'cartaoSUS',
                      somenteNumeros(
                        e.target.value
                      ).slice(0, 15)
                    )
                  }
                  placeholder="Número do Cartão SUS"
                  inputMode="numeric"
                  maxLength={15}
                  disabled={salvando}
                />

                <p className="mt-2 text-xs text-slate-400">
                  Informe até 15 dígitos, sem pontuação.
                </p>
              </div>
            </div>
          </section>
        </div>

        {/* Rodapé e ações */}
        <div className="flex flex-col gap-4 border-t border-slate-100 bg-slate-50/70 px-5 py-5 sm:flex-row sm:items-center sm:justify-between sm:px-7">
          <div className="flex items-center gap-2 text-xs text-slate-500">
            <ClipboardList
              size={16}
              className="text-teal-600"
            />

            <span>
              {editando
                ? 'Revise as informações antes de atualizar.'
                : 'Confira os dados antes de cadastrar.'}
            </span>
          </div>

          <div className="flex flex-col gap-3 sm:flex-row">
            <button
              type="button"
              onClick={onCancel}
              disabled={salvando}
              className="inline-flex items-center justify-center gap-2 rounded-xl border border-slate-200 bg-white px-5 py-3 text-sm font-semibold text-slate-700 transition hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-50"
            >
              <X size={16} />
              Cancelar
            </button>

            <button
              type="submit"
              disabled={salvando}
              className="inline-flex items-center justify-center gap-2 rounded-xl bg-teal-600 px-6 py-3 text-sm font-semibold text-white shadow-sm shadow-teal-100 transition hover:bg-teal-700 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {salvando ? (
                <Loader2
                  size={17}
                  className="animate-spin"
                />
              ) : (
                <Save size={17} />
              )}

              {salvando
                ? 'Salvando...'
                : editando
                  ? 'Salvar alterações'
                  : 'Cadastrar paciente'}
            </button>
          </div>
        </div>
      </form>
    </section>
  )
}
