import { useState, useEffect } from 'react'
import { Table, Card, Input, Tag, Button, Drawer, Descriptions, Spin } from 'antd'
import dayjs from 'dayjs'
import { adminApi } from '../../../api'

export default function AdminUsers() {
  const [loading, setLoading] = useState(true)
  const [users, setUsers] = useState<any[]>([])
  const [pagination, setPagination] = useState({ current: 1, pageSize: 20, total: 0 })
  const [keyword, setKeyword] = useState('')
  const [selectedUser, setSelectedUser] = useState<any>(null)
  const [drawerVisible, setDrawerVisible] = useState(false)
  const [userDetail, setUserDetail] = useState<any>(null)
  const [detailLoading, setDetailLoading] = useState(false)

  useEffect(() => {
    loadUsers()
  }, [pagination.current, keyword])

  const loadUsers = async () => {
    setLoading(true)
    try {
      const res: any = await adminApi.getUsers({
        page: pagination.current,
        per_page: pagination.pageSize,
        keyword,
      })
      if (res.success) {
        setUsers(res.data || [])
        setPagination(prev => ({ ...prev, total: res.total }))
      }
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handleViewUser = async (user: any) => {
    setSelectedUser(user)
    setDrawerVisible(true)
    setDetailLoading(true)
    try {
      const res: any = await adminApi.getUser(user.id)
      if (res.success) {
        setUserDetail(res.data)
      }
    } catch (error) {
      console.error(error)
    } finally {
      setDetailLoading(false)
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '用户名', dataIndex: 'username', key: 'username' },
    { title: '邮箱', dataIndex: 'email', key: 'email' },
    {
      title: '订阅状态',
      key: 'subscription',
      render: (_: any, record: any) => record.subscription ? (
        <Tag color="green">{record.subscription.plan?.name}</Tag>
      ) : (
        <Tag>无订阅</Tag>
      ),
    },
    {
      title: 'new-api',
      key: 'newapi',
      render: (_: any, record: any) => record.newapi_bound ? (
        <Tag color="blue">{record.newapi_username}</Tag>
      ) : (
        <Tag>未绑定</Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: number) => (
        <Tag color={status === 1 ? 'green' : 'red'}>{status === 1 ? '启用' : '禁用'}</Tag>
      ),
    },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at', render: (time: string) => dayjs(time).format('YYYY-MM-DD') },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: any) => (
        <Button type="link" onClick={() => handleViewUser(record)}>详情</Button>
      ),
    },
  ]

  return (
    <div>
      <Card>
        <div style={{ marginBottom: 16 }}>
          <Input.Search
            placeholder="搜索用户名或邮箱"
            allowClear
            onSearch={setKeyword}
            style={{ width: 300 }}
          />
        </div>

        <Table
          rowKey="id"
          columns={columns}
          dataSource={users}
          loading={loading}
          pagination={{
            ...pagination,
            onChange: (page) => setPagination(prev => ({ ...prev, current: page })),
          }}
        />
      </Card>

      <Drawer
        title="用户详情"
        width={500}
        open={drawerVisible}
        onClose={() => setDrawerVisible(false)}
      >
        {detailLoading ? (
          <div style={{ textAlign: 'center', padding: 50 }}><Spin /></div>
        ) : userDetail ? (
          <>
            <Descriptions column={1} bordered size="small">
              <Descriptions.Item label="用户名">{userDetail.user?.username}</Descriptions.Item>
              <Descriptions.Item label="邮箱">{userDetail.user?.email || '-'}</Descriptions.Item>
              <Descriptions.Item label="角色">{userDetail.user?.role >= 10 ? '管理员' : '普通用户'}</Descriptions.Item>
              <Descriptions.Item label="new-api 账号">{userDetail.user?.newapi_username || '-'}</Descriptions.Item>
              <Descriptions.Item label="当前余额">{userDetail.current_quota}</Descriptions.Item>
            </Descriptions>

            {userDetail.subscription && (
              <>
                <h4 style={{ marginTop: 24 }}>订阅信息</h4>
                <Descriptions column={1} bordered size="small">
                  <Descriptions.Item label="套餐">{userDetail.subscription.plan?.name}</Descriptions.Item>
                  <Descriptions.Item label="状态">
                    <Tag color={userDetail.subscription.status === 'active' ? 'green' : 'red'}>
                      {userDetail.subscription.status}
                    </Tag>
                  </Descriptions.Item>
                  <Descriptions.Item label="每日额度">{userDetail.subscription.daily_quota}</Descriptions.Item>
                  <Descriptions.Item label="今日额度">{userDetail.subscription.today_quota}</Descriptions.Item>
                  <Descriptions.Item label="结转额度">{userDetail.subscription.carried_quota}</Descriptions.Item>
                  <Descriptions.Item label="到期时间">{dayjs(userDetail.subscription.end_date).format('YYYY-MM-DD')}</Descriptions.Item>
                </Descriptions>
              </>
            )}
          </>
        ) : null}
      </Drawer>
    </div>
  )
}
