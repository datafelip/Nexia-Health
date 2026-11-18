
import { NavLink } from 'react-router-dom'
import {
  Activity,
  BrainCircuit,
  ArrowUpRight,
  ClipboardList,
  HeartPulse,
  LayoutDashboard,
  ListOrdered,
  Stethoscope,
  Users,
} from 'lucide-react'

interface SidebarProps {
  onNavigate?: () => void
}

const menuItems = [
  {
    name: 'Dashboard',
    path: '/',
    icon: LayoutDashboard,
  },
  {
    name: 'Pacientes',
    path: '/pacientes',
    icon: Users,
  },
  {
    name: 'Triagens',
    path: '/triagens',
    icon: ClipboardList,
  },
  {
    name: 'Fila de Atendimento',
    path: '/filas',
    icon: ListOrdered,
  },
  {
    name: 'Atendimentos',
    path: '/atendimentos',
    icon: Stethoscope,
  },
  {
    name: 'Inteligência Artificial',
    path: '/inteligencia-artificial',
    icon: BrainCircuit,
  },
]

export default function Sidebar({
  onNavigate,
}: SidebarProps) {
  return (
    <aside
      className="
        flex h-full w-full min-w-0 flex-col
        overflow-x-hidden overflow-y-auto
        bg-[#0F172A] text-white
        shadow-xl shadow-slate-950/10
        [scrollbar-width:thin]
        [scrollbar-color:#334155_transparent]
      "
    >
      {/* Identidade visual */}
      <div
        className="
          flex min-h-20 shrink-0 items-center
          gap-3 border-b border-white/10
          px-4 sm:px-5
        "
      >
        <div
          className="
            flex h-11 w-11 shrink-0
            items-center justify-center
            rounded-2xl
            bg-gradient-to-br
            from-teal-400 to-teal-600
            shadow-lg shadow-teal-950/30
          "
        >
          <HeartPulse
            size={25}
            strokeWidth={2}
            className="text-white"
            aria-hidden="true"
          />
        </div>

        <div className="min-w-0">
          <p
            className="
              truncate text-lg font-bold
              tracking-tight text-white
            "
          >
            Nexia Health
          </p>

          <p
            className="
              mt-0.5 text-[10px]
              font-medium tracking-wider
              text-slate-400
            "
          >
            MEDICAL PLATFORM
          </p>
        </div>
      </div>

      {/* Navegação principal */}
      <nav
        aria-label="Navegação principal"
        className="min-h-0 flex-1 px-3 py-6"
      >
        <div className="mb-4 px-3">
          <p
            className="
              text-[10px] font-bold
              uppercase tracking-[0.2em]
              text-slate-500
            "
          >
            Menu principal
          </p>
        </div>

        <div className="space-y-1.5">
          {menuItems.map((item) => {
            const Icon = item.icon

            return (
              <NavLink
                key={item.path}
                to={item.path}
                end={item.path === '/'}
                onClick={onNavigate}
                className={({ isActive }) =>
                  [
                    'group relative flex min-h-12',
                    'items-center gap-3',
                    'overflow-hidden rounded-xl',
                    'px-3.5 py-3',
                    'transition-all duration-200',
                    'focus-visible:outline',
                    'focus-visible:outline-2',
                    'focus-visible:outline-offset-2',
                    'focus-visible:outline-teal-400',
                    isActive
                      ? [
                          'bg-teal-500/15',
                          'text-teal-300',
                          'ring-1 ring-inset',
                          'ring-teal-400/20',
                        ].join(' ')
                      : [
                          'text-slate-400',
                          'hover:bg-white/[0.06]',
                          'hover:text-white',
                        ].join(' '),
                  ].join(' ')
                }
              >
                {({ isActive }) => (
                  <>
                    {/* Indicador lateral ativo */}
                    {isActive && (
                      <span
                        aria-hidden="true"
                        className="
                          absolute bottom-2 left-0 top-2
                          w-1 rounded-r-full
                          bg-teal-400
                        "
                      />
                    )}

                    {/* Ícone da página */}
                    <Icon
                      size={19}
                      strokeWidth={isActive ? 2.3 : 1.9}
                      aria-hidden="true"
                      className={
                        isActive
                          ? 'shrink-0 text-teal-400'
                          : [
                              'shrink-0 text-slate-400',
                              'transition-colors',
                              'group-hover:text-teal-300',
                            ].join(' ')
                      }
                    />

                    {/* Nome da página */}
                    <span
                      className="
                        min-w-0 flex-1
                        text-[13px] font-semibold
                        leading-snug
                      "
                    >
                      {item.name}
                    </span>

                    {/* Indicador de página selecionada */}
                    {isActive && (
                      <span
                        aria-hidden="true"
                        className="
                          h-1.5 w-1.5 shrink-0
                          rounded-full bg-teal-400
                        "
                      />
                    )}
                  </>
                )}
              </NavLink>
            )
          })}
        </div>
      </nav>

      {/* Identificação institucional */}
      <div className="shrink-0 px-3 pb-4">
        <div
          className="
            rounded-2xl
            border border-white/10
            bg-white/[0.04]
            p-3.5
          "
        >
          <div
            className="
              mb-3 flex items-center
              justify-between
            "
          >
            <div
              className="
                flex h-9 w-9
                items-center justify-center
                rounded-xl bg-teal-400/10
              "
            >
              <Activity
                size={19}
                className="text-teal-400"
                aria-hidden="true"
              />
            </div>

            <ArrowUpRight
              size={16}
              className="text-slate-500"
              aria-hidden="true"
            />
          </div>

          <p className="text-xs font-semibold text-slate-200">
            Gestão Inteligente de UPA
          </p>

          <p
            className="
              mt-1 text-[11px]
              leading-relaxed
              text-slate-400
            "
          >
            Tecnologia e organização para o cuidado com a saúde.
          </p>
        </div>

        {/* Rodapé da Sidebar */}
        <div
          className="
            mt-4 flex items-center
            justify-between
            border-t border-white/10
            px-2 pt-4
          "
        >
          <span className="text-[11px] font-medium text-slate-400">
            Nexia Health
          </span>

          <span
            className="
              rounded-md
              bg-white/[0.06]
              px-2 py-1
              text-[10px] font-semibold
              text-slate-400
            "
          >
            v1.0
          </span>
        </div>
      </div>
    </aside>
  )
}
