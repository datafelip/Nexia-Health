
import { useState, type FormEvent } from 'react'
import {
  Activity,
  AlertCircle,
  ArrowRight,
  BrainCircuit,
  CalendarDays,
  CheckCircle2,
  Database,
  Info,
  Loader2,
  Sparkles,
  TrendingUp,
} from 'lucide-react'

import {
  preverDemanda,
  type PrevisaoDemandaResponse,
} from '../services/aiService'

function obterDataAtual(): string {
  const hoje = new Date()

  const ano = hoje.getFullYear()
  const mes = String(hoje.getMonth() + 1).padStart(2, '0')
  const dia = String(hoje.getDate()).padStart(2, '0')

  return `${ano}-${mes}-${dia}`
}

function formatarData(data: string): string {
  const [ano, mes, dia] = data.split('-').map(Number)

  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: 'long',
    year: 'numeric',
  }).format(new Date(ano, mes - 1, dia))
}

export default function InteligenciaArtificial() {
  const [data, setData] = useState(obterDataAtual)
  const [feriado, setFeriado] = useState(false)
  const [carregando, setCarregando] = useState(false)
  const [erro, setErro] = useState<string | null>(null)

  const [previsao, setPrevisao] =
    useState<PrevisaoDemandaResponse | null>(null)

  async function gerarPrevisao(evento: FormEvent<HTMLFormElement>) {
    evento.preventDefault()

    if (!data) {
      setErro('Selecione uma data para gerar a previsão.')
      return
    }

    setCarregando(true)
    setErro(null)
    setPrevisao(null)

    try {
      const resultado = await preverDemanda({
        data,
        feriado: feriado ? 1 : 0,
      })

      setPrevisao(resultado)
    } catch (error) {
      setErro(
        error instanceof Error
          ? error.message
          : 'Não foi possível consultar o Nexia AI.'
      )
    } finally {
      setCarregando(false)
    }
  }

  return (
    <div className="min-w-0 space-y-6 pb-8">
      {/* Cabeçalho */}
      <section
        className="
          relative overflow-hidden rounded-3xl
          bg-[#0F172A] p-5 text-white
          shadow-xl shadow-slate-200/60
          sm:p-8
        "
      >
        <div
          className="
            pointer-events-none absolute -right-16
            -top-20 h-64 w-64 rounded-full
            bg-teal-400/10 blur-3xl
          "
        />

        <div className="relative z-10">
          <div
            className="
              mb-5 inline-flex items-center gap-2
              rounded-full border border-teal-400/20
              bg-teal-400/10 px-3 py-1.5
              text-xs font-semibold text-teal-300
            "
          >
            <Sparkles size={14} aria-hidden="true" />
            NEXIA AI
          </div>

          <div className="flex flex-wrap items-start justify-between gap-5">
            <div className="max-w-2xl">
              <h1
                className="
                  text-2xl font-bold tracking-tight
                  sm:text-3xl
                "
              >
                Inteligência Operacional
              </h1>

              <p
                className="
                  mt-3 text-sm leading-relaxed
                  text-slate-300 sm:text-base
                "
              >
                Explore previsões experimentais de demanda
                hospitalar utilizando Machine Learning
                para análise de padrões operacionais.
              </p>
            </div>

            <div
              className="
                flex h-14 w-14 shrink-0
                items-center justify-center
                rounded-2xl border border-teal-400/20
                bg-teal-400/10 text-teal-300
              "
            >
              <BrainCircuit size={29} aria-hidden="true" />
            </div>
          </div>

          <div
            className="
              mt-6 flex flex-wrap items-center
              gap-3 border-t border-white/10 pt-5
            "
          >
            <span
              className="
                inline-flex items-center gap-2
                text-xs text-slate-300
              "
            >
              <Activity
                size={15}
                className="text-teal-400"
                aria-hidden="true"
              />
              Previsão de demanda
            </span>

            <span
              className="
                inline-flex items-center gap-2
                text-xs text-slate-300
              "
            >
              <Database
                size={15}
                className="text-teal-400"
                aria-hidden="true"
              />
              Dados sintéticos
            </span>
          </div>
        </div>
      </section>

      {/* Formulário e resultado */}
      <section className="grid min-w-0 grid-cols-1 gap-5 xl:grid-cols-5">
        <div
          className="
            min-w-0 rounded-2xl border
            border-slate-200 bg-white
            p-5 shadow-sm sm:p-6
            xl:col-span-2
          "
        >
          <div className="mb-6">
            <div
              className="
                mb-3 flex h-11 w-11
                items-center justify-center
                rounded-xl bg-teal-50
                text-teal-700
              "
            >
              <CalendarDays size={22} aria-hidden="true" />
            </div>

            <h2 className="text-lg font-bold text-slate-900">
              Nova previsão
            </h2>

            <p className="mt-1 text-sm text-slate-500">
              Informe os parâmetros para consultar
              o modelo de Machine Learning.
            </p>
          </div>

          <form onSubmit={gerarPrevisao} className="space-y-5">
            <div>
              <label
                htmlFor="data-previsao"
                className="
                  mb-2 block text-sm
                  font-semibold text-slate-700
                "
              >
                Data da previsão
              </label>

              <input
                id="data-previsao"
                type="date"
                required
                value={data}
                onChange={(evento) => setData(evento.target.value)}
                disabled={carregando}
                className="
                  w-full min-w-0 rounded-xl
                  border border-slate-200
                  bg-slate-50 px-4 py-3
                  text-sm text-slate-800
                  outline-none transition
                  focus:border-teal-500
                  focus:ring-2 focus:ring-teal-500/10
                  disabled:opacity-60
                "
              />
            </div>

            <label
              className="
                flex cursor-pointer items-start gap-3
                rounded-xl border border-slate-200
                bg-slate-50 p-4
              "
            >
              <input
                type="checkbox"
                checked={feriado}
                onChange={(evento) =>
                  setFeriado(evento.target.checked)
                }
                disabled={carregando}
                className="
                  mt-0.5 h-4 w-4
                  accent-teal-600
                "
              />

              <span className="min-w-0">
                <span
                  className="
                    block text-sm font-semibold
                    text-slate-800
                  "
                >
                  Considerar como feriado
                </span>

                <span
                  className="
                    mt-1 block text-xs
                    leading-relaxed text-slate-500
                  "
                >
                  Informa ao modelo o indicador
                  de feriado para a data selecionada.
                </span>
              </span>
            </label>

            <button
              type="submit"
              disabled={carregando}
              className="
                flex w-full items-center
                justify-center gap-2 rounded-xl
                bg-teal-600 px-4 py-3.5
                text-sm font-bold text-white
                shadow-sm shadow-teal-100
                transition hover:bg-teal-700
                disabled:cursor-not-allowed
                disabled:opacity-60
              "
            >
              {carregando ? (
                <>
                  <Loader2
                    size={18}
                    className="animate-spin"
                    aria-hidden="true"
                  />
                  Gerando previsão...
                </>
              ) : (
                <>
                  <BrainCircuit size={18} aria-hidden="true" />
                  Gerar previsão
                  <ArrowRight size={17} aria-hidden="true" />
                </>
              )}
            </button>
          </form>
        </div>

        {/* Painel de resultados */}
        <div
          className="
            min-w-0 rounded-2xl border
            border-slate-200 bg-white
            p-5 shadow-sm sm:p-6
            xl:col-span-3
          "
        >
          <div className="mb-6 flex flex-wrap items-center justify-between gap-3">
            <div>
              <h2 className="text-lg font-bold text-slate-900">
                Resultado da análise
              </h2>

              <p className="mt-1 text-sm text-slate-500">
                Estimativa gerada pelo Nexia AI
              </p>
            </div>

            {previsao && (
              <span
                className="
                  inline-flex items-center gap-1.5
                  rounded-full bg-emerald-50
                  px-3 py-1.5 text-xs
                  font-semibold text-emerald-700
                "
              >
                <CheckCircle2 size={14} aria-hidden="true" />
                Previsão gerada
              </span>
            )}
          </div>

          {erro && (
            <div
              role="alert"
              className="
                mb-5 flex items-start gap-3
                rounded-xl border border-red-200
                bg-red-50 p-4 text-sm
                text-red-700
              "
            >
              <AlertCircle
                size={19}
                className="shrink-0"
                aria-hidden="true"
              />
              <span>{erro}</span>
            </div>
          )}

          {carregando ? (
            <div
              role="status"
              className="
                flex min-h-64 flex-col
                items-center justify-center
                text-center
              "
            >
              <Loader2
                size={36}
                className="animate-spin text-teal-600"
                aria-hidden="true"
              />

              <p className="mt-4 font-semibold text-slate-800">
                Processando previsão
              </p>

              <p className="mt-1 text-sm text-slate-500">
                Consultando o modelo de Machine Learning...
              </p>
            </div>
          ) : previsao ? (
            <div className="space-y-5">
              <div
                className="
                  rounded-2xl border border-teal-100
                  bg-gradient-to-br
                  from-teal-50 to-white
                  p-6 sm:p-8
                "
              >
                <div className="flex items-center gap-2 text-teal-700">
                  <TrendingUp size={18} aria-hidden="true" />

                  <span
                    className="
                      text-xs font-bold
                      uppercase tracking-wider
                    "
                  >
                    Demanda estimada
                  </span>
                </div>

                <div className="mt-5 flex flex-wrap items-end gap-3">
                  <span
                    className="
                      text-6xl font-bold
                      tracking-tight text-slate-900
                      sm:text-7xl
                    "
                  >
                    {previsao.atendimentos_previstos}
                  </span>

                  <span className="mb-2 text-sm text-slate-500">
                    atendimentos previstos
                  </span>
                </div>

                <p className="mt-4 text-sm text-slate-600">
                  Previsão para{' '}
                  <strong>{formatarData(previsao.data)}</strong>
                </p>
              </div>

              <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                <div
                  className="
                    rounded-xl border border-slate-200
                    bg-slate-50 p-4
                  "
                >
                  <p className="text-xs text-slate-500">
                    Algoritmo utilizado
                  </p>

                  <p
                    className="
                      mt-2 break-words text-sm
                      font-bold text-slate-800
                    "
                  >
                    {previsao.modelo}
                  </p>
                </div>

                <div
                  className="
                    rounded-xl border border-slate-200
                    bg-slate-50 p-4
                  "
                >
                  <p className="text-xs text-slate-500">
                    Origem dos dados
                  </p>

                  <p className="mt-2 text-sm font-bold text-slate-800">
                    {previsao.origem_dados === 'sinteticos'
                      ? 'Dados sintéticos'
                      : previsao.origem_dados}
                  </p>
                </div>
              </div>

              <div
                className="
                  flex items-start gap-3 rounded-xl
                  border border-amber-200
                  bg-amber-50 p-4
                "
              >
                <Info
                  size={19}
                  className="shrink-0 text-amber-700"
                  aria-hidden="true"
                />

                <p className="text-xs leading-relaxed text-amber-900">
                  {previsao.observacao}
                </p>
              </div>
            </div>
          ) : (
            <div
              className="
                flex min-h-64 flex-col
                items-center justify-center
                rounded-2xl border-2
                border-dashed border-slate-200
                bg-slate-50/50 px-5
                text-center
              "
            >
              <div
                className="
                  flex h-16 w-16
                  items-center justify-center
                  rounded-2xl bg-teal-50
                  text-teal-600
                "
              >
                <BrainCircuit size={32} aria-hidden="true" />
              </div>

              <h3 className="mt-5 font-bold text-slate-800">
                Pronto para analisar
              </h3>

              <p
                className="
                  mt-2 max-w-sm text-sm
                  leading-relaxed text-slate-500
                "
              >
                Selecione uma data e clique em
                Gerar previsão para visualizar
                a estimativa de demanda hospitalar.
              </p>
            </div>
          )}
        </div>
      </section>

      {/* Transparência científica */}
      <section
        className="
          flex flex-col gap-4 rounded-2xl
          border border-slate-200 bg-white
          p-5 shadow-sm sm:flex-row
          sm:items-start sm:p-6
        "
      >
        <div
          className="
            flex h-11 w-11 shrink-0
            items-center justify-center
            rounded-xl bg-blue-50
            text-blue-700
          "
        >
          <Info size={22} aria-hidden="true" />
        </div>

        <div>
          <h2 className="text-sm font-bold text-slate-900">
            Sobre o modelo experimental
          </h2>

          <p
            className="
              mt-2 text-sm leading-relaxed
              text-slate-600
            "
          >
            O Nexia AI utiliza o algoritmo Random Forest
            Regressor, treinado com registros sintéticos
            de atendimentos. As previsões demonstram
            a integração entre Inteligência Artificial
            e sistemas de gestão hospitalar, mas não
            representam estimativas validadas com
            dados reais de unidades de saúde.
          </p>
        </div>
      </section>
    </div>
  )
}
