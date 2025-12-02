import { useState, useEffect } from 'react'
import { Card, Form, Input, Switch, Button, message, Divider, Spin } from 'antd'
import { adminApi } from '../../../api'

export default function AdminSettings() {
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [syncing, setSyncing] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => {
    loadSettings()
  }, [])

  const loadSettings = async () => {
    try {
      const res: any = await adminApi.getSettings()
      if (res.success) {
        const data = res.data || {}
        form.setFieldsValue({
          site_name: data.site_name,
          site_description: data.site_description,
          require_login: data.require_login === '1',
          allow_register: data.allow_register === '1',
          newapi_login_enabled: data.newapi_login_enabled === '1',
        })
      }
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handleSave = async (values: any) => {
    setSaving(true)
    try {
      const data = {
        site_name: values.site_name,
        site_description: values.site_description,
        require_login: values.require_login ? '1' : '0',
        allow_register: values.allow_register ? '1' : '0',
        newapi_login_enabled: values.newapi_login_enabled ? '1' : '0',
      }
      const res: any = await adminApi.updateSettings(data)
      if (res.success) {
        message.success('保存成功')
      } else {
        message.error(res.message)
      }
    } catch (error: any) {
      message.error(error.message || '保存失败')
    } finally {
      setSaving(false)
    }
  }

  const handleSync = async () => {
    setSyncing(true)
    try {
      const res: any = await adminApi.triggerSync()
      if (res.success) {
        message.success('同步任务已启动')
      } else {
        message.error(res.message)
      }
    } catch (error: any) {
      message.error(error.message || '操作失败')
    } finally {
      setSyncing(false)
    }
  }

  if (loading) {
    return <div style={{ textAlign: 'center', padding: 50 }}><Spin /></div>
  }

  return (
    <div style={{ maxWidth: 800 }}>
      <Card title="站点设置" style={{ marginBottom: 24 }}>
        <Form form={form} layout="vertical" onFinish={handleSave}>
          <Form.Item name="site_name" label="站点名称">
            <Input />
          </Form.Item>
          <Form.Item name="site_description" label="站点描述">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Divider />
          <Form.Item name="require_login" label="访问需要登录" valuePropName="checked" extra="开启后未登录用户无法查看套餐列表">
            <Switch />
          </Form.Item>
          <Form.Item name="allow_register" label="允许注册" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item name="newapi_login_enabled" label="允许 new-api 登录" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={saving}>保存设置</Button>
          </Form.Item>
        </Form>
      </Card>

      <Card title="系统操作">
        <p style={{ marginBottom: 16, color: '#666' }}>
          手动触发订阅额度同步，通常每天 0:00 自动执行
        </p>
        <Button onClick={handleSync} loading={syncing}>
          立即同步额度
        </Button>
      </Card>
    </div>
  )
}
