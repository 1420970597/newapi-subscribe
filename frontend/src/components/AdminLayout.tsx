import { Outlet, Link, useLocation, useNavigate } from 'react-router-dom'
import { Layout, Menu } from 'antd'
import { UserOutlined, AppstoreOutlined, ShoppingCartOutlined, SettingOutlined, HomeOutlined } from '@ant-design/icons'

const { Header, Content, Sider } = Layout

export default function AdminLayout() {
  const location = useLocation()
  const navigate = useNavigate()

  const menuItems = [
    { key: '/admin/users', label: '用户管理', icon: <UserOutlined /> },
    { key: '/admin/plans', label: '套餐管理', icon: <AppstoreOutlined /> },
    { key: '/admin/orders', label: '订单管理', icon: <ShoppingCartOutlined /> },
    { key: '/admin/settings', label: '系统设置', icon: <SettingOutlined /> },
  ]

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider theme="light" width={200}>
        <div style={{ height: 64, display: 'flex', alignItems: 'center', justifyContent: 'center', borderBottom: '1px solid #f0f0f0' }}>
          <Link to="/admin" style={{ fontSize: 18, fontWeight: 'bold', color: '#1890ff' }}>
            管理后台
          </Link>
        </div>
        <Menu
          mode="inline"
          selectedKeys={[location.pathname]}
          style={{ height: 'calc(100% - 64px)', borderRight: 0 }}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header style={{ background: '#fff', padding: '0 24px', display: 'flex', alignItems: 'center', justifyContent: 'space-between', borderBottom: '1px solid #f0f0f0' }}>
          <span style={{ fontSize: 16 }}>
            {menuItems.find(item => item.key === location.pathname)?.label || '管理后台'}
          </span>
          <Link to="/">
            <HomeOutlined /> 返回前台
          </Link>
        </Header>
        <Content style={{ margin: 24, padding: 24, background: '#fff', minHeight: 280 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
