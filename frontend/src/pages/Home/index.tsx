import { useState, useEffect } from 'react'
import { Card, Row, Col, Button, Tag, Spin, message, Modal } from 'antd'
import { useNavigate } from 'react-router-dom'
import { planApi } from '../../api'
import { useAuthStore } from '../../store/auth'

interface Plan {
  id: number
  name: string
  description: string
  period_type: string
  period_days: number
  daily_quota: number
  carry_over: number
  price_type: string
  price: number
  newapi_group: string
}

const periodTypeMap: Record<string, string> = {
  day: '天',
  week: '周',
  month: '月',
  custom: '自定义',
}

export default function Home() {
  const navigate = useNavigate()
  const { isAuthenticated } = useAuthStore()
  const [plans, setPlans] = useState<Plan[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadPlans()
  }, [])

  const loadPlans = async () => {
    try {
      const res: any = await planApi.list()
      if (res.success) {
        setPlans(res.data || [])
      }
    } catch (error: any) {
      message.error(error.message || '加载失败')
    } finally {
      setLoading(false)
    }
  }

  const handlePurchase = (plan: Plan) => {
    if (!isAuthenticated) {
      Modal.confirm({
        title: '请先登录',
        content: '购买订阅需要先登录账号',
        okText: '去登录',
        cancelText: '取消',
        onOk: () => navigate('/login'),
      })
      return
    }
    navigate(`/user?plan=${plan.id}`)
  }

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 100 }}>
        <Spin size="large" />
      </div>
    )
  }

  return (
    <div>
      <div style={{ textAlign: 'center', marginBottom: 40 }}>
        <h1 style={{ fontSize: 32, marginBottom: 16 }}>选择适合您的套餐</h1>
        <p style={{ color: '#666' }}>灵活的订阅方案，满足不同需求</p>
      </div>

      <Row gutter={[24, 24]} justify="center">
        {plans.map((plan) => (
          <Col key={plan.id} xs={24} sm={12} md={8} lg={6}>
            <Card
              hoverable
              style={{ height: '100%' }}
              styles={{ body: { display: 'flex', flexDirection: 'column', height: '100%' } }}
            >
              <div style={{ textAlign: 'center', marginBottom: 16 }}>
                <h2 style={{ fontSize: 24, marginBottom: 8 }}>{plan.name}</h2>
                <p style={{ color: '#666', minHeight: 44 }}>{plan.description}</p>
              </div>

              <div style={{ textAlign: 'center', marginBottom: 24 }}>
                <span style={{ fontSize: 36, fontWeight: 'bold', color: '#1890ff' }}>
                  ¥{plan.price}
                </span>
                <span style={{ color: '#666' }}>
                  /{plan.price_type === 'fixed' ? `${plan.period_days}${periodTypeMap[plan.period_type]}` : '天'}
                </span>
              </div>

              <div style={{ marginBottom: 24, flex: 1 }}>
                <div style={{ marginBottom: 8 }}>
                  <Tag color="blue">每日额度: {plan.daily_quota}</Tag>
                </div>
                <div style={{ marginBottom: 8 }}>
                  <Tag color={plan.carry_over ? 'green' : 'default'}>
                    {plan.carry_over ? '支持结转' : '不结转'}
                  </Tag>
                </div>
                <div>
                  <Tag color="purple">分组: {plan.newapi_group}</Tag>
                </div>
              </div>

              <Button type="primary" block size="large" onClick={() => handlePurchase(plan)}>
                立即订阅
              </Button>
            </Card>
          </Col>
        ))}
      </Row>

      {plans.length === 0 && (
        <div style={{ textAlign: 'center', padding: 60, color: '#999' }}>
          暂无可用套餐
        </div>
      )}
    </div>
  )
}
