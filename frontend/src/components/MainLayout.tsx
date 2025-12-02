import { Outlet, Link, useNavigate } from 'react-router-dom'
import { Layout, Menu, Button, Dropdown, Space, Avatar } from 'antd'
import { UserOutlined, HomeOutlined, ShoppingOutlined, BarChartOutlined, SettingOutlined, LogoutOutlined, CrownOutlined } from '@ant-design/icons'
import { useAuthStore } from '../store/auth'

const { Header, Content, Footer } = Layout

export default function MainLayout() {
  const navigate = useNavigate()
  const { user, isAuthenticated, logout } = useAuthStore()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const userMenu = {
    items: [
      { key: 'dashboard', label: <Link to="/user">我的订阅</Link>, icon: <CrownOutlined /> },
      { key: 'orders', label: <Link to="/user/orders">订单记录</Link>, icon: <ShoppingOutlined /> },
      { key: 'usage', label: <Link to="/user/usage">使用统计</Link>, icon: <BarChartOutlined /> },
      { key: 'settings', label: <Link to="/user/settings">账户设置</Link>, icon: <SettingOutlined /> },
      { type: 'divider' as const },
      ...(user?.role >= 10 ? [{ key: 'admin', label: <Link to="/admin">管理后台</Link>, icon: <SettingOutlined /> }] : []),
      { key: 'logout', label: '退出登录', icon: <LogoutOutlined />, onClick: handleLogout },
    ],
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', background: '#fff', borderBottom: '1px solid #f0f0f0' }}>
        <div style={{ display: 'flex', alignItems: 'center' }}>
          <Link to="/" style={{ fontSize: 20, fontWeight: 'bold', color: '#1890ff', marginRight: 40 }}>
            订阅中心
          </Link>
          <Menu
            mode="horizontal"
            selectedKeys={[]}
            style={{ border: 'none' }}
            items={[
              { key: 'home', label: <Link to="/">套餐列表</Link>, icon: <HomeOutlined /> },
            ]}
          />
        </div>
        <Space>
          {isAuthenticated ? (
            <Dropdown menu={userMenu}>
              <Space style={{ cursor: 'pointer' }}>
                <Avatar icon={<UserOutlined />} />
                <span>{user?.username}</span>
              </Space>
            </Dropdown>
          ) : (
            <Button type="primary" onClick={() => navigate('/login')}>
              登录
            </Button>
          )}
        </Space>
      </Header>
      <Content style={{ padding: '24px 50px', background: '#f5f5f5' }}>
        <Outlet />
      </Content>
      <Footer style={{ textAlign: 'center', background: '#fff' }}>
        订阅中心 ©{new Date().getFullYear()}
      </Footer>
    </Layout>
  )
}
