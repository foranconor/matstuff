import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Routes, Route, NavLink, Navigate } from 'react-router-dom'
import './index.css'
import { getToken, logout } from './auth'
import Login from './pages/Login.jsx'
import Materials from './pages/Materials.jsx'
import MaterialDetail from './pages/MaterialDetail.jsx'
import Suppliers from './pages/Suppliers.jsx'
import SupplierDetail from './pages/SupplierDetail.jsx'
import Contacts from './pages/Contacts.jsx'
import ContactDetail from './pages/ContactDetail.jsx'
import Brackets from './pages/Brackets.jsx'
import BracketDetail from './pages/BracketDetail.jsx'
import Profiles from './pages/Profiles.jsx'
import ProfileDetail from './pages/ProfileDetail.jsx'

function ProtectedRoute({ children }) {
  const token = getToken()
  if (!token) {
    return <Navigate to="/login" replace />
  }
  return children
}

function Layout({ children }) {
  function handleLogout() {
    logout()
    window.location.href = '/login'
  }

  return (
    <>
      <nav>
        <NavLink to="/materials">Materials</NavLink>
        <NavLink to="/suppliers">Suppliers</NavLink>
        <NavLink to="/contacts">Contacts</NavLink>
        <NavLink to="/brackets">Brackets</NavLink>
        <NavLink to="/profiles">Profiles</NavLink>
        <button className="btn-logout" onClick={handleLogout}>Sign out</button>
      </nav>
      <main>{children}</main>
    </>
  )
}

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/*" element={
          <ProtectedRoute>
            <Layout>
              <Routes>
                <Route index element={<Navigate to="/materials" replace />} />
                <Route path="materials" element={<Materials />} />
                <Route path="materials/:id" element={<MaterialDetail />} />
                <Route path="suppliers" element={<Suppliers />} />
                <Route path="suppliers/:id" element={<SupplierDetail />} />
                <Route path="contacts" element={<Contacts />} />
                <Route path="contacts/:id" element={<ContactDetail />} />
                <Route path="brackets" element={<Brackets />} />
                <Route path="brackets/:id" element={<BracketDetail />} />
                <Route path="profiles" element={<Profiles />} />
                <Route path="profiles/:id" element={<ProfileDetail />} />
              </Routes>
            </Layout>
          </ProtectedRoute>
        } />
      </Routes>
    </BrowserRouter>
  </StrictMode>,
)
