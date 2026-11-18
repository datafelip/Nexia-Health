
import { useEffect, useState, type FormEvent } from 'react'
import axios from 'axios'

import {
  Activity,
  AlertCircle,
  CheckCircle2,
  ClipboardPlus,
  HeartPulse,
  Info,
  Loader2,
  Save,
  ShieldCheck,
  Stethoscope,
  Thermometer,
  UserRound,
  X,
} from 'lucide-react'

import { pacienteService } from '../../services/pacienteService'
import { profissionalService } from '../../services/profissionalService'
import { usuarioService } from '../../services/usuarioService'
import { triagemService } from '../../services/triagemService'

import { filtrarPacientesVisiveis } from '../../utils/pacienteUtils'

import type { Paciente } from '../../types/paciente'
import type { Profissional } from '../../types/profissional'
import type { Usuario } from '../../types/usuario'
import type {
  ClassificacaoRisco,
  TriagemPayload,
} from '../../types/triagem'

interface TriagemFormProps {
  onSuccess: () => void
  onCancel: () => void
}

interface ProfissionalDisponivel extends Profissional {
  nome: string
  perfil: string
}

const classificacoes: {
  valor: ClassificacaoRisco
  titulo: string
  descricao: string
  estiloSelecionado: string
  corIndicador: string
  corTexto: string
}[] = [
  {
    valor: 'Vermelho',
    titulo: 'Vermelho',
    descricao: 'Prioridade máxima',
    estiloSelecionado:
      'border-red-400 bg-red-50 ring-2 ring-red-100',
    corIndicador: 'bg-red-500',
    corTexto: 'text-red-700',
  },
  {
    valor: 'Amarelo',
    titulo: 'Amarelo',
    descricao: 'Prioridade intermediária',
    estiloSelecionado:
      'border-amber-400 bg-amber-50 ring-2 ring-amber-100',
    corIndicador: 'bg-amber-400',
    corTexto: 'text-amber-700',
  },
  {
    valor: 'Verde',
    titulo: 'Verde',
    descricao: 'Menor prioridade',
    estiloSelecionado:
      'border-emerald-400 bg-emerald-50 ring-2 ring-emerald-100',
    corIndicador: 'bg-emerald-500',
    corTexto: 'text-emerald-700',
  },
]

function obterMensagemErro(erro: unknown): string {
  if (axios.isAxiosError(erro)) {
    const mensagem = erro.response?.data?.error

    if (typeof mensagem === 'string') {
      return mensagem
    }

    const mensagemAlternativa = erro.response?.data?.message

    if (typeof mensagemAlternativa === 'string') {
      return mensagemAlternativa
    }

    if (!erro.response) {
      return 'Não foi possível conectar à API. Verifique se o backend está em execução.'
    }
  }

  return 'Ocorreu um erro inesperado. Tente novamente.'
}

