import { useState, useEffect } from 'react'
import { Table, Card, Button, Modal, Form, Input, InputNumber, Select, Switch, message, Tag, Popconfirm } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { adminApi, planApi } from '../../../api'

export default function AdminPlans() {
  const [loading, setLoading] = useState(true)
  const [plans, setPlans] = useState<any[]>([])
  const [modalVisible, setModalVisible] = useState(false)
  const [editingPlan, setEditingPlan] = useState<any>(null)
  const [groups, setGroups] = useState<string[]>([])
  const [submitLoading, setSubmitLoading] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => {
    loadPlans()
    loadGroups()
  }, [])

  const loadPlans = async () => {
    setLoading(true)
    try {
      const res: any = await planApi.list()
      if (res.success) {
        setPlans(res.data || [])
      }
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const loadGroups = async () => {
    try {
      const res: any = await adminApi.getNewAPIGroups()
      if (res.success) {
        setGroups(res.data || [])
      }
    } catch (error) {
      console.error(error)
    }
  }

  const handleSubmit = async (values: any) => {
    setSubmitLoading(true)
    try {
      const data = { ...values, status: values.status ? 1 : 0, carry_over: values.carry_over ? 1 : 0 }
      let res: any
      if (editingPlan) {
        res = await adminApi.updatePlan(editingPlan.id, data)
      } else {
        res = await adminApi.createPlan(data)
      }

      if (res.success) {
        message.success(editingPlan ? '更新成功' : '创建成功')
        setModalVisible(false)
        form.resetFields()
        loadPlans()
      } else {
        message.error(res.message)
      }
    } catch (error: any) {
      message.error(error.message || '操作失败')
    } finally {
      setSubmitLoading(false)
    }
  }

  const handleDelete = async (id: number) => {
    try {
      const res: any = await adminApi.deletePlan(id)
      if (res.success) {
        message.success('删除成功')
        loadPlans()
      } else {
        message.error(res.message)
      }
    } catch (error: any) {
      message.error(error.message || '删除失败')
    }
  }

  const openModal = (plan?: any) => {
    setEditingPlan(plan || null)
    if (plan) {
      form.setFieldsValue({ ...plan, status: plan.status === 1, carry_over: plan.carry_over === 1 })
    } else {
      form.resetFields()
    }
    setModalVisible(true)
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '周期', key: 'period', render: (_: any, r: any) => `${r.period_days}天` },
    { title: '每日额度', dataIndex: 'daily_quota', key: 'daily_quota' },
    { title: '结转', dataIndex: 'carry_over', key: 'carry_over', render: (v: number) => v ? <Tag color="green">是</Tag> : <Tag>否</Tag> },
    { title: '价格', key: 'price', render: (_: any, r: any) => `¥${r.price}${r.price_type === 'daily' ? '/天' : ''}` },
    { title: '分组', dataIndex: 'newapi_group', key: 'newapi_group' },
    { title: '状态', dataIndex: 'status', key: 'status', render: (s: number) => <Tag color={s === 1 ? 'green' : 'default'}>{s === 1 ? '上架' : '下架'}</Tag> },
    {
      title: '操作', key: 'action',
      render: (_: any, record: any) => (
        <>
          <Button type="link" onClick={() => openModal(record)}>编辑</Button>
          <Popconfirm title="确定删除？" onConfirm={() => handleDelete(record.id)}>
            <Button type="link" danger>删除</Button>
          </Popconfirm>
        </>
      ),
    },
  ]

  return (
    <div>
      <Card>
        <div style={{ marginBottom: 16 }}>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => openModal()}>新建套餐</Button>
        </div>
        <Table rowKey="id" columns={columns} dataSource={plans} loading={loading} />
      </Card>

      <Modal
        title={editingPlan ? '编辑套餐' : '新建套餐'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        onOk={() => form.submit()}
        confirmLoading={submitLoading}
        width={600}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit} initialValues={{ status: true, carry_over: false, period_type: 'month', price_type: 'fixed' }}>
          <Form.Item name="name" label="套餐名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item name="period_type" label="周期类型" rules={[{ required: true }]}>
            <Select options={[{ value: 'day', label: '天' }, { value: 'week', label: '周' }, { value: 'month', label: '月' }, { value: 'custom', label: '自定义' }]} />
          </Form.Item>
          <Form.Item name="period_days" label="周期天数" rules={[{ required: true }]}>
            <InputNumber min={1} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="daily_quota" label="每日额度" rules={[{ required: true }]}>
            <InputNumber min={1} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="carry_over" label="支持结转" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item name="max_carry_over" label="最大结转额度" extra="0 表示无限制">
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="price_type" label="价格类型" rules={[{ required: true }]}>
            <Select options={[{ value: 'fixed', label: '固定价格' }, { value: 'daily', label: '按天计价' }]} />
          </Form.Item>
          <Form.Item name="price" label="价格" rules={[{ required: true }]}>
            <InputNumber min={0} precision={2} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="newapi_group" label="new-api 分组" rules={[{ required: true }]}>
            <Select options={groups.map(g => ({ value: g, label: g }))} />
          </Form.Item>
          <Form.Item name="status" label="上架" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item name="sort_order" label="排序">
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
