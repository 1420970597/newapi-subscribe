import { useState } from 'react'
import { Card, Form, Input, Button, Switch, InputNumber, message, Divider, Tag } from 'antd'
import { userApi } from '../../../api'
import { useAuthStore } from '../../../store/auth'

export default function Settings() {
  const { user } = useAuthStore()
  const [loading, setLoading] = useState(false)
  const [bindLoading, setBindLoading] = useState(false)
  const [bindForm] = Form.useForm()

  const handleUpdateProfile = async (values: any) => {
    setLoading(true)
    try {
      const res: any = await userApi.updateProfile(values)
      if (res.success) {
        message.success('更新成功')
      } else {
        message.error(res.message)
      }
    } catch (error: any) {
      message.error(error.message || '更新失败')
    } finally {
      setLoading(false)
    }
  }

  const handleUpdateEmailSettings = async (values: any) => {
    setLoading(true)
    try {
      const res: any = await userApi.updateEmailSettings({
        email_remind: values.email_remind ? 1 : 0,
        remind_days: values.remind_days,
      })
      if (res.success) {
        message.success('更新成功')
      } else {
        message.error(res.message)
      }
    } catch (error: any) {
      message.error(error.message || '更新失败')
    } finally {
      setLoading(false)
    }
  }

  const handleBindNewAPI = async (values: any) => {
    setBindLoading(true)
    try {
      const res: any = await userApi.bindNewAPI(values)
      if (res.success) {
        message.success('绑定成功')
        bindForm.resetFields()
        window.location.reload()
      } else {
        message.error(res.message)
      }
    } catch (error: any) {
      message.error(error.message || '绑定失败')
    } finally {
      setBindLoading(false)
    }
  }

  return (
    <div style={{ maxWidth: 800, margin: '0 auto' }}>
      <h2 style={{ marginBottom: 24 }}>账户设置</h2>

      <Card title="个人信息" style={{ marginBottom: 24 }}>
        <Form
          layout="vertical"
          initialValues={{ email: user?.email }}
          onFinish={handleUpdateProfile}
        >
          <Form.Item label="用户名">
            <Input value={user?.username} disabled />
          </Form.Item>
          <Form.Item name="email" label="邮箱" rules={[{ type: 'email', message: '请输入有效的邮箱' }]}>
            <Input placeholder="用于接收到期提醒" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>保存</Button>
          </Form.Item>
        </Form>
      </Card>

      <Card title="new-api 账号绑定" style={{ marginBottom: 24 }}>
        {user?.newapi_bound ? (
          <div>
            <p>已绑定账号: <Tag color="green">{user.newapi_username}</Tag></p>
          </div>
        ) : (
          <Form form={bindForm} layout="vertical" onFinish={handleBindNewAPI}>
            <Form.Item name="username" label="new-api 用户名" rules={[{ required: true }]}>
              <Input />
            </Form.Item>
            <Form.Item name="password" label="new-api 密码" rules={[{ required: true }]}>
              <Input.Password />
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit" loading={bindLoading}>绑定</Button>
            </Form.Item>
          </Form>
        )}
      </Card>

      <Card title="邮件提醒">
        <Form
          layout="vertical"
          initialValues={{
            email_remind: user?.email_remind === 1,
            remind_days: user?.remind_days || 3,
          }}
          onFinish={handleUpdateEmailSettings}
        >
          <Form.Item name="email_remind" label="开启到期提醒" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item name="remind_days" label="提前提醒天数">
            <InputNumber min={1} max={30} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>保存</Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}
