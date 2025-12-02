import { useState, useEffect } from 'react'
import { Table, Tag, Card } from 'antd'
import dayjs from 'dayjs'
import { orderApi } from '../../../api'

const statusMap: Record<string, { color: string; text: string }> = {
  pending: { color: 'orange', text: '待支付' },
  paid: { color: 'green', text: '已支付' },
  cancelled: { color: 'default', text: '已取消' },
  refunded: { color: 'red', text: '已退款' },
}

export default function Orders() {
  const [loading, setLoading] = useState(true)
  const [orders, setOrders] = useState<any[]>([])
  const [pagination, setPagination] = useState({ current: 1, pageSize: 10, total: 0 })

  useEffect(() => {
    loadOrders()
  }, [pagination.current])

  const loadOrders = async () => {
    setLoading(true)
    try {
      const res: any = await orderApi.list({ page: pagination.current, per_page: pagination.pageSize })
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
    { title: '订单号', dataIndex: 'order_no', key: 'order_no' },
    { title: '套餐', dataIndex: ['plan', 'name'], key: 'plan' },
    { title: '类型', dataIndex: 'order_type', key: 'order_type', render: (type: string) => type === 'new' ? '新购' : '续费' },
    { title: '天数', dataIndex: 'period_days', key: 'period_days', render: (days: number) => `${days}天` },
    { title: '金额', dataIndex: 'amount', key: 'amount', render: (amount: number) => `¥${amount}` },
    { title: '状态', dataIndex: 'status', key: 'status', render: (status: string) => <Tag color={statusMap[status]?.color}>{statusMap[status]?.text}</Tag> },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at', render: (time: string) => dayjs(time).format('YYYY-MM-DD HH:mm') },
  ]

  return (
    <div style={{ maxWidth: 1200, margin: '0 auto' }}>
      <h2 style={{ marginBottom: 24 }}>订单记录</h2>
      <Card>
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
