/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { type FC, useMemo, useState } from 'react';

import { useRequest } from 'ahooks';
import { Button, Empty, Spin, Toast, Typography } from '@coze-arch/coze-design';

import { CRM_QUERY_EXAMPLES } from '../mock-data';
import { getDefaultCRMNLQuestion, runCRMNLQuery } from '../nl-query-service';
import type { CRMNLQueryResult, CRMQueryHistoryItem } from '../types';
import { QueryResultChart } from './query-result-chart';
import { QueryResultTable } from './query-result-table';

interface CRMNLQueryPanelProps {
  spaceId: string;
}

const CARD_CLASS_NAME =
  'rounded-[20px] border border-solid coz-stroke-primary coz-bg-max p-[20px] shadow-[0_16px_40px_rgba(15,23,42,0.04)]';
const TEXTAREA_CLASS_NAME =
  'min-h-[112px] w-full resize-y rounded-[16px] border border-solid coz-stroke-primary bg-white px-[14px] py-[12px] text-[14px] leading-[22px] outline-none transition focus:border-[#2563eb]';

export const CRMNLQueryPanel: FC<CRMNLQueryPanelProps> = ({ spaceId }) => {
  const [question, setQuestion] = useState(getDefaultCRMNLQuestion());
  const [result, setResult] = useState<CRMNLQueryResult>();
  const [errorMessage, setErrorMessage] = useState('');
  const [history, setHistory] = useState<CRMQueryHistoryItem[]>([]);

  const queryRequest = useRequest(
    async (nextQuestion: string) =>
      runCRMNLQuery({
        space_id: spaceId,
        question: nextQuestion,
      }),
    {
      manual: true,
      onSuccess: response => {
        setErrorMessage('');
        setResult(response);
        setHistory(prev => {
          const nextHistory: CRMQueryHistoryItem[] = [
            {
              id: `${Date.now()}`,
              question: response.question,
              answer: response.answer,
              intent_type: response.intent_type,
              created_at: formatNow(),
            },
            ...prev,
          ];

          return nextHistory.slice(0, 6);
        });
      },
      onError: error => {
        const nextError =
          error instanceof Error ? error.message : 'Query failed. Please retry.';
        setErrorMessage(nextError);
        Toast.error(nextError);
      },
    },
  );

  const scopeBadges = useMemo(() => {
    const items = [
      result?.intent_type ? `Intent: ${result.intent_type}` : '',
      result?.meta?.scope_note ? `Scope: ${result.meta.scope_note}` : '',
      result?.meta?.updated_at ? `Updated: ${result.meta.updated_at}` : '',
    ];

    return items.filter(Boolean);
  }, [result]);

  const handleRunQuery = () => {
    const trimmedQuestion = question.trim();

    if (!trimmedQuestion) {
      const nextError =
        'Please enter a CRM question, for example: How many customers do I have now?';
      setErrorMessage(nextError);
      Toast.warning(nextError);
      return;
    }

    queryRequest.run(trimmedQuestion);
  };

  return (
    <section
      id="crm-ai-agent-panel"
      data-testid="crm-ai-agent-panel"
      className={`${CARD_CLASS_NAME} bg-[linear-gradient(135deg,#fffaf0_0%,#ffffff_48%,#e0f2fe_100%)]`}
    >
      <div className="mb-[16px] flex flex-wrap items-start justify-between gap-[16px]">
        <div>
          <Typography.Title heading={5} className="!mb-[6px]">
            CRM Natural Language Query
          </Typography.Title>
          <Typography.Paragraph className="!mb-0 max-w-[820px] coz-fg-secondary">
            Let business users ask CRM questions in natural language and read
            the result as a summary, table, chart and business definition.
          </Typography.Paragraph>
        </div>
        <div className="rounded-full bg-[rgba(255,255,255,0.88)] px-[12px] py-[6px] text-[12px] font-[600] text-[#0f172a]">
          MVP: customer count / top sales / forecast
        </div>
      </div>

      <div className="grid grid-cols-1 gap-[16px] xl:grid-cols-[minmax(0,1.2fr)_minmax(320px,0.8fr)]">
        <div className="rounded-[18px] bg-white/90 p-[16px]">
          <Typography.Text className="mb-[8px] block text-[12px] font-[600] coz-fg-secondary">
            Query Input
          </Typography.Text>
          <textarea
            className={TEXTAREA_CLASS_NAME}
            value={question}
            placeholder="Example: How many customers do I have now?"
            onChange={event => setQuestion(event.target.value)}
            onKeyDown={event => {
              if ((event.metaKey || event.ctrlKey) && event.key === 'Enter') {
                event.preventDefault();
                handleRunQuery();
              }
            }}
          />
          <div className="mt-[12px] flex flex-wrap items-center justify-between gap-[12px]">
            <Typography.Text className="text-[12px] coz-fg-secondary">
              Adding a time range, owner or product dimension usually makes the
              answer more stable.
            </Typography.Text>
            <div className="flex gap-[12px]">
              <Button
                color="secondary"
                htmlType="button"
                onClick={() => {
                  setQuestion(getDefaultCRMNLQuestion());
                  setErrorMessage('');
                }}
              >
                Reset
              </Button>
              <Button
                htmlType="button"
                data-testid="crm-run-query-button"
                loading={queryRequest.loading}
                onClick={handleRunQuery}
              >
                Run Query
              </Button>
            </div>
          </div>
        </div>

        <div className="rounded-[18px] bg-[rgba(248,250,252,0.82)] p-[16px]">
          <Typography.Text className="mb-[8px] block text-[12px] font-[600] coz-fg-secondary">
            Example Questions
          </Typography.Text>
          <div className="flex flex-col gap-[10px]">
            {CRM_QUERY_EXAMPLES.map(item => (
              <button
                key={item.key}
                type="button"
                className="rounded-[16px] border border-solid coz-stroke-primary bg-white px-[14px] py-[12px] text-left transition hover:border-[#93c5fd] hover:bg-[rgba(239,246,255,0.72)]"
                onClick={() => {
                  setQuestion(item.question);
                  setErrorMessage('');
                }}
              >
                <div className="text-[14px] font-[600] text-[#0f172a]">
                  {item.title}
                </div>
                <div className="mt-[4px] text-[13px] text-[#334155]">
                  {item.question}
                </div>
                <div className="mt-[6px] text-[12px] coz-fg-secondary">
                  {item.description}
                </div>
              </button>
            ))}
          </div>
        </div>
      </div>

      {errorMessage ? (
        <div className="mt-[16px] rounded-[16px] border border-solid border-[#fecaca] bg-[rgba(254,242,242,0.9)] px-[16px] py-[12px]">
          <Typography.Text className="text-[13px] text-[#b91c1c]">
            Error: {errorMessage}
          </Typography.Text>
        </div>
      ) : null}

      {history.length ? (
        <div className="mt-[16px] rounded-[18px] bg-[rgba(255,255,255,0.86)] p-[16px]">
          <Typography.Text className="mb-[10px] block text-[12px] font-[600] coz-fg-secondary">
            Recent Queries
          </Typography.Text>
          <div className="grid grid-cols-1 gap-[10px] lg:grid-cols-2">
            {history.map(item => (
              <button
                key={item.id}
                type="button"
                className="rounded-[14px] border border-solid coz-stroke-primary bg-white px-[14px] py-[12px] text-left transition hover:border-[#93c5fd]"
                onClick={() => setQuestion(item.question)}
              >
                <div className="text-[13px] font-[600] text-[#0f172a]">
                  {item.question}
                </div>
                <div className="mt-[6px] text-[12px] coz-fg-secondary">
                  {item.answer}
                </div>
                <div className="mt-[8px] text-[11px] text-[#64748b]">
                  {item.created_at}
                </div>
              </button>
            ))}
          </div>
        </div>
      ) : null}

      <div className="mt-[16px] rounded-[18px] bg-white/92 p-[16px]">
        <div className="mb-[12px] flex flex-wrap items-center justify-between gap-[12px]">
          <div>
            <Typography.Title heading={6} className="!mb-[4px]">
              Result Area
            </Typography.Title>
            <Typography.Paragraph className="!mb-0 coz-fg-secondary">
              The page surfaces the short answer first, then expands into table,
              chart and definition details.
            </Typography.Paragraph>
          </div>
          <div className="flex flex-wrap gap-[8px]">
            {scopeBadges.map(item => (
              <span
                key={item}
                className="rounded-full bg-[rgba(239,246,255,0.9)] px-[10px] py-[5px] text-[12px] text-[#1d4ed8]"
              >
                {item}
              </span>
            ))}
          </div>
        </div>

        <Spin spinning={queryRequest.loading}>
          {result ? (
            <div className="flex flex-col gap-[16px]">
              <div className="rounded-[18px] border border-solid coz-stroke-primary bg-[linear-gradient(135deg,#eff6ff_0%,#ffffff_52%,#dbeafe_100%)] p-[18px]">
                <Typography.Text className="coz-fg-secondary text-[12px]">
                  One-line Conclusion
                </Typography.Text>
                <Typography.Title heading={4} className="!mb-[8px] !mt-[12px]">
                  {result.answer}
                </Typography.Title>
                <Typography.Paragraph className="!mb-0 coz-fg-secondary">
                  {result.meta?.definition ||
                    'The result is interpreted with the default CRM definition.'}
                </Typography.Paragraph>
              </div>

              <div className="grid grid-cols-1 gap-[16px] 2xl:grid-cols-[minmax(0,1.1fr)_minmax(380px,0.9fr)]">
                <QueryResultTable result={result} />
                <QueryResultChart result={result} />
              </div>

              <div className="grid grid-cols-1 gap-[16px] xl:grid-cols-2">
                <div className="rounded-[16px] border border-solid coz-stroke-primary bg-[rgba(248,250,252,0.72)] p-[16px]">
                  <Typography.Title heading={6} className="!mb-[10px]">
                    Query Definition
                  </Typography.Title>
                  <div className="flex flex-col gap-[8px] text-[13px] text-[#334155]">
                    <div>
                      Default definition:{' '}
                      {result.meta?.definition || 'Metric default definition'}
                    </div>
                    <div>
                      Data updated at:{' '}
                      {result.meta?.updated_at || 'After query execution'}
                    </div>
                    <div>
                      Visible scope:{' '}
                      {result.meta?.scope_note || 'Current authorized scope'}
                    </div>
                  </div>
                </div>

                <div className="rounded-[16px] border border-solid coz-stroke-primary bg-[rgba(248,250,252,0.72)] p-[16px]">
                  <Typography.Title heading={6} className="!mb-[10px]">
                    Applied Filters
                  </Typography.Title>
                  {result.meta?.applied_filters?.length ? (
                    <div className="flex flex-wrap gap-[8px]">
                      {result.meta.applied_filters.map(item => (
                        <span
                          key={item}
                          className="rounded-full bg-white px-[10px] py-[5px] text-[12px] text-[#334155]"
                        >
                          {item}
                        </span>
                      ))}
                    </div>
                  ) : (
                    <Typography.Paragraph className="!mb-0 coz-fg-secondary">
                      No explicit filters were returned. The query ran with the
                      metric default scope.
                    </Typography.Paragraph>
                  )}
                </div>
              </div>

              {result.disclaimer ? (
                <div className="rounded-[16px] border border-solid border-[#fed7aa] bg-[rgba(255,247,237,0.92)] px-[16px] py-[14px]">
                  <Typography.Text className="text-[13px] text-[#c2410c]">
                    Disclaimer: {result.disclaimer}
                  </Typography.Text>
                </div>
              ) : null}
            </div>
          ) : (
            <div className="flex min-h-[280px] items-center justify-center rounded-[18px] border border-dashed coz-stroke-primary bg-[rgba(248,250,252,0.72)]">
              <Empty title="Pick an example question or enter your own query to start." />
            </div>
          )}
        </Spin>
      </div>
    </section>
  );
};

function formatNow(): string {
  const now = new Date();
  return `${now.getFullYear()}-${pad(now.getMonth() + 1)}-${pad(now.getDate())} ${pad(now.getHours())}:${pad(now.getMinutes())}`;
}

function pad(value: number): string {
  return String(value).padStart(2, '0');
}
