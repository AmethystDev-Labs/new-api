import React, { useEffect, useState, useMemo } from 'react';
import { Table, Tag, Progress, Typography, Card, Space, Spin, Toast, Input } from '@douyinfe/semi-ui';
import { IconSearch } from '@douyinfe/semi-icons';
import { API } from '../../helpers';

const { Title, Text } = Typography;

const Status = () => {
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState([]);
  const [searchKeyword, setSearchKeyword] = useState('');

  const loadStats = async () => {
    try {
      const res = await API.get('/api/model/stats');
      const { success, data } = res.data;
      if (success) {
        setStats(data || []);
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

  const filteredStats = useMemo(() => {
    if (!searchKeyword) return stats;
    return stats.filter(item => 
      item.model_name.toLowerCase().includes(searchKeyword.toLowerCase())
    );
  }, [stats, searchKeyword]);

  const columns = [
    {
      title: '模型名称',
      dataIndex: 'model_name',
      key: 'model_name',
      sorter: (a, b) => a.model_name.localeCompare(b.model_name),
      render: (text) => (
        <Text strong style={{ fontSize: '14px' }}>
          {text}
        </Text>
      ),
    },
    {
      title: '请求次数 (24h)',
      dataIndex: 'total_count',
      key: 'total_count',
      sorter: (a, b) => (a.success_count + a.error_count) - (b.success_count + b.error_count),
      render: (_, record) => {
        return (
          <Space vertical align='start' spacing='extra-tight'>
            <Tag color='green' size='small' variant='light'>
              成功: {record.success_count}
            </Tag>
            <Tag color='red' size='small' variant='light'>
              失败: {record.error_count}
            </Tag>
          </Space>
        );
      },
    },
    {
      title: '成功率',
      key: 'success_rate',
      sorter: (a, b) => {
        const rateA = a.success_count / (a.success_count + a.error_count || 1);
        const rateB = b.success_count / (b.success_count + b.error_count || 1);
        return rateA - rateB;
      },
      render: (_, record) => {
        const total = record.success_count + record.error_count;
        const rate = total === 0 ? 0 : Math.round((record.success_count / total) * 100);
        let color = '#3bb346'; // green
        if (rate < 80) color = '#ff9d00'; // orange
        if (rate < 50) color = '#f5222d'; // red

        return (
          <div style={{ width: '200px' }}>
            <Space>
              <Progress
                percent={rate}
                stroke={color}
                showInfo={false}
                style={{ width: '120px' }}
              />
              <Text strong>{rate}%</Text>
            </Space>
          </div>
        );
      },
    },
    {
      title: '平均延迟',
      dataIndex: 'avg_latency',
      key: 'avg_latency',
      sorter: (a, b) => a.avg_latency - b.avg_latency,
      render: (text) => (
        <Tag color='blue' variant='light'>
          {text.toFixed(2)}ms
        </Tag>
      ),
    },
  ];

  return (
    <div style={{ padding: '80px 24px 24px 24px', maxWidth: '1200px', margin: '0 auto' }}>
      <Card
        title={
          <Space vertical align='start'>
            <Title heading={2}>模型运行状态</Title>
            <Text type='secondary'>
              展示过去 24 小时内各模型的历史请求成功率及平均响应时间。数据每 10 分钟聚合更新一次。
            </Text>
          </Space>
        }
        headerExtraContent={
          <Input
            prefix={<IconSearch />}
            placeholder="搜索模型名称..."
            value={searchKeyword}
            onChange={value => setSearchKeyword(value)}
            style={{ width: 250 }}
            showClear
          />
        }
      >
        <Table
          columns={columns}
          dataSource={filteredStats}
          loading={loading}
          pagination={false}
          empty='暂无统计数据'
        />
      </Card>
    </div>
  );
};

export default Status;