export default function TriagemForm({
  onSuccess,
  onCancel,
}: TriagemFormProps) {
  const [pacientes, setPacientes] = useState<Paciente[]>([])
  const [profissionais, setProfissionais] = useState<
    ProfissionalDisponivel[]
  >([])

  const [pacienteId, setPacienteId] = useState('')
  const [profissionalId, setProfissionalId] = useState('')

  const [classificacaoRisco, setClassificacaoRisco] =
    useState<ClassificacaoRisco | null>(null)

  const [pressaoArterial, setPressaoArterial] = useState('')
  const [temperatura, setTemperatura] = useState('')
  const [frequenciaCardiaca, setFrequenciaCardiaca] =
    useState('')
  const [observacoes, setObservacoes] = useState('')

  const [carregandoDados, setCarregandoDados] = useState(true)
  const [salvando, setSalvando] = useState(false)
  const [erro, setErro] = useState<string | null>(null)

  useEffect(() => {
    let ativo = true

    async function carregarDados() {
      try {
        setCarregandoDados(true)
        setErro(null)

        const [
          pacientesResponse,
          profissionaisResponse,
          usuariosResponse,
        ] = await Promise.all([
          pacienteService.listar(),
          profissionalService.listar(),
          usuarioService.listar(),
        ])

        if (!ativo) return

        const pacientesVisiveis = filtrarPacientesVisiveis(
          pacientesResponse ?? []
        )

        const usuariosPorId = new Map<string, Usuario>()

        for (const usuario of usuariosResponse ?? []) {
          usuariosPorId.set(usuario.id_usuario, usuario)
        }

        const profissionaisDisponiveis: ProfissionalDisponivel[] =
          []

        for (const profissional of profissionaisResponse ?? []) {
          if (profissional.status !== 'Ativo') {
            continue
          }

          const usuario = usuariosPorId.get(
            profissional.id_usuario
          )

          if (!usuario || usuario.status !== 'Ativo') {
            continue
          }

          if (
            usuario.perfil !== 'Enfermeiro' &&
            usuario.perfil !== 'Médico'
          ) {
            continue
          }

          profissionaisDisponiveis.push({
            ...profissional,
            nome: usuario.nome,
            perfil: usuario.perfil,
          })
        }

        profissionaisDisponiveis.sort((a, b) =>
          a.nome.localeCompare(b.nome, 'pt-BR')
        )

        const pacientesOrdenados = [
          ...pacientesVisiveis,
        ].sort((a, b) =>
          a.nome.localeCompare(b.nome, 'pt-BR')
        )

        setPacientes(pacientesOrdenados)
        setProfissionais(profissionaisDisponiveis)
      } catch (erroCapturado) {
        if (ativo) {
          setErro(obterMensagemErro(erroCapturado))
        }
      } finally {
        if (ativo) {
          setCarregandoDados(false)
        }
      }
    }

    void carregarDados()

    return () => {
      ativo = false
    }
  }, [])

  async function handleSubmit(
    event: FormEvent<HTMLFormElement>
  ) {
    event.preventDefault()

    if (salvando) return

    setErro(null)

    if (!pacienteId) {
      setErro('Selecione um paciente.')
      return
    }

    if (!profissionalId) {
      setErro('Selecione o profissional responsável.')
      return
    }

    if (!classificacaoRisco) {
      setErro('Selecione a classificação de risco.')
      return
    }

    if (temperatura.trim()) {
      const valorTemperatura = Number(temperatura)

      if (
        !Number.isFinite(valorTemperatura) ||
        valorTemperatura < 25 ||
        valorTemperatura > 45
      ) {
        setErro(
          'Informe uma temperatura entre 25 e 45 °C.'
        )
        return
      }
    }

    if (frequenciaCardiaca.trim()) {
      const valorFrequencia = Number(frequenciaCardiaca)

      if (
        !Number.isInteger(valorFrequencia) ||
        valorFrequencia < 1 ||
        valorFrequencia > 32767
      ) {
        setErro(
          'Informe uma frequência cardíaca válida.'
        )
        return
      }
    }

    const pressao = pressaoArterial.trim()

    if (
      pressao &&
      !/^\d{2,3}\/\d{2,3}$/.test(pressao)
    ) {
      setErro(
        'Informe a pressão arterial no formato 120/80.'
      )
      return
    }

    const dados: TriagemPayload = {
      id_paciente: pacienteId,
      id_profissional: profissionalId,
      classificacao_risco: classificacaoRisco,
      pressao_arterial: pressao || null,
      temperatura: temperatura.trim()
        ? Number(temperatura)
        : null,
      frequencia_cardiaca: frequenciaCardiaca.trim()
        ? Number(frequenciaCardiaca)
        : null,
      observacoes: observacoes.trim() || null,
    }

    try {
      setSalvando(true)

      await triagemService.cadastrar(dados)

      onSuccess()
    } catch (erroCapturado) {
      setErro(obterMensagemErro(erroCapturado))
    } finally {
      setSalvando(false)
    }
  }

  const inputClassName =
    'w-full rounded-xl border border-slate-200 ' +
    'bg-slate-50/70 px-4 py-3 text-sm text-slate-800 ' +
    'outline-none transition placeholder:text-slate-400 ' +
    'focus:border-teal-500 focus:bg-white ' +
    'focus:ring-2 focus:ring-teal-500/10 ' +
    'disabled:cursor-not-allowed disabled:opacity-60'

  const labelClassName =
    'mb-2 block text-sm font-semibold text-slate-700'

  return (
    <section className="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
      {/* Cabeçalho */}
      <div className="relative overflow-hidden border-b border-slate-100 bg-gradient-to-r from-white via-white to-teal-50/70 p-6 sm:p-7">
        <div className="relative flex items-start justify-between gap-4">
          <div className="flex items-start gap-4">
            <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-teal-600 text-white shadow-lg shadow-teal-100">
              <ClipboardPlus size={24} />
            </div>

            <div>
              <span className="inline-flex rounded-full border border-teal-100 bg-teal-50 px-2.5 py-1 text-[10px] font-bold uppercase tracking-wider text-teal-700">
                Avaliação inicial
              </span>

              <h2 className="mt-2 text-xl font-bold tracking-tight text-slate-900 sm:text-2xl">
                Nova triagem
              </h2>

              <p className="mt-2 max-w-xl text-sm leading-relaxed text-slate-500">
                Registre a avaliação clínica, defina a
                classificação de risco e encaminhe o paciente
                à fila de atendimento.
              </p>
            </div>
          </div>

          <button
            type="button"
            onClick={onCancel}
            disabled={salvando}
            aria-label="Fechar formulário"
            className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl border border-slate-200 bg-white text-slate-500 transition hover:bg-slate-50 disabled:opacity-50"
          >
            <X size={18} />
          </button>
        </div>
      </div>

      {carregandoDados ? (
        <div
          role="status"
          className="flex flex-col items-center justify-center gap-3 px-6 py-16"
        >
          <Loader2
            size={30}
            className="animate-spin text-teal-600"
          />

          <p className="text-sm text-slate-500">
            Carregando pacientes e profissionais...
          </p>
        </div>
      ) : (
        <form onSubmit={handleSubmit}>
          <div className="space-y-9 p-5 sm:p-7">
            {/* Orientação */}
            <div className="flex items-start gap-3 rounded-xl border border-blue-100 bg-blue-50/70 p-4">
              <ShieldCheck
                size={19}
                className="mt-0.5 shrink-0 text-blue-600"
              />

              <div>
                <p className="text-sm font-semibold text-blue-900">
                  Registro de triagem
                </p>

                <p className="mt-1 text-xs leading-relaxed text-blue-700">
                  Os campos com asterisco (*) são obrigatórios.
                  A classificação deve ser definida pelo
                  profissional responsável conforme o protocolo
                  clínico da unidade.
                </p>
              </div>
            </div>

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

            {/* Identificação */}
            <section className="space-y-5">
              <div className="flex items-center gap-3 border-b border-slate-100 pb-4">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-teal-50 text-teal-600">
                  <UserRound size={20} />
                </div>

                <div>
                  <h3 className="text-sm font-bold text-slate-900">
                    Identificação
                  </h3>

                  <p className="mt-0.5 text-xs text-slate-500">
                    Paciente e profissional responsável
                  </p>
                </div>
              </div>

              <div className="grid gap-5 md:grid-cols-2">
                <div>
                  <label
                    htmlFor="triagem-paciente"
                    className={labelClassName}
                  >
                    Paciente{' '}
                    <span className="text-red-500">*</span>
                  </label>

                  <select
                    id="triagem-paciente"
                    required
                    value={pacienteId}
                    onChange={(event) =>
                      setPacienteId(event.target.value)
                    }
                    disabled={salvando}
                    className={inputClassName}
                  >
                    <option value="">
                      Selecione um paciente
                    </option>

                    {pacientes.map((paciente) => (
                      <option
                        key={paciente.id}
                        value={paciente.id}
                      >
                        {paciente.nome}
                      </option>
                    ))}
                  </select>

                  {pacientes.length === 0 && (
                    <p className="mt-2 text-xs text-amber-700">
                      Nenhum paciente disponível. Cadastre um
                      paciente antes de realizar a triagem.
                    </p>
                  )}
                </div>

                <div>
                  <label
                    htmlFor="triagem-profissional"
                    className={labelClassName}
                  >
                    Profissional responsável{' '}
                    <span className="text-red-500">*</span>
                  </label>

                  <select
                    id="triagem-profissional"
                    required
                    value={profissionalId}
                    onChange={(event) =>
                      setProfissionalId(event.target.value)
                    }
                    disabled={salvando}
                    className={inputClassName}
                  >
                    <option value="">
                      Selecione um profissional
                    </option>

                    {profissionais.map((profissional) => (
                      <option
                        key={profissional.id_profissional}
                        value={profissional.id_profissional}
                      >
                        {profissional.nome}
                        {' — '}
                        {profissional.perfil}
                        {' • '}
                        {profissional.conselho_profissional}
                        {' '}
                        {profissional.numero_registro}
                      </option>
                    ))}
                  </select>

                  {profissionais.length === 0 && (
                    <p className="mt-2 text-xs text-amber-700">
                      Nenhum profissional ativo com perfil
                      Enfermeiro ou Médico foi encontrado.
                    </p>
                  )}
                </div>
              </div>
            </section>

            {/* Sinais vitais */}
            <section className="space-y-5">
              <div className="flex items-center gap-3 border-b border-slate-100 pb-4">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-blue-50 text-blue-600">
                  <Activity size={20} />
                </div>

                <div>
                  <h3 className="text-sm font-bold text-slate-900">
                    Sinais vitais
                  </h3>

                  <p className="mt-0.5 text-xs text-slate-500">
                    Dados clínicos coletados na avaliação
                  </p>
                </div>
              </div>

              <div className="grid gap-5 md:grid-cols-3">
                <div>
                  <label
                    htmlFor="triagem-pressao"
                    className={labelClassName}
                  >
                    Pressão arterial
                  </label>

                  <div className="relative">
                    <input
                      id="triagem-pressao"
                      type="text"
                      placeholder="120/80"
                      value={pressaoArterial}
                      onChange={(event) =>
                        setPressaoArterial(event.target.value)
                      }
                      disabled={salvando}
                      className={inputClassName}
                    />
                  </div>

                  <p className="mt-2 text-xs text-slate-400">
                    mmHg • Exemplo: 120/80
                  </p>
                </div>

                <div>
                  <label
                    htmlFor="triagem-temperatura"
                    className={labelClassName}
                  >
                    Temperatura
                  </label>

                  <div className="relative">
                    <input
                      id="triagem-temperatura"
                      type="number"
                      step="0.1"
                      min="25"
                      max="45"
                      placeholder="36.5"
                      value={temperatura}
                      onChange={(event) =>
                        setTemperatura(event.target.value)
                      }
                      disabled={salvando}
                      className={`${inputClassName} pr-10`}
                    />

                    <Thermometer
                      size={17}
                      className="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-slate-400"
                    />
                  </div>

                  <p className="mt-2 text-xs text-slate-400">
                    Graus Celsius (°C)
                  </p>
                </div>

                <div>
                  <label
                    htmlFor="triagem-frequencia"
                    className={labelClassName}
                  >
                    Frequência cardíaca
                  </label>

                  <div className="relative">
                    <input
                      id="triagem-frequencia"
                      type="number"
                      step="1"
                      min="1"
                      max="32767"
                      placeholder="80"
                      value={frequenciaCardiaca}
                      onChange={(event) =>
                        setFrequenciaCardiaca(
                          event.target.value
                        )
                      }
                      disabled={salvando}
                      className={`${inputClassName} pr-10`}
                    />

                    <HeartPulse
                      size={17}
                      className="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-slate-400"
                    />
                  </div>

                  <p className="mt-2 text-xs text-slate-400">
                    Batimentos por minuto (bpm)
                  </p>
                </div>
              </div>
            </section>

            {/* Classificação de risco */}
            <section className="space-y-5">
              <div className="flex items-center gap-3 border-b border-slate-100 pb-4">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-violet-50 text-violet-600">
                  <HeartPulse size={20} />
                </div>

                <div>
                  <h3 className="text-sm font-bold text-slate-900">
                    Classificação de risco{' '}
                    <span className="text-red-500">*</span>
                  </h3>

                  <p className="mt-0.5 text-xs text-slate-500">
                    Prioridade clínica atribuída pelo profissional
                  </p>
                </div>
              </div>

              <div className="flex items-start gap-2 rounded-xl border border-slate-200 bg-slate-50 p-3">
                <Info
                  size={17}
                  className="mt-0.5 shrink-0 text-slate-500"
                />

                <p className="text-xs leading-relaxed text-slate-600">
                  Selecione a classificação definida pelo
                  profissional responsável, conforme o protocolo
                  clínico adotado pela unidade.
                </p>
              </div>

              <div
                className="grid gap-3 md:grid-cols-3"
                role="group"
                aria-label="Classificação de risco"
              >
                {classificacoes.map((classificacao) => {
                  const selecionado =
                    classificacaoRisco === classificacao.valor

                  return (
                    <button
                      key={classificacao.valor}
                      type="button"
                      aria-pressed={selecionado}
                      disabled={salvando}
                      onClick={() =>
                        setClassificacaoRisco(
                          classificacao.valor
                        )
                      }
                      className={`rounded-2xl border-2 p-5 text-left transition disabled:cursor-not-allowed disabled:opacity-60 ${
                        selecionado
                          ? classificacao.estiloSelecionado
                          : 'border-slate-200 bg-white hover:border-slate-300 hover:bg-slate-50'
                      }`}
                    >
                      <div className="flex items-center justify-between gap-2">
                        <div
                          className={`flex h-10 w-10 items-center justify-center rounded-xl ${
                            selecionado
                              ? 'bg-white/80'
                              : 'bg-slate-50'
                          }`}
                        >
                          <span
                            className={`h-4 w-4 rounded-full ${classificacao.corIndicador}`}
                          />
                        </div>

                        {selecionado && (
                          <CheckCircle2
                            size={20}
                            className={classificacao.corTexto}
                          />
                        )}
                      </div>

                      <p
                        className={`mt-4 text-base font-bold ${
                          selecionado
                            ? classificacao.corTexto
                            : 'text-slate-800'
                        }`}
                      >
                        {classificacao.titulo}
                      </p>

                      <p className="mt-1 text-xs text-slate-500">
                        {classificacao.descricao}
                      </p>
                    </button>
                  )
                })}
              </div>
            </section>

            {/* Observações clínicas */}
            <section className="space-y-5">
              <div className="flex items-center gap-3 border-b border-slate-100 pb-4">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-emerald-50 text-emerald-600">
                  <Stethoscope size={20} />
                </div>

                <div>
                  <h3 className="text-sm font-bold text-slate-900">
                    Observações clínicas
                  </h3>

                  <p className="mt-0.5 text-xs text-slate-500">
                    Informações complementares da avaliação
                  </p>
                </div>
              </div>

              <div>
                <label
                  htmlFor="triagem-observacoes"
                  className={labelClassName}
                >
                  Observações
                </label>

                <textarea
                  id="triagem-observacoes"
                  rows={4}
                  placeholder="Registre as informações relevantes da avaliação..."
                  value={observacoes}
                  onChange={(event) =>
                    setObservacoes(event.target.value)
                  }
                  disabled={salvando}
                  className={`${inputClassName} resize-y`}
                />

                <p className="mt-2 text-xs text-slate-400">
                  Campo opcional para informações clínicas
                  complementares.
                </p>
              </div>
            </section>
          </div>

          {/* Rodapé */}
          <div className="flex flex-col gap-4 border-t border-slate-100 bg-slate-50/70 px-5 py-5 sm:flex-row sm:items-center sm:justify-between sm:px-7">
            <div className="flex items-center gap-2 text-xs text-slate-500">
              <ShieldCheck
                size={16}
                className="shrink-0 text-teal-600"
              />

              <span>
                Ao registrar, o paciente será incluído
                automaticamente na fila.
              </span>
            </div>

            <div className="flex flex-col gap-3 sm:flex-row">
              <button
                type="button"
                onClick={onCancel}
                disabled={salvando}
                className="inline-flex items-center justify-center gap-2 rounded-xl border border-slate-200 bg-white px-5 py-3 text-sm font-semibold text-slate-700 transition hover:bg-slate-100 disabled:opacity-50"
              >
                <X size={16} />
                Cancelar
              </button>

              <button
                type="submit"
                disabled={
                  salvando ||
                  pacientes.length === 0 ||
                  profissionais.length === 0
                }
                className="inline-flex items-center justify-center gap-2 rounded-xl bg-teal-600 px-6 py-3 text-sm font-semibold text-white shadow-sm shadow-teal-100 transition hover:bg-teal-700 disabled:cursor-not-allowed disabled:opacity-50"
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
                  ? 'Registrando triagem...'
                  : 'Registrar triagem'}
              </button>
            </div>
          </div>
        </form>
      )}

      {/* Erro de carregamento inicial */}
      {!carregandoDados && erro && (
        <div className="sr-only" aria-live="polite">
          {erro}
        </div>
      )}
    </section>
  )
}
