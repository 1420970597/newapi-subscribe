import { useState, useEffect } from 'react'
import { Card, Table, DatePicker, Empty, Spin, Alert } from 'antd'
import dayjs from 'dayjs'
import { subscriptionApi } from '../../../api'
import { useAuthStore } from '../../../store/auth'

const { RangePicker } = DatePicker

export default function Usage() {
  const { user } = useAuthStore()
  const [loading, setLoading] = useState(false)
  const [logs, setLogs] = useState<any[]>([])
  const [dateRange, setDateRange] = useState<[dayjs.Dayjs, dayjs.Dayjs] | null>(null)

  useEffect(() => {
    if (user?.newapi_bound) {
      loadLogs()
    }
  }, [dateRange])

  const loadLogs = async () => {
    setLoading(true)
    try {
      const params: any = {}
      if (dateRange) {
        params.start_date = dateRange[0].unix()
        params.end_date = dateRange[1].unix()
      }
      const res: any = await subscriptionApi.usageDetail(params)
      if (res.success) {
        setLogs(res.data || [])
      }
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  if (!user?.newapi_bound) {
    return (
      <div style={{ maxWidth: 1200, margin: '0 auto' }}>
        <h2 style={{ marginBottom: 24 }}>使用统计</h2>
        <Alert
          message="未绑定 new-api 账号"
          description="请先在账户设置中绑定 new-api 账号，或前往 new-api 站点查询使用记录"
          type="warning"
          showIcon
        />
      </div>
    )
  }

  const columns = [
    { title: '时间', dataIndex: 'created_at', key: 'created_at', render: (time: number) => dayjs.unix(time).format('YYYY-MM-DD HH:mm:ss') },
    { title: '模型', dataIndex: 'model_name', key: 'model_name' },
    { title: '消耗', dataIndex: 'quota', key: 'quota' },
    { title: '输入Token', dataIndex: 'prompt_tokens', key: 'prompt_tokens' },
    { title: '输出Token', dataIndex: 'completion_tokens', key: 'completion_tokens' },
  ]

  return (
    <div style={{ maxWidth: 1200, margin: '0 auto' }}>
      <h2 style={{ marginBottom: 24 }}>使用统计</h2>

      <Card>
        <div style={{ marginBottom: 16 }}>
          <RangePicker
            value={dateRange}
            onChange={(dates) => setDateRange(dates as [dayjs.Dayjs, dayjs.Dayjs])}
          />
        </div>

        {loading ? (
          <div style={{ textAlign: 'center', padding: 50 }}><Spin /></div>
        ) : logs.length > 0 ? (
          <Table rowKey="id" columns={columns} dataSource={logs} pagination={{ pageSize: 20 }} />
        ) : (
          <Empty description="暂无使用记录" />
        )}
      </Card>
    </div>
  )
}
