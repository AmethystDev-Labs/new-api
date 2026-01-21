import React, { useEffect, useState } from 'react';
import { Table, Tag, Progress, Typography, Card, Space, Spin, Toast } from '@douyinfe/semi-ui';
import { API } from '../../helpers';

const { Title, Text } = Typography;

const Status = () => {
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState([]);

  const loadStats = async () => {
    try {
      const res = await API.get('/api/model/stats');
      const { success, data } = res.data;
      if (success) {
        setStats(data);
      } else {
        Toast.error('获取数据失败');
      }
    } catch (e) {
      console.error(e);
      Toast.error('网络错误，请稍后再试');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadStats();
    const interval = setInterval(loadStats, 60000); // 每分钟刷新一次
    return () => clearInterval(interval);
  }, []);

  const columns = [
    {
      title: '模型名称',
      dataIndex: 'model_name',
      key: 'model_name',
      render: (text) => <Text strong>{text}</Text>,
    },
    {
      title: '请求次数 (24h)',
      key: 'total_count',
      render: (_, record) => (
        <Space vertical align='start' spacing='tight'>
          <Tag color='blue' shape='ghost'>成功: {record.success_count}</Tag>
          <Tag color='red' shape='ghost'>失败: {record.error_count}</Tag>
        </Space>
      ),
      sorter: (a, b) => (a.success_count + a.error_count) - (b.success_count + b.error_count),
    },
    {
      title: '成功率',
      key: 'success_rate',
      render: (_, record) => {
        const total = record.success_count + record.error_count;
        const rate = total === 0 ? 0 : (record.success_count / total) * 100;
        let stroke = 'var(--semi-color-success)';
        if (rate < 90) stroke = 'var(--semi-color-warning)';
        if (rate < 70) stroke = 'var(--semi-color-danger)';
        
        return (
          <div style={{ width: 150 }}>
            <Progress 
              percent={parseFloat(rate.toFixed(1))} 
              showInfo
              stroke={stroke}
            />
          </div>
        );
      },
      sorter: (a, b) => {
        const rateA = (a.success_count + a.error_count) === 0 ? 0 : (a.success_count / (a.success_count + a.error_count));
        const rateB = (b.success_count + b.error_count) === 0 ? 0 : (b.success_count / (b.success_count + b.error_count));
        return rateA - rateB;
      },
    },
    {
      title: '平均延迟',
      dataIndex: 'avg_latency',
      key: 'avg_latency',
      render: (latency) => (
        <Tag color={latency < 1000 ? 'green' : latency < 3000 ? 'amber' : 'red'}>
          {latency.toFixed(2)}ms
        </Tag>
      ),
      sorter: (a, b) => a.avg_latency - b.avg_latency,
    },
  ];

  return (
    <div style={{ padding: '24px', maxWidth: '1200px', margin: '0 auto' }}>
      <Card bordered={false} style={{ marginBottom: '24px' }}>
        <Title heading={2}>模型运行状态</Title>
        <Text type="secondary">展示过去 24 小时内各模型的历史请求成功率及平均响应时间。数据每 10 分钟聚合更新一次。</Text>
      </Card>

      <Card bordered={false}>
        {loading ? (
          <div style={{ textAlign: 'center', padding: '50px' }}>
            <Spin size="large" />
            <div style={{ marginTop: '12px' }}>加载中...</div>
          </div>
        ) : (
          <Table 
            dataSource={stats} 
            columns={columns} 
            pagination={{ pageSize: 20 }}
          />
        )}
      </Card>
    </div>
  );
};

export default Status;
