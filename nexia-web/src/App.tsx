
import {
  BrowserRouter,
  Navigate,
  Route,
  Routes,
} from 'react-router-dom'

import MainLayout from './components/layout/MainLayout'

import Dashboard from './pages/Dashboard'
import Pacientes from './pages/Pacientes'
import Triagens from './pages/Triagens'
import Filas from './pages/Filas'
import Atendimentos from './pages/Atendimentos'
import InteligenciaArtificial from './pages/InteligenciaArtificial'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<MainLayout />}>
          {/* Dashboard */}
          <Route
            index
            element={<Dashboard />}
          />

          {/* Pacientes */}
          <Route
            path="pacientes"
            element={<Pacientes />}
          />

          {/* Triagens */}
          <Route
            path="triagens"
            element={<Triagens />}
          />

          {/* Fila de Atendimento */}
          <Route
            path="filas"
            element={<Filas />}
          />

          {/* Atendimentos */}
          <Route
            path="atendimentos"
            element={<Atendimentos />}
          />

          {/* Inteligência Artificial */}
          <Route
            path="inteligencia-artificial"
            element={<InteligenciaArtificial />}
          />
        </Route>

        {/* Redirecionamento para rotas inexistentes */}
        <Route
          path="*"
          element={<Navigate to="/" replace />}
        />
      </Routes>
    </BrowserRouter>
  )
}
