
import { useCallback, useEffect, useMemo, useState } from 'react'
import {
  Activity,
  AlertCircle,
  ArrowRight,
  CheckCircle2,
  ClipboardList,
  Clock3,
  HeartPulse,
  Loader2,
  RefreshCw,
  ShieldCheck,
  Stethoscope,
  TrendingUp,
  Users,
} from 'lucide-react'
import { Link } from 'react-router-dom'

import { pacienteService } from '../services/pacienteService'
import { triagemService } from '../services/triagemService'
import { filaService } from '../services/filaService'
import { atendimentoService } from '../services/atendimentoService'

import type { Paciente } from '../types/paciente'
import type { Triagem } from '../types/triagem'
import type { Fila } from '../types/fila'
import type { Atendimento } from '../types/atendimento'

interface DashboardData {
  pacientes: Paciente[]
  triagens: Triagem[]
  filas: Fila[]
  atendimentos: Atendimento[]
}

const dadosIniciais: DashboardData = {
  pacientes: [],
  triagens: [],
  filas: [],
  atendimentos: [],
}

const nomesPacientesTeste = new Set([
  'paciente teste',
  'paciente teste 2',
])

function ehPacienteTeste(nome: string): boolean {
  return nomesPacientesTeste.has(
    nome.trim().toLocaleLowerCase('pt-BR')
  )
}

function formatarData(valor?: string | null): string {
  if (!valor) return '—'

  const data = new Date(valor)

  return Number.isNaN(data.getTime())
    ? '—'
    : data.toLocaleString('pt-BR')
}

function obterMensagemErro(error: unknown): string {
  if (
    typeof error === 'object' &&
    error !== null &&
    'response' in error
  ) {
    const response = (
      error as {
        response?: {
          data?: {
            error?: string
            message?: string
          }
        }
      }
    ).response

    return (
      response?.data?.error ??
      response?.data?.message ??
      'Erro ao consultar os dados do dashboard.'
    )
  }

  return error instanceof Error
    ? error.message
    : 'Não foi possível carregar o dashboard.'
}

