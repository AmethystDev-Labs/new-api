/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

import React, { useMemo } from 'react';
import {
  Modal,
  Button,
  Empty,
  Space,
  Tag,
  Typography,
} from '@douyinfe/semi-ui';
import { IconCopy } from '@douyinfe/semi-icons';
import { copy, showError, showSuccess } from '../../../../helpers';

const { Text } = Typography;

const parseAuditLine = (line) => {
  if (typeof line !== 'string') {
    return null;
  }
  const colonIndex = line.indexOf(': ');
  const arrowIndex = line.indexOf(' -> ', colonIndex + 2);
  if (colonIndex <= 0 || arrowIndex <= colonIndex) {
    return null;
  }

  return {
    field: line.slice(0, colonIndex),
    before: line.slice(colonIndex + 2, arrowIndex),
    after: line.slice(arrowIndex + 4),
    raw: line,
  };
};

const ValuePanel = ({ label, value, tone }) => (
  <div
    style={{
      flex: 1,
      minWidth: 0,
      padding: 12,
      borderRadius: 12,
      border: '1px solid var(--semi-color-border)',
      background:
        tone === 'after'
          ? 'rgba(var(--semi-blue-5), 0.08)'
          : 'var(--semi-color-fill-0)',
    }}
  >
    <div
      style={{
        marginBottom: 6,
        fontSize: 12,
        fontWeight: 600,
        color: 'var(--semi-color-text-2)',
      }}
    >
      {label}
    </div>
    <Text
      style={{
        display: 'block',
        fontFamily:
          'ui-monospace, SFMono-Regular, SF Mono, Menlo, Consolas, monospace',
        fontSize: 13,
        lineHeight: 1.65,
        whiteSpace: 'pre-wrap',
        wordBreak: 'break-word',
      }}
    >
      {value}
    </Text>
  </div>
);

const ParamOverrideModal = ({
  showParamOverrideModal,
  setShowParamOverrideModal,
  paramOverrideTarget,
  t,
}) => {
  const lines = Array.isArray(paramOverrideTarget?.lines)
    ? paramOverrideTarget.lines
    : [];

  const parsedLines = useMemo(() => {
    return lines.map(parseAuditLine);
  }, [lines]);

  const copyAll = async () => {
    const content = lines.join('\n');
    if (!content) {
      return;
    }
    if (await copy(content)) {
      showSuccess(t('参数覆盖已复制'));
      return;
    }
    showError(t('无法复制到剪贴板，请手动复制'));
  };

  return (
    <Modal
      title={t('参数覆盖详情')}
      visible={showParamOverrideModal}
      onCancel={() => setShowParamOverrideModal(false)}
      footer={null}
      centered
      closable
      maskClosable
      width={760}
    >
      <div style={{ padding: 20 }}>
        <div
          style={{
            marginBottom: 16,
            padding: 16,
            borderRadius: 16,
            background:
              'linear-gradient(135deg, rgba(var(--semi-blue-5), 0.08), rgba(var(--semi-teal-5), 0.12))',
            border: '1px solid rgba(var(--semi-blue-5), 0.16)',
          }}
        >
          <div
            style={{
              display: 'flex',
              justifyContent: 'space-between',
              gap: 12,
              flexWrap: 'wrap',
              alignItems: 'flex-start',
            }}
          >
            <div>
              <div
                style={{
                  fontSize: 18,
                  fontWeight: 700,
                  color: 'var(--semi-color-text-0)',
                  marginBottom: 8,
                }}
              >
                {t('已应用参数覆盖')}
              </div>
              <Space wrap spacing={8}>
                <Tag color='blue' size='large'>
                  {t('{{count}} 项变更', { count: lines.length })}
                </Tag>
                {paramOverrideTarget?.modelName ? (
                  <Tag color='cyan' size='large'>
                    {paramOverrideTarget.modelName}
                  </Tag>
                ) : null}
                {paramOverrideTarget?.requestId ? (
                  <Tag color='grey' size='large'>
                    {t('Request ID')}: {paramOverrideTarget.requestId}
                  </Tag>
                ) : null}
              </Space>
            </div>

            <Button
              icon={<IconCopy />}
              theme='solid'
              type='tertiary'
              onClick={copyAll}
              disabled={lines.length === 0}
            >
              {t('复制全部')}
            </Button>
          </div>

          {paramOverrideTarget?.requestPath ? (
            <div style={{ marginTop: 12 }}>
              <Text type='tertiary' size='small'>
                {t('请求路径')}: {paramOverrideTarget.requestPath}
              </Text>
            </div>
          ) : null}
        </div>

        {lines.length === 0 ? (
          <Empty
            description={t('暂无参数覆盖记录')}
            style={{ padding: '32px 0 12px' }}
          />
        ) : (
          <div
            style={{
              display: 'flex',
              flexDirection: 'column',
              gap: 12,
              maxHeight: '60vh',
              overflowY: 'auto',
              paddingRight: 4,
            }}
          >
            {parsedLines.map((item, index) => {
              if (!item) {
                return (
                  <div
                    key={`raw-${index}`}
                    style={{
                      padding: 14,
                      borderRadius: 14,
                      border: '1px solid var(--semi-color-border)',
                      background: 'var(--semi-color-fill-0)',
                    }}
                  >
                    <Text
                      style={{
                        fontFamily:
                          'ui-monospace, SFMono-Regular, SF Mono, Menlo, Consolas, monospace',
                        fontSize: 13,
                        lineHeight: 1.65,
                        whiteSpace: 'pre-wrap',
                        wordBreak: 'break-word',
                      }}
                    >
                      {lines[index]}
                    </Text>
                  </div>
                );
              }

              return (
                <div
                  key={`${item.field}-${index}`}
                  style={{
                    padding: 14,
                    borderRadius: 16,
                    border: '1px solid var(--semi-color-border)',
                    background: 'var(--semi-color-bg-0)',
                    boxShadow: '0 8px 24px rgba(15, 23, 42, 0.04)',
                  }}
                >
                  <div style={{ marginBottom: 12 }}>
                    <Tag color='blue' shape='circle' size='large'>
                      {item.field}
                    </Tag>
                  </div>
                  <div
                    style={{
                      display: 'flex',
                      gap: 12,
                      flexWrap: 'wrap',
                      alignItems: 'stretch',
                    }}
                  >
                    <ValuePanel label={t('变更前')} value={item.before} />
                    <ValuePanel label={t('变更后')} value={item.after} tone='after' />
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </Modal>
  );
};

export default ParamOverrideModal;
