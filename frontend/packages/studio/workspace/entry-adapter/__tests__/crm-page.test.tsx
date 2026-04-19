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

import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { CRMManagePage } from '../src/pages/crm/page';

const getDashboardOverviewMock = vi.fn();
const listCustomersMock = vi.fn();
const listOpportunitiesMock = vi.fn();
const listSalesOrdersMock = vi.fn();
const createCustomerMock = vi.fn();

vi.mock('@coze-studio/api-schema', () => ({
  crm: {
    GetDashboardOverview: getDashboardOverviewMock,
    ListCustomers: listCustomersMock,
    ListOpportunities: listOpportunitiesMock,
    ListSalesOrders: listSalesOrdersMock,
    CreateCustomer: createCustomerMock,
  },
}));

describe('CRMManagePage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    getDashboardOverviewMock.mockResolvedValue({
      data: {
        customer_total: 36,
        new_customers_this_month: 4,
        opportunity_total_amount: '128000.5',
        new_opportunities_this_month: 7,
        sales_order_total_amount: '86500.25',
        recent_order_trend: [
          { date: '2026-04-10', order_count: 2, order_amount: '16800' },
          { date: '2026-04-11', order_count: 1, order_amount: '9200' },
        ],
      },
    });
    listCustomersMock.mockResolvedValue({
      data: {
        list: [
          {
            customer_id: '1',
            customer_name: 'Star Manufacturing',
            owner_user_name: 'Alice',
            status: 'active',
            updated_at: '1776038400000',
          },
        ],
      },
    });
    listOpportunitiesMock.mockResolvedValue({
      data: {
        list: [
          {
            opportunity_id: '2',
            opportunity_name: 'East Region Renewal',
            stage: 'proposal',
            amount: '52000',
            status: 'open',
          },
        ],
      },
    });
    listSalesOrdersMock.mockResolvedValue({
      data: {
        list: [
          {
            sales_order_id: '3',
            product_name: 'AI Seat Package',
            amount: '32000',
            order_date: '2026-04-12',
            status: 'draft',
          },
        ],
      },
    });
    createCustomerMock.mockResolvedValue({
      data: {
        customer_id: '99',
        customer_name: 'New Customer',
      },
    });
  });

  it('renders dashboard metrics, trend section and summary tables', async () => {
    render(<CRMManagePage spaceId="1" />);

    await waitFor(() => {
      expect(getDashboardOverviewMock).toHaveBeenCalledWith(
        { space_id: '1' },
        { __disableErrorToast: true },
      );
    });

    expect(await screen.findByText('CRM Dashboard')).toBeInTheDocument();
    expect(
      screen.getByRole('button', { name: 'Overview' }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole('button', { name: 'Customers' }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole('button', { name: 'Opportunities' }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole('button', { name: 'Orders' }),
    ).toBeInTheDocument();
    expect(screen.getByText('Customers')).toBeInTheDocument();
    expect(screen.getByText('36')).toBeInTheDocument();
    expect(screen.getByText('Recent 30-Day Order Trend')).toBeInTheDocument();
    expect(screen.getByText('Customer Snapshot')).toBeInTheDocument();
    expect(screen.getByText('Star Manufacturing')).toBeInTheDocument();
    expect(screen.getByText('East Region Renewal')).toBeInTheDocument();
    expect(screen.getByText('AI Seat Package')).toBeInTheDocument();
  });

  it('switches crm sub menu views', async () => {
    render(<CRMManagePage spaceId="1" />);

    await screen.findByText('CRM Dashboard');

    fireEvent.click(screen.getByRole('button', { name: 'Customers' }));

    expect(await screen.findByText('Customer Management')).toBeInTheDocument();
    expect(screen.getByText('Customer Snapshot')).toBeInTheDocument();
    expect(screen.queryByText('Opportunity Pipeline')).not.toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'Orders' }));

    expect(await screen.findByText('Sales Orders')).toBeInTheDocument();
    expect(screen.getByText('Order Snapshot')).toBeInTheDocument();
    expect(screen.getByText('Recent 30-Day Order Trend')).toBeInTheDocument();
  });

  it('opens quick create form and submits the minimal customer payload', async () => {
    render(<CRMManagePage spaceId="1" />);

    await screen.findByText('CRM Dashboard');

    fireEvent.click(screen.getByRole('button', { name: 'Quick Create Customer' }));

    fireEvent.change(screen.getByLabelText('Customer Name'), {
      target: { value: 'New Customer' },
    });
    fireEvent.change(screen.getByLabelText('Owner'), {
      target: { value: 'Bob' },
    });
    fireEvent.change(screen.getByLabelText('Industry'), {
      target: { value: 'AI SaaS' },
    });
    fireEvent.change(screen.getByLabelText('Remark'), {
      target: { value: 'seeded from dashboard form' },
    });

    fireEvent.click(screen.getByRole('button', { name: 'Create Customer' }));

    await waitFor(() => {
      expect(createCustomerMock).toHaveBeenCalledWith(
        {
          space_id: '1',
          customer_name: 'New Customer',
          owner_user_name: 'Bob',
          industry: 'AI SaaS',
          remark: 'seeded from dashboard form',
          status: 'active',
        },
        { __disableErrorToast: true },
      );
    });

    await waitFor(() => {
      expect(listCustomersMock).toHaveBeenCalledTimes(2);
      expect(getDashboardOverviewMock).toHaveBeenCalledTimes(2);
    });
  });

  it('runs the crm nl query panel and renders the mock result', async () => {
    render(<CRMManagePage spaceId="1" />);

    await screen.findByText('CRM Dashboard');

    fireEvent.click(screen.getByTestId('crm-run-query-button'));

    expect(
      await screen.findByText(
        'There are 1,277 active customers in the current workspace.',
      ),
    ).toBeInTheDocument();
    expect(screen.getByText('Structured Result')).toBeInTheDocument();
    expect(screen.getByText('Active Customer Count')).toBeInTheDocument();
  });
});
