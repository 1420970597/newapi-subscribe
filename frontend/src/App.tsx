import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './store/auth'
import MainLayout from './components/MainLayout'
import AdminLayout from './components/AdminLayout'
import Home from './pages/Home'
import Login from './pages/Login'
import Dashboard from './pages/User/Dashboard'
import Orders from './pages/User/Orders'
import Usage from './pages/User/Usage'
import Settings from './pages/User/Settings'
import AdminUsers from './pages/Admin/Users'
import AdminPlans from './pages/Admin/Plans'
import AdminOrders from './pages/Admin/Orders'
import AdminSettings from './pages/Admin/Settings'

function PrivateRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuthStore()
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" />
}

function AdminRoute({ children }: { children: React.ReactNode }) {
  const { user, isAuthenticated } = useAuthStore()
  if (!isAuthenticated) return <Navigate to="/login" />
  if ((user?.role ?? 0) < 10) return <Navigate to="/" />
  return <>{children}</>
}

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />

        <Route path="/" element={<MainLayout />}>
          <Route index element={<Home />} />
          <Route path="user" element={<PrivateRoute><Dashboard /></PrivateRoute>} />
          <Route path="user/orders" element={<PrivateRoute><Orders /></PrivateRoute>} />
          <Route path="user/usage" element={<PrivateRoute><Usage /></PrivateRoute>} />
          <Route path="user/settings" element={<PrivateRoute><Settings /></PrivateRoute>} />
        </Route>

        <Route path="/admin" element={<AdminRoute><AdminLayout /></AdminRoute>}>
          <Route index element={<Navigate to="/admin/users" />} />
          <Route path="users" element={<AdminUsers />} />
          <Route path="plans" element={<AdminPlans />} />
          <Route path="orders" element={<AdminOrders />} />
          <Route path="settings" element={<AdminSettings />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
