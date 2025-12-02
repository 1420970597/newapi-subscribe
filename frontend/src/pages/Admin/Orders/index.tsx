import { useState, useEffect } from 'react'
import { Table, Card, Select, Tag } from 'antd'
import dayjs from 'dayjs'
import { adminApi } from '../../../api'

const statusMap: Record<string, { color: string; text: string }> = {
  pending: { color: 'orange', text: '待支付' },
  paid: { color: 'green', text: '已支付' },
  cancelled: { color: 'default', text: '已取消' },
  refunded: { color: 'red', text: '已退款' },
}

export default function AdminOrders() {
  const [loading, setLoading] = useState(true)
  const [orders, setOrders] = useState<any[]>([])
  const [pagination, setPagination] = useState({ current: 1, pageSize: 20, total: 0 })
  const [status, setStatus] = useState<string>('')

  useEffect(() => {
    loadOrders()
  }, [pagination.current, status])

  const loadOrders = async () => {
    setLoading(true)
    try {
      const res: any = await adminApi.getOrders({
        page: pagination.current,
        per_page: pagination.pageSize,
        status: status || undefined,
      })
      if (res.success) {
        setOrders(res.data || [])
        setPagination(prev => ({ ...prev, total: res.total }))
      }
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '订单号', dataIndex: 'order_no', key: 'order_no' },
    { title: '用户', dataIndex: ['user', 'username'], key: 'user' },
    { title: '套餐', dataIndex: ['plan', 'name'], key: 'plan' },
    { title: '类型', dataIndex: 'order_type', key: 'order_type', render: (type: string) => type === 'new' ? '新购' : '续费' },
    { title: '天数', dataIndex: 'period_days', key: 'period_days' },
    { title: '金额', dataIndex: 'amount', key: 'amount', render: (amount: number) => `¥${amount}` },
    { title: '支付方式', dataIndex: 'payment_method', key: 'payment_method' },
    { title: '状态', dataIndex: 'status', key: 'status', render: (s: string) => <Tag color={statusMap[s]?.color}>{statusMap[s]?.text}</Tag> },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at', render: (time: string) => dayjs(time).format('YYYY-MM-DD HH:mm') },
    { title: '支付时间', dataIndex: 'paid_at', key: 'paid_at', render: (time: string) => time ? dayjs(time).format('YYYY-MM-DD HH:mm') : '-' },
  ]

  return (
    <div>
      <Card>
        <div style={{ marginBottom: 16 }}>
          <Select
            style={{ width: 150 }}
            placeholder="筛选状态"
            allowClear
            value={status || undefined}
            onChange={setStatus}
            options={[
              { value: 'pending', label: '待支付' },
              { value: 'paid', label: '已支付' },
              { value: 'cancelled', label: '已取消' },
              { value: 'refunded', label: '已退款' },
            ]}
          />
        </div>
        <Table
          rowKey="id"
          columns={columns}
          dataSource={orders}
          loading={loading}
          pagination={{
            ...pagination,
            onChange: (page) => setPagination(prev => ({ ...prev, current: page })),
          }}
        />
      </Card>
    </div>
  )
}
