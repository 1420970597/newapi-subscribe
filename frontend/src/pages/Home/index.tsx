import { useState, useEffect } from 'react'
import { Button, Spin, message, Modal } from 'antd'
import { CheckOutlined, ThunderboltOutlined, CrownOutlined, RocketOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { planApi } from '../../api'
import { useAuthStore } from '../../store/auth'
import './home.css'

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

const planIcons: Record<number, React.ReactNode> = {
  0: <ThunderboltOutlined />,
  1: <CrownOutlined />,
  2: <RocketOutlined />,
}

const planColors: string[] = [
  'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
  'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
  'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
  'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
  'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
]

export default function Home() {
  const navigate = useNavigate()
  const { isAuthenticated } = useAuthStore()
  const [plans, setPlans] = useState<Plan[]>([])
  const [loading, setLoading] = useState(true)
  const [hoveredPlan, setHoveredPlan] = useState<number | null>(null)

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
      <div className="loading-container">
        <Spin size="large" />
      </div>
    )
  }

  return (
    <div className="home-container">
      {/* 背景装饰 */}
      <div className="bg-decoration">
        <div className="bg-blob bg-blob-1" />
        <div className="bg-blob bg-blob-2" />
        <div className="bg-blob bg-blob-3" />
      </div>

      {/* 标题区域 */}
      <div className="hero-section">
        <h1 className="hero-title">
          <span className="title-gradient">选择您的专属套餐</span>
        </h1>
        <p className="hero-subtitle">
          灵活的订阅方案，为您提供最优质的 AI 服务体验
        </p>
      </div>

      {/* 套餐卡片 */}
      <div className="plans-grid">
        {plans.map((plan, index) => (
          <div
            key={plan.id}
            className={`plan-card ${hoveredPlan === plan.id ? 'hovered' : ''} ${index === 0 ? 'featured' : ''}`}
            onMouseEnter={() => setHoveredPlan(plan.id)}
            onMouseLeave={() => setHoveredPlan(null)}
            style={{
              animationDelay: `${index * 0.1}s`,
            }}
          >
            {/* 卡片光效 */}
            <div className="card-glow" style={{ background: planColors[index % planColors.length] }} />

            {/* 卡片内容 */}
            <div className="card-content">
              {/* 图标 */}
              <div className="plan-icon" style={{ background: planColors[index % planColors.length] }}>
                {planIcons[index % 3]}
              </div>

              {/* 套餐名称 */}
              <h2 className="plan-name">{plan.name}</h2>

              {/* 描述 */}
              <p className="plan-description">{plan.description}</p>

              {/* 价格 */}
              <div className="price-section">
                <span className="currency">¥</span>
                <span className="price-value">{plan.price}</span>
                <span className="price-period">
                  /{plan.price_type === 'fixed' ? `${plan.period_days}${periodTypeMap[plan.period_type]}` : '天'}
                </span>
              </div>

              {/* 特性列表 */}
              <ul className="features-list">
                <li className="feature-item">
                  <CheckOutlined className="feature-icon" />
                  <span>每日额度 {plan.daily_quota.toLocaleString()}</span>
                </li>
                <li className="feature-item">
                  <CheckOutlined className="feature-icon" />
                  <span>{plan.carry_over ? '支持额度结转' : '额度不结转'}</span>
                </li>
                <li className="feature-item">
                  <CheckOutlined className="feature-icon" />
                  <span>分组: {plan.newapi_group}</span>
                </li>
                <li className="feature-item">
                  <CheckOutlined className="feature-icon" />
                  <span>有效期 {plan.period_days} 天</span>
                </li>
              </ul>

              {/* 订阅按钮 */}
              <Button
                type="primary"
                size="large"
                className="subscribe-btn"
                onClick={() => handlePurchase(plan)}
                style={{
                  background: planColors[index % planColors.length],
                  border: 'none',
                }}
              >
                立即订阅
              </Button>
            </div>
          </div>
        ))}
      </div>

      {plans.length === 0 && (
        <div className="empty-state">
          <p>暂无可用套餐</p>
        </div>
      )}
    </div>
  )
}