export default function Dashboard() {
  const [dados, setDados] = useState<DashboardData>(dadosIniciais)
  const [carregando, setCarregando] = useState(true)
  const [erro, setErro] = useState('')
  const [atualizadoEm, setAtualizadoEm] = useState<Date | null>(
    null
  )

  const carregarDados = useCallback(async () => {
    setCarregando(true)
    setErro('')

    try {
      const [
        pacientes,
        triagens,
        filas,
        atendimentos,
      ] = await Promise.all([
        pacienteService.listar(),
        triagemService.listar(),
        filaService.listar(),
        atendimentoService.listar(),
      ])

      setDados({
        pacientes: pacientes ?? [],
        triagens: triagens ?? [],
        filas: filas ?? [],
        atendimentos: atendimentos ?? [],
      })

      setAtualizadoEm(new Date())
    } catch (error) {
      setErro(obterMensagemErro(error))
    } finally {
      setCarregando(false)
    }
  }, [])

  useEffect(() => {
    void carregarDados()
  }, [carregarDados])

  // Oculta os pacientes de teste apenas da interface.
  // Nenhum registro é excluído do banco.
  const dadosVisiveis = useMemo(() => {
    const pacientes = dados.pacientes.filter(
      (paciente) => !ehPacienteTeste(paciente.nome)
    )

    const idsPacientesVisiveis = new Set(
      pacientes.map((paciente) => paciente.id)
    )

    return {
      pacientes,
      triagens: dados.triagens.filter((triagem) =>
        idsPacientesVisiveis.has(triagem.id_paciente)
      ),
      filas: dados.filas.filter((fila) =>
        idsPacientesVisiveis.has(fila.id_paciente)
      ),
      atendimentos: dados.atendimentos.filter(
        (atendimento) =>
          idsPacientesVisiveis.has(atendimento.id_paciente)
      ),
    }
  }, [dados])

  const indicadores = useMemo(
    () => [
      {
        titulo: 'Pacientes cadastrados',
        valor: dadosVisiveis.pacientes.length,
        descricao: 'Registros de pacientes',
        icone: Users,
        corIcone: 'text-blue-600',
        corFundo: 'bg-blue-50',
        corBarra: 'bg-blue-500',
        rota: '/pacientes',
      },
      {
        titulo: 'Triagens realizadas',
        valor: dadosVisiveis.triagens.length,
        descricao: 'Avaliações registradas',
        icone: ClipboardList,
        corIcone: 'text-violet-600',
        corFundo: 'bg-violet-50',
        corBarra: 'bg-violet-500',
        rota: '/triagens',
      },
      {
        titulo: 'Aguardando atendimento',
        valor: dadosVisiveis.filas.filter(
          (fila) => fila.status === 'Aguardando'
        ).length,
        descricao: 'Pacientes aguardando',
        icone: Clock3,
        corIcone: 'text-amber-600',
        corFundo: 'bg-amber-50',
        corBarra: 'bg-amber-500',
        rota: '/filas',
      },
      {
        titulo: 'Em atendimento',
        valor: dadosVisiveis.atendimentos.filter(
          (atendimento) => !atendimento.data_fim
        ).length,
        descricao: 'Atendimentos em andamento',
        icone: Stethoscope,
        corIcone: 'text-teal-600',
        corFundo: 'bg-teal-50',
        corBarra: 'bg-teal-500',
        rota: '/atendimentos',
      },
    ],
    [dadosVisiveis]
  )

  const riscos = useMemo(() => {
    const classificacoes = [
      {
        nome: 'Vermelho',
        cor: 'bg-red-500',
        texto: 'text-red-600',
        fundo: 'bg-red-50',
      },
      {
        nome: 'Amarelo',
        cor: 'bg-amber-400',
        texto: 'text-amber-600',
        fundo: 'bg-amber-50',
      },
      {
        nome: 'Verde',
        cor: 'bg-emerald-500',
        texto: 'text-emerald-600',
        fundo: 'bg-emerald-50',
      },
    ]

    return classificacoes.map((item) => ({
      ...item,
      quantidade: dadosVisiveis.triagens.filter(
        (triagem) =>
          triagem.classificacao_risco === item.nome
      ).length,
    }))
  }, [dadosVisiveis.triagens])

  const totalTriagensClassificadas = riscos.reduce(
    (total, item) => total + item.quantidade,
    0
  )

  const resumoFila = useMemo(
    () => [
      {
        nome: 'Aguardando',
        quantidade: dadosVisiveis.filas.filter(
          (fila) => fila.status === 'Aguardando'
        ).length,
        cor: 'bg-amber-400',
      },
      {
        nome: 'Chamado',
        quantidade: dadosVisiveis.filas.filter(
          (fila) => fila.status === 'Chamado'
        ).length,
        cor: 'bg-blue-500',
      },
      {
        nome: 'Em Atendimento',
        quantidade: dadosVisiveis.filas.filter(
          (fila) => fila.status === 'Em Atendimento'
        ).length,
        cor: 'bg-teal-500',
      },
      {
        nome: 'Finalizado',
        quantidade: dadosVisiveis.filas.filter(
          (fila) => fila.status === 'Finalizado'
        ).length,
        cor: 'bg-emerald-500',
      },
      {
        nome: 'Cancelado',
        quantidade: dadosVisiveis.filas.filter(
          (fila) => fila.status === 'Cancelado'
        ).length,
        cor: 'bg-slate-400',
      },
    ],
    [dadosVisiveis.filas]
  )

  const totalFilas = dadosVisiveis.filas.length

  const pacientesPorID = useMemo(
    () =>
      new Map(
        dadosVisiveis.pacientes.map((paciente) => [
          paciente.id,
          paciente.nome,
        ])
      ),
    [dadosVisiveis.pacientes]
  )

  const atendimentosRecentes = useMemo(
    () =>
      [...dadosVisiveis.atendimentos]
        .sort(
          (a, b) =>
            new Date(b.data_inicio).getTime() -
            new Date(a.data_inicio).getTime()
        )
        .slice(0, 5),
    [dadosVisiveis.atendimentos]
  )

  const totalFinalizados = dadosVisiveis.atendimentos.filter(
    (atendimento) => Boolean(atendimento.data_fim)
  ).length

  const totalChamados = dadosVisiveis.filas.filter(
    (fila) => fila.status === 'Chamado'
  ).length

  const dadosDisponiveis = atualizadoEm !== null

  return (
    <div className="min-w-0 w-full max-w-full space-y-5 pb-4 sm:space-y-7">

      {/* Apresentação do dashboard */}
      <section
        className="
          relative min-w-0 max-w-full overflow-hidden
          rounded-2xl border border-slate-200
          bg-white p-3 shadow-sm
          min-[360px]:p-4 sm:p-7
        "
      >
        <div
          aria-hidden="true"
          className="
            pointer-events-none absolute
            -right-16 -top-20 h-56 w-56
            rounded-full bg-teal-50 blur-3xl
          "
        />

        <div
          className="
            relative flex min-w-0 flex-col
            justify-between gap-5
            md:flex-row md:items-center
          "
        >
          <div
            className="
              flex min-w-0 flex-col gap-3
              min-[360px]:flex-row min-[360px]:items-start
              min-[360px]:gap-4
            "
          >
            <div
              className="
                flex h-10 w-10 shrink-0
                items-center justify-center
                rounded-xl bg-teal-600 text-white
                shadow-lg shadow-teal-100
                sm:h-12 sm:w-12 sm:rounded-2xl
              "
            >
              <HeartPulse size={23} aria-hidden="true" />
            </div>

            <div className="min-w-0 flex-1">
              <div
                className="
                  mb-2 flex min-w-0 flex-wrap
                  items-center gap-2
                "
              >
                <span
                  className="
                    max-w-full rounded-full
                    border border-teal-100
                    bg-teal-50 px-2 py-1
                    text-[9px] font-bold uppercase
                    tracking-wide text-teal-700
                    sm:px-2.5 sm:text-[10px]
                  "
                >
                  Visão operacional
                </span>

                <span
                  className="
                    inline-flex min-w-0 items-center
                    gap-1.5 text-[10px]
                    font-medium text-slate-500
                    sm:text-xs
                  "
                >
                  <span
                    className="
                      h-2 w-2 shrink-0
                      rounded-full bg-teal-500
                    "
                  />

                  <span className="break-words">
                    Plataforma de gestão
                  </span>
                </span>
              </div>

              <h2
                className="
                  max-w-full break-words
                  text-lg font-bold leading-tight
                  tracking-tight text-slate-900
                  min-[360px]:text-xl
                  sm:text-3xl
                "
              >
                Visão geral da unidade
              </h2>

              <p
                className="
                  mt-2 max-w-xl break-words
                  text-xs leading-relaxed
                  text-slate-500 sm:text-sm
                "
              >
                Acompanhe os indicadores, a classificação de risco
                e o fluxo de atendimentos em um só lugar.
              </p>

              {atualizadoEm && (
                <p
                  className="
                    mt-3 flex min-w-0
                    flex-wrap items-center gap-1.5
                    break-words text-[10px]
                    leading-relaxed text-slate-400
                    sm:text-xs
                  "
                >
                  <Clock3
                    size={13}
                    className="shrink-0"
                    aria-hidden="true"
                  />

                  <span>
                    Última atualização:{' '}
                    {atualizadoEm.toLocaleString('pt-BR')}
                  </span>
                </p>
              )}
            </div>
          </div>

          <button
            type="button"
            onClick={() => void carregarDados()}
            disabled={carregando}
            className="
              inline-flex max-w-full shrink-0
              items-center justify-center gap-2
              self-start rounded-xl
              border border-slate-200 bg-white
              px-3 py-2.5 text-xs font-semibold
              text-slate-700 shadow-sm transition
              hover:border-teal-200 hover:bg-teal-50
              hover:text-teal-700
              disabled:cursor-not-allowed
              disabled:opacity-50
              sm:px-4 sm:text-sm
            "
          >
            <RefreshCw
              size={16}
              className={carregando ? 'animate-spin' : ''}
              aria-hidden="true"
            />

            Atualizar dados
          </button>
        </div>
      </section>

      {/* Mensagem de erro */}
      {erro && (
        <div
          role="alert"
          className="
            flex min-w-0 items-start gap-3
            rounded-xl border border-red-200
            bg-red-50 p-3 text-xs
            text-red-700 sm:p-4 sm:text-sm
          "
        >
          <AlertCircle
            size={20}
            className="shrink-0"
            aria-hidden="true"
          />

          <span className="min-w-0 break-words">
            {erro}
          </span>
        </div>
      )}

      {/* Indicador de carregamento */}
      {carregando && (
        <div
          role="status"
          className="
            flex min-w-0 items-center gap-2
            text-xs text-slate-500 sm:text-sm
          "
        >
          <Loader2
            size={17}
            className="shrink-0 animate-spin"
            aria-hidden="true"
          />

          Carregando indicadores...
        </div>
      )}

      {/* Indicadores principais */}
      <section className="min-w-0">
        <div
          className="
            mb-4 flex min-w-0 flex-wrap
            items-start justify-between gap-3
          "
        >
          <div className="min-w-0">
            <h3 className="text-base font-bold text-slate-900">
              Indicadores principais
            </h3>

            <p
              className="
                mt-1 break-words text-xs
                text-slate-500
              "
            >
              Panorama dos registros e atendimentos da unidade
            </p>
          </div>

          <span
            className="
              inline-flex max-w-full items-center gap-1.5
              rounded-full border border-slate-200
              bg-white px-3 py-1.5
              text-[11px] font-medium text-slate-500
              sm:text-xs
            "
          >
            <TrendingUp
              size={14}
              className="shrink-0 text-teal-600"
              aria-hidden="true"
            />

            Dados consolidados
          </span>
        </div>

        <div
          className="
            grid min-w-0 grid-cols-1
            gap-3 sm:grid-cols-2
            xl:grid-cols-4
          "
        >
          {indicadores.map((indicador) => {
            const Icone = indicador.icone

            return (
              <Link
                key={indicador.titulo}
                to={indicador.rota}
                className="
                  group relative min-w-0 max-w-full
                  overflow-hidden rounded-2xl
                  border border-slate-200 bg-white
                  p-4 shadow-sm
                  transition-all duration-200
                  hover:-translate-y-0.5
                  hover:border-teal-200
                  hover:shadow-lg
                  hover:shadow-slate-200/60
                  sm:p-5
                "
              >
                <div
                  className={`
                    absolute inset-x-0 top-0 h-1
                    ${indicador.corBarra}
                  `}
                />

                <div
                  className="
                    flex items-start
                    justify-between gap-3
                  "
                >
                  <div
                    className={`
                      flex h-11 w-11 shrink-0
                      items-center justify-center
                      rounded-2xl sm:h-12 sm:w-12
                      ${indicador.corFundo}
                      ${indicador.corIcone}
                    `}
                  >
                    <Icone
                      size={23}
                      aria-hidden="true"
                    />
                  </div>

                  <ArrowRight
                    size={17}
                    className="
                      shrink-0 text-slate-300
                      transition
                      group-hover:translate-x-1
                      group-hover:text-teal-600
                    "
                    aria-hidden="true"
                  />
                </div>

                <p
                  className="
                    mt-4 break-words
                    text-sm font-medium
                    text-slate-500 sm:mt-5
                  "
                >
                  {indicador.titulo}
                </p>

                <p
                  className="
                    mt-2 break-all
                    text-3xl font-bold
                    tracking-tight text-slate-900
                    sm:text-4xl
                  "
                >
                  {dadosDisponiveis
                    ? indicador.valor
                    : '—'}
                </p>

                <p
                  className="
                    mt-2 break-words
                    text-xs text-slate-400
                  "
                >
                  {indicador.descricao}
                </p>
              </Link>
            )
          })}
        </div>
      </section>

      {/* Análises operacionais */}
      <section
        className="
          grid min-w-0 grid-cols-1
          gap-4 xl:grid-cols-2
        "
      >
        {/* Classificação de risco */}
        <div
          className="
            min-w-0 max-w-full rounded-2xl
            border border-slate-200 bg-white
            p-3 shadow-sm
            min-[360px]:p-4 sm:p-6
          "
        >
          <div
            className="
              mb-6 flex min-w-0 items-start
              justify-between gap-2 sm:mb-7
            "
          >
            <div className="min-w-0 flex-1">
              <div
                className="
                  mb-2 flex min-w-0
                  items-center gap-2
                "
              >
                <span
                  className="
                    h-2 w-2 shrink-0
                    rounded-full bg-teal-500
                  "
                />

                <span
                  className="
                    break-words text-[9px]
                    font-bold uppercase
                    tracking-wider text-teal-700
                    sm:text-[10px]
                  "
                >
                  Triagem clínica
                </span>
              </div>

              <h3
                className="
                  break-words text-sm
                  font-bold text-slate-900
                  sm:text-base
                "
              >
                Classificação de risco
              </h3>

              <p
                className="
                  mt-1 break-words
                  text-xs text-slate-500
                "
              >
                Distribuição das avaliações realizadas
              </p>
            </div>

            <div
              className="
                shrink-0 rounded-xl
                bg-slate-50 p-2
                text-slate-500 sm:p-2.5
              "
            >
              <ClipboardList
                size={20}
                aria-hidden="true"
              />
            </div>
          </div>

          <div className="min-w-0 space-y-6">
            {riscos.map((risco) => {
              const percentual =
                totalTriagensClassificadas > 0
                  ? (risco.quantidade /
                      totalTriagensClassificadas) * 100
                  : 0

              return (
                <div
                  key={risco.nome}
                  className="min-w-0"
                >
                  <div
                    className="
                      mb-2.5 flex min-w-0
                      flex-wrap items-center
                      justify-between gap-2
                    "
                  >
                    <div
                      className="
                        flex min-w-0 items-center gap-2
                      "
                    >
                      <span
                        className={`
                          h-2.5 w-2.5 shrink-0
                          rounded-full ${risco.cor}
                        `}
                      />

                      <span
                        className={`
                          text-sm font-semibold
                          ${risco.texto}
                        `}
                      >
                        {risco.nome}
                      </span>
                    </div>

                    <div
                      className="
                        shrink-0 text-right text-sm
                      "
                    >
                      <span
                        className="
                          font-bold text-slate-800
                        "
                      >
                        {risco.quantidade}
                      </span>

                      <span
                        className="
                          ml-2 text-xs text-slate-400
                        "
                      >
                        {percentual.toFixed(1)}%
                      </span>
                    </div>
                  </div>

                  <div
                    className={`
                      h-2.5 w-full overflow-hidden
                      rounded-full ${risco.fundo}
                    `}
                    role="progressbar"
                    aria-label={`Risco ${risco.nome}`}
                    aria-valuenow={risco.quantidade}
                    aria-valuemin={0}
                    aria-valuemax={totalTriagensClassificadas}
                  >
                    <div
                      className={`
                        h-full rounded-full ${risco.cor}
                      `}
                      style={{
                        width: `${percentual}%`,
                      }}
                    />
                  </div>
                </div>
              )
            })}
          </div>

          <div
            className="
              mt-7 flex min-w-0 flex-wrap
              items-center justify-between gap-3
              border-t border-slate-100 pt-5
            "
          >
            <span
              className="
                min-w-0 flex-1 break-words
                text-xs text-slate-500
              "
            >
              Total de triagens classificadas
            </span>

            <span
              className="
                shrink-0 rounded-lg
                bg-slate-100 px-3 py-1.5
                text-sm font-bold text-slate-800
              "
            >
              {totalTriagensClassificadas}
            </span>
          </div>
        </div>

        {/* Situação da fila */}
        <div
          className="
            min-w-0 max-w-full rounded-2xl
            border border-slate-200 bg-white
            p-3 shadow-sm
            min-[360px]:p-4 sm:p-6
          "
        >
          <div
            className="
              mb-6 flex min-w-0 items-start
              justify-between gap-2 sm:mb-7
            "
          >
            <div className="min-w-0 flex-1">
              <div
                className="
                  mb-2 flex min-w-0
                  items-center gap-2
                "
              >
                <span
                  className="
                    h-2 w-2 shrink-0
                    rounded-full bg-teal-500
                  "
                />

                <span
                  className="
                    break-words text-[9px]
                    font-bold uppercase
                    tracking-wider text-teal-700
                    sm:text-[10px]
                  "
                >
                  Fluxo de pacientes
                </span>
              </div>

              <h3
                className="
                  break-words text-sm
                  font-bold text-slate-900
                  sm:text-base
                "
              >
                Situação da fila
              </h3>

              <p
                className="
                  mt-1 break-words
                  text-xs text-slate-500
                "
              >
                Distribuição dos registros por status
              </p>
            </div>

            <div
              className="
                shrink-0 rounded-xl
                bg-slate-50 p-2
                text-slate-500 sm:p-2.5
              "
            >
              <Activity
                size={20}
                aria-hidden="true"
              />
            </div>
          </div>

          <div
            className="
              mb-7 flex h-4 w-full
              min-w-0 overflow-hidden
              rounded-full bg-slate-100
            "
          >
            {resumoFila.map((item) => (
              <div
                key={item.nome}
                className={`
                  h-full min-w-0
                  ${item.cor} transition-all
                `}
                title={`${item.nome}: ${item.quantidade}`}
                style={{
                  width:
                    totalFilas > 0
                      ? `${
                          (item.quantidade / totalFilas) * 100
                        }%`
                      : '0%',
                }}
              />
            ))}
          </div>

          <div className="min-w-0 space-y-4">
            {resumoFila.map((item) => (
              <div
                key={item.nome}
                className="
                  flex min-w-0 items-center
                  justify-between gap-2
                "
              >
                <div
                  className="
                    flex min-w-0 items-center gap-2
                    sm:gap-3
                  "
                >
                  <span
                    className={`
                      h-2.5 w-2.5 shrink-0
                      rounded-full ${item.cor}
                    `}
                  />

                  <span
                    className="
                      min-w-0 break-words
                      text-xs text-slate-600
                      sm:text-sm
                    "
                  >
                    {item.nome}
                  </span>
                </div>

                <span
                  className="
                    shrink-0 rounded-lg
                    bg-slate-50 px-2.5 py-1
                    text-sm font-bold text-slate-800
                  "
                >
                  {item.quantidade}
                </span>
              </div>
            ))}
          </div>

          <div
            className="
              mt-7 flex min-w-0 flex-wrap
              items-center justify-between gap-3
              border-t border-slate-100 pt-5
            "
          >
            <span
              className="
                min-w-0 flex-1 break-words
                text-xs text-slate-500
              "
            >
              Total de entradas registradas
            </span>

            <span
              className="
                shrink-0 rounded-lg
                bg-slate-100 px-3 py-1.5
                text-sm font-bold text-slate-800
              "
            >
              {totalFilas}
            </span>
          </div>
        </div>
      </section>

      {/* Atendimentos e resumo */}
      <section
        className="
          grid min-w-0 grid-cols-1
          gap-4 lg:grid-cols-3
        "
      >
        {/* Atendimentos recentes */}
        <div
          className="
            min-w-0 max-w-full rounded-2xl
            border border-slate-200 bg-white
            p-3 shadow-sm
            min-[360px]:p-4 sm:p-6
            lg:col-span-2
          "
        >
          <div
            className="
              mb-6 flex min-w-0 flex-wrap
              items-center justify-between gap-3
            "
          >
            <div className="min-w-0">
              <h3
                className="
                  break-words text-base
                  font-bold text-slate-900
                "
              >
                Atendimentos recentes
              </h3>

              <p
                className="
                  mt-1 break-words
                  text-xs text-slate-500
                "
              >
                Últimos atendimentos iniciados
              </p>
            </div>

            <Link
              to="/atendimentos"
              className="
                inline-flex shrink-0
                items-center gap-1.5
                rounded-lg px-2 py-1.5
                text-xs font-semibold
                text-teal-700 transition
                hover:bg-teal-50 sm:text-sm
              "
            >
              Ver todos

              <ArrowRight
                size={16}
                aria-hidden="true"
              />
            </Link>
          </div>

          {atendimentosRecentes.length === 0 ? (
            <div
              className="
                flex min-w-0 flex-col
                items-center justify-center
                py-10 text-center sm:py-12
              "
            >
              <div
                className="
                  mb-3 rounded-2xl
                  bg-slate-50 p-4 text-slate-400
                "
              >
                <Stethoscope
                  size={25}
                  aria-hidden="true"
                />
              </div>

              <p
                className="
                  break-words text-sm
                  font-semibold text-slate-700
                "
              >
                Nenhum atendimento registrado
              </p>

              <p
                className="
                  mt-1 break-words
                  text-xs text-slate-400
                "
              >
                Os atendimentos aparecerão aqui.
              </p>
            </div>
          ) : (
            <div
              className="
                min-w-0 divide-y
                divide-slate-100
              "
            >
              {atendimentosRecentes.map((atendimento) => {
                const finalizado = Boolean(
                  atendimento.data_fim
                )

                return (
                  <div
                    key={atendimento.id_atendimento}
                    className="
                      flex min-w-0 flex-col
                      items-start gap-3 py-4
                      first:pt-0 last:pb-0
                      sm:flex-row sm:items-center
                      sm:justify-between
                    "
                  >
                    <div
                      className="
                        flex min-w-0 w-full
                        items-start gap-3
                        sm:w-auto sm:flex-1
                      "
                    >
                      <div
                        className="
                          flex h-9 w-9 shrink-0
                          items-center justify-center
                          rounded-xl bg-teal-50
                          text-teal-700 sm:h-10 sm:w-10
                        "
                      >
                        <Users
                          size={18}
                          aria-hidden="true"
                        />
                      </div>

                      <div className="min-w-0 flex-1">
                        <p
                          className="
                            break-words text-sm
                            font-semibold text-slate-800
                          "
                        >
                          {pacientesPorID.get(
                            atendimento.id_paciente
                          ) ?? 'Paciente não identificado'}
                        </p>

                        <p
                          className="
                            mt-1 break-words
                            text-xs leading-relaxed
                            text-slate-500
                          "
                        >
                          Início:{' '}
                          {formatarData(
                            atendimento.data_inicio
                          )}
                        </p>
                      </div>
                    </div>

                    <span
                      className={`
                        inline-flex max-w-full
                        shrink-0 items-center
                        gap-1.5 rounded-full
                        px-3 py-1.5 text-xs
                        font-semibold
                        ${
                          finalizado
                            ? 'bg-emerald-50 text-emerald-700'
                            : 'bg-amber-50 text-amber-700'
                        }
                      `}
                    >
                      {finalizado ? (
                        <CheckCircle2
                          size={13}
                          aria-hidden="true"
                        />
                      ) : (
                        <Clock3
                          size={13}
                          aria-hidden="true"
                        />
                      )}

                      {finalizado
                        ? 'Finalizado'
                        : 'Em atendimento'}
                    </span>
                  </div>
                )
              })}
            </div>
          )}
        </div>

        {/* Resumo operacional */}
        <div
          className="
            flex min-w-0 max-w-full
            flex-col rounded-2xl
            border border-slate-200
            bg-white p-3 shadow-sm
            min-[360px]:p-4 sm:p-6
          "
        >
          <div className="min-w-0">
            <h3
              className="
                break-words text-base
                font-bold text-slate-900
              "
            >
              Resumo operacional
            </h3>

            <p
              className="
                mt-1 break-words
                text-xs text-slate-500
              "
            >
              Acompanhamento da unidade
            </p>
          </div>

          <div
            className="
              mt-6 min-w-0 flex-1
              space-y-6 sm:mt-7
            "
          >
            <div
              className="
                flex min-w-0 items-center gap-3
              "
            >
              <div
                className="
                  shrink-0 rounded-xl
                  bg-blue-50 p-2.5
                  text-blue-600 sm:p-3
                "
              >
                <Clock3
                  size={20}
                  aria-hidden="true"
                />
              </div>

              <div className="min-w-0 flex-1">
                <p
                  className="
                    break-words text-xs
                    text-slate-500
                  "
                >
                  Pacientes chamados
                </p>

                <p
                  className="
                    mt-1 break-all
                    text-2xl font-bold
                    text-slate-900
                  "
                >
                  {totalChamados}
                </p>
              </div>
            </div>

            <div
              className="
                flex min-w-0 items-center gap-3
              "
            >
              <div
                className="
                  shrink-0 rounded-xl
                  bg-emerald-50 p-2.5
                  text-emerald-600 sm:p-3
                "
              >
                <CheckCircle2
                  size={20}
                  aria-hidden="true"
                />
              </div>

              <div className="min-w-0 flex-1">
                <p
                  className="
                    break-words text-xs
                    text-slate-500
                  "
                >
                  Atendimentos finalizados
                </p>

                <p
                  className="
                    mt-1 break-all
                    text-2xl font-bold
                    text-slate-900
                  "
                >
                  {totalFinalizados}
                </p>
              </div>
            </div>

            <div
              className="
                flex min-w-0 items-center gap-3
              "
            >
              <div
                className="
                  shrink-0 rounded-xl
                  bg-violet-50 p-2.5
                  text-violet-600 sm:p-3
                "
              >
                <ClipboardList
                  size={20}
                  aria-hidden="true"
                />
              </div>

              <div className="min-w-0 flex-1">
                <p
                  className="
                    break-words text-xs
                    text-slate-500
                  "
                >
                  Total de triagens
                </p>

                <p
                  className="
                    mt-1 break-all
                    text-2xl font-bold
                    text-slate-900
                  "
                >
                  {dadosVisiveis.triagens.length}
                </p>
              </div>
            </div>
          </div>

          <div
            className="
              mt-7 min-w-0 space-y-3
              border-t border-slate-100 pt-5
            "
          >
            <div
              className="
                flex min-w-0 items-start
                gap-2 text-xs text-slate-500
              "
            >
              <ShieldCheck
                size={16}
                className="shrink-0 text-teal-600"
                aria-hidden="true"
              />

              <span className="min-w-0 break-words">
                Gestão integrada de atendimentos
              </span>
            </div>

            <Link
              to="/filas"
              className="
                flex min-w-0 w-full
                items-center justify-center gap-2
                rounded-xl bg-teal-600
                px-3 py-3 text-center
                text-xs font-semibold
                text-white shadow-sm
                shadow-teal-100 transition
                hover:bg-teal-700
                sm:px-4 sm:text-sm
              "
            >
              Acessar fila

              <ArrowRight
                size={16}
                className="shrink-0"
                aria-hidden="true"
              />
            </Link>
          </div>
        </div>
      </section>
    </div>
  )
}
