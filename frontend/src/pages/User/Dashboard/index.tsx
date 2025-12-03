import { useState, useEffect } from 'react'
import { Card, Row, Col, Progress, Button, Tag, Statistic, Empty, Spin, message, Modal, Form, Select, Input, Radio } from 'antd'
import { useSearchParams, useNavigate } from 'react-router-dom'
import { subscriptionApi, planApi, orderApi } from '../../../api'

export default function Dashboard() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const planId = searchParams.get('plan')

  const [loading, setLoading] = useState(true)
  const [subscription, setSubscription] = useState<any>(null)
  const [currentQuota, setCurrentQuota] = useState(0)
  const [todayUsage, setTodayUsage] = useState<any>(null)
  const [purchaseModal, setPurchaseModal] = useState(false)
  const [plans, setPlans] = useState<any[]>([])
  const [selectedPlan, setSelectedPlan] = useState<any>(null)
  const [purchaseLoading, setPurchaseLoading] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => {
    loadData()
  }, [])

  useEffect(() => {
    if (planId) {
      loadPlans().then(() => {
        setPurchaseModal(true)
      })
    }
  }, [planId])

  const loadData = async () => {
    try {
      const res: any = await subscriptionApi.current()
      if (res.success && res.data) {
        setSubscription(res.data.subscription)
        setCurrentQuota(res.data.current_quota)
      }
      // 加载今日用量
      const usageRes: any = await subscriptionApi.todayUsage()
      if (usageRes.success) {
        setTodayUsage(usageRes.data)
      }
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const loadPlans = async () => {
    const res: any = await planApi.list()
    if (res.success) {
      setPlans(res.data || [])
      if (planId) {
        const plan = res.data?.find((p: any) => p.id === parseInt(planId))
        if (plan) setSelectedPlan(plan)
      }
    }
  }

  const handlePurchase = async (values: any) => {
    if (!selectedPlan) return
    setPurchaseLoading(true)
    try {
      const purchaseRes: any = await subscriptionApi.purchase({
        plan_id: selectedPlan.id,
        newapi_action: values.newapi_action,
        newapi_username: values.newapi_username,
        newapi_password: values.newapi_password,
      })

      if (purchaseRes.success) {
        const payRes: any = await orderApi.pay({
          order_id: purchaseRes.data.order.id,
          payment_method: values.payment_method,
        })

        if (payRes.success) {
          window.location.href = payRes.data.pay_url
        } else {
          message.error(payRes.message)
        }
      } else {
        message.error(purchaseRes.message)
      }
    } catch (error: any) {
      message.error(error.message || '操作失败')
    } finally {
      setPurchaseLoading(false)
    }
  }

  if (loading) {
    return <div style={{ textAlign: 'center', padding: 100 }}><Spin size="large" /></div>
  }

  return (
    <div style={{ maxWidth: 1200, margin: '0 auto' }}>
      <h2 style={{ marginBottom: 24 }}>我的订阅</h2>

      {subscription ? (
        <Row gutter={24}>
          <Col span={16}>
            <Card>
              <Row gutter={24}>
                <Col span={12}>
                  <h3>{subscription.plan?.name || '订阅套餐'}</h3>
                  <Tag color={subscription.status === 'active' ? 'green' : 'red'}>
                    {subscription.status === 'active' ? '生效中' : '已过期'}
                  </Tag>
                </Col>
                <Col span={12} style={{ textAlign: 'right' }}>
                  <Button type="primary" onClick={() => { loadPlans(); setPurchaseModal(true) }}>
                    续费
                  </Button>
                </Col>
              </Row>

              <div style={{ marginTop: 24 }}>
                <p style={{ marginBottom: 8 }}>订阅周期</p>
                <Progress
                  percent={Math.round((1 - subscription.days_remaining / subscription.period_days) * 100)}
                  format={() => `剩余 ${subscription.days_remaining || 0} 天`}
                />
              </div>

              <div style={{ marginTop: 24 }}>
                <p style={{ marginBottom: 8 }}>今日额度</p>
                <Progress
                  percent={subscription.today_quota > 0 ? Math.round((currentQuota / subscription.today_quota) * 100) : 0}
                  format={() => `${currentQuota} / ${subscription.today_quota}`}
                  status={currentQuota > 0 ? 'active' : 'exception'}
                />
              </div>
            </Card>
          </Col>

          <Col span={8}>
            <Card>
              <Statistic title="每日额度" value={subscription.daily_quota} />
              <Statistic title="结转额度" value={subscription.carried_quota} style={{ marginTop: 16 }} />
              <Statistic title="模型分组" value={subscription.newapi_group} style={{ marginTop: 16 }} />
              {todayUsage && (
                <Statistic
                  title="今日已用"
                  value={todayUsage.today_used || 0}
                  suffix={`/ ${todayUsage.daily_quota || subscription.daily_quota}`}
                  style={{ marginTop: 16 }}
                  valueStyle={{ color: (todayUsage.today_used || 0) > (todayUsage.daily_quota || subscription.daily_quota) * 0.8 ? '#cf1322' : '#3f8600' }}
                />
              )}
            </Card>
          </Col>
        </Row>
      ) : (
        <Card>
          <Empty description="暂无订阅">
            <Button type="primary" onClick={() => navigate('/')}>
              选择套餐
            </Button>
          </Empty>
        </Card>
      )}

      <Modal
        title="购买订阅"
        open={purchaseModal}
        onCancel={() => setPurchaseModal(false)}
        footer={null}
        width={500}
      >
        <Form form={form} onFinish={handlePurchase} layout="vertical">
          <Form.Item label="选择套餐" required>
            <Select
              value={selectedPlan?.id}
              onChange={(id) => setSelectedPlan(plans.find(p => p.id === id))}
              placeholder="请选择套餐"
            >
              {plans.map(plan => (
                <Select.Option key={plan.id} value={plan.id}>
                  {plan.name} - ¥{plan.price}/{plan.period_days}天
                </Select.Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item name="newapi_action" label="new-api 账号" rules={[{ required: true }]} initialValue="create_new">
            <Radio.Group>
              <Radio value="create_new">创建新账号</Radio>
              <Radio value="bind_existing">绑定现有账号</Radio>
              <Radio value="overwrite">覆盖当前账号</Radio>
            </Radio.Group>
          </Form.Item>

          <Form.Item noStyle shouldUpdate={(prev, cur) => prev.newapi_action !== cur.newapi_action}>
            {({ getFieldValue }) => getFieldValue('newapi_action') === 'bind_existing' && (
              <>
                <Form.Item name="newapi_username" label="new-api 用户名" rules={[{ required: true }]}>
                  <Input />
                </Form.Item>
                <Form.Item name="newapi_password" label="new-api 密码" rules={[{ required: true }]}>
                  <Input.Password />
                </Form.Item>
              </>
            )}
          </Form.Item>

          <Form.Item name="payment_method" label="支付方式" rules={[{ required: true }]} initialValue="alipay">
            <Radio.Group>
              <Radio value="alipay">支付宝</Radio>
              <Radio value="wxpay">微信</Radio>
            </Radio.Group>
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" block loading={purchaseLoading}>
              立即支付 {selectedPlan && `¥${selectedPlan.price}`}
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
