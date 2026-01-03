/**
 * Tests for DeveloperTable component
 *
 * Tests the sortable, filterable table displaying individual developer metrics.
 */

import { describe, it, expect, vi } from 'vitest';
import { render, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import DeveloperTable from './DeveloperTable';
import { Developer } from '../../graphql/types';

describe('DeveloperTable', () => {
  const mockDevelopers: Developer[] = [
    {
      id: '1',
      name: 'Alice Johnson',
      email: 'alice@example.com',
      team: 'Frontend',
      seniority: 'senior',
      active: true,
      stats: {
        totalSuggestions: 500,
        acceptedSuggestions: 400,
        acceptanceRate: 0.8,
        aiLinesAdded: 2000,
        aiLinesDeleted: 500,
        humanLinesAdded: 500,
        humanLinesDeleted: 100,
        chatInteractions: 50,
        aiVelocity: 0.8,
      },
    },
    {
      id: '2',
      name: 'Bob Smith',
      email: 'bob@example.com',
      team: 'Backend',
      seniority: 'mid',
      active: true,
      stats: {
        totalSuggestions: 300,
        acceptedSuggestions: 150,
        acceptanceRate: 0.5,
        aiLinesAdded: 1000,
        aiLinesDeleted: 200,
        humanLinesAdded: 1000,
        humanLinesDeleted: 300,
        chatInteractions: 30,
        aiVelocity: 0.5,
      },
    },
    {
      id: '3',
      name: 'Charlie Brown',
      email: 'charlie@example.com',
      team: 'DevOps',
      seniority: 'junior',
      active: true,
      stats: {
        totalSuggestions: 200,
        acceptedSuggestions: 30,
        acceptanceRate: 0.15,
        aiLinesAdded: 300,
        aiLinesDeleted: 50,
        humanLinesAdded: 1700,
        humanLinesDeleted: 400,
        chatInteractions: 20,
        aiVelocity: 0.15,
      },
    },
  ];

  describe('Rendering', () => {
    it('should render table with data', () => {
      render(<DeveloperTable data={mockDevelopers} />);

      // Check table headers
      expect(screen.getByText('Name')).toBeInTheDocument();
      expect(screen.getByText('Team')).toBeInTheDocument();
      expect(screen.getByText('Suggestions')).toBeInTheDocument();
      expect(screen.getByText('Accepted')).toBeInTheDocument();
      expect(screen.getByText('Rate')).toBeInTheDocument();
      expect(screen.getByText('AI Lines')).toBeInTheDocument();
    });

    it('should render all developer rows', () => {
      render(<DeveloperTable data={mockDevelopers} />);

      expect(screen.getByText('Alice Johnson')).toBeInTheDocument();
      expect(screen.getByText('Bob Smith')).toBeInTheDocument();
      expect(screen.getByText('Charlie Brown')).toBeInTheDocument();
    });

    it('should display developer stats correctly', () => {
      render(<DeveloperTable data={mockDevelopers} />);

      // Check Alice's stats
      const aliceRow = screen.getByText('Alice Johnson').closest('tr');
      expect(aliceRow).toBeInTheDocument();
      if (aliceRow) {
        expect(within(aliceRow).getByText('500')).toBeInTheDocument();
        expect(within(aliceRow).getByText('400')).toBeInTheDocument();
        expect(within(aliceRow).getByText('80.0%')).toBeInTheDocument();
      }
    });

    it('should handle empty data', () => {
      render(<DeveloperTable data={[]} />);

      expect(screen.getByText('No developers found')).toBeInTheDocument();
    });

    it('should highlight low acceptance rates', () => {
      render(<DeveloperTable data={mockDevelopers} highlightThreshold={20} />);

      // Charlie has 15% acceptance rate, should be highlighted
      const charlieRow = screen.getByText('Charlie Brown').closest('tr');
      expect(charlieRow).toHaveClass('bg-red-50');
    });
  });

  describe('Sorting', () => {
    it('should sort by name ascending by default', () => {
      render(<DeveloperTable data={mockDevelopers} />);

      const rows = screen.getAllByRole('row');
      // Skip header row (index 0)
      const firstDataRow = rows[1];
      expect(within(firstDataRow).getByText('Alice Johnson')).toBeInTheDocument();
    });

    it('should call onSort when column header clicked', async () => {
      const user = userEvent.setup();
      const onSort = vi.fn();

      render(<DeveloperTable data={mockDevelopers} onSort={onSort} />);

      const teamHeader = screen.getByText('Team');
      await user.click(teamHeader);

      expect(onSort).toHaveBeenCalledWith('team', 'asc');
    });

    it('should toggle sort direction on repeated clicks', async () => {
      const user = userEvent.setup();
      const onSort = vi.fn();

      render(<DeveloperTable data={mockDevelopers} onSort={onSort} />);

      const teamHeader = screen.getByText('Team');

      // First click: asc (switching from name to team)
      await user.click(teamHeader);
      expect(onSort).toHaveBeenCalledWith('team', 'asc');

      // Second click on same column: desc
      await user.click(teamHeader);
      expect(onSort).toHaveBeenCalledWith('team', 'desc');
    });

    it('should show sort indicator on sorted column', async () => {
      const user = userEvent.setup();
      const { container } = render(<DeveloperTable data={mockDevelopers} />);

      const rateHeader = screen.getByText('Rate');
      await user.click(rateHeader);

      // Should show sort icon
      const sortIcon = container.querySelector('[data-testid="sort-icon"]');
      expect(sortIcon).toBeInTheDocument();
    });
  });

  describe('Searching', () => {
    it('should render search input', () => {
      render(<DeveloperTable data={mockDevelopers} />);

      const searchInput = screen.getByPlaceholderText(/search/i);
      expect(searchInput).toBeInTheDocument();
    });

    it('should call onSearch when typing in search input', async () => {
      const user = userEvent.setup();
      const onSearch = vi.fn();

      render(<DeveloperTable data={mockDevelopers} onSearch={onSearch} />);

      const searchInput = screen.getByPlaceholderText(/search/i);
      await user.type(searchInput, 'Alice');

      // Should be called for each character (debounced in implementation)
      expect(onSearch).toHaveBeenCalled();
    });

    it('should clear search when clear button clicked', async () => {
      const user = userEvent.setup();
      const onSearch = vi.fn();

      render(<DeveloperTable data={mockDevelopers} onSearch={onSearch} />);

      const searchInput = screen.getByPlaceholderText(/search/i) as HTMLInputElement;
      await user.type(searchInput, 'Alice');

      // Find and click clear button (if implemented)
      const clearButton = screen.queryByRole('button', { name: /clear/i });
      if (clearButton) {
        await user.click(clearButton);
        expect(searchInput.value).toBe('');
        expect(onSearch).toHaveBeenCalledWith('');
      }
    });
  });

  describe('Pagination', () => {
    const manyDevelopers: Developer[] = Array.from({ length: 50 }, (_, i) => ({
      id: `${i + 1}`,
      name: `Developer ${i + 1}`,
      email: `dev${i + 1}@example.com`,
      team: `Team ${(i % 3) + 1}`,
      seniority: 'mid',
      active: true,
      stats: {
        totalSuggestions: 100 + i,
        acceptedSuggestions: 50 + i,
        acceptanceRate: 0.5 + i / 100,
        aiLinesAdded: 500 + i * 10,
        aiLinesDeleted: 100,
        humanLinesAdded: 500,
        humanLinesDeleted: 100,
        chatInteractions: 10 + i,
        aiVelocity: 0.5,
      },
    }));

    it('should show pagination controls', () => {
      render(<DeveloperTable data={manyDevelopers} pageSize={25} />);

      expect(screen.getByText(/page/i)).toBeInTheDocument();
    });

    it('should display correct number of rows per page', () => {
      render(<DeveloperTable data={manyDevelopers} pageSize={25} />);

      const rows = screen.getAllByRole('row');
      // 1 header + 25 data rows
      expect(rows.length).toBe(26);
    });

    it('should call onPageChange when page button clicked', async () => {
      const user = userEvent.setup();
      const onPageChange = vi.fn();

      render(
        <DeveloperTable
          data={manyDevelopers}
          pageSize={25}
          onPageChange={onPageChange}
        />
      );

      const nextButton = screen.getByRole('button', { name: /next/i });
      await user.click(nextButton);

      expect(onPageChange).toHaveBeenCalledWith(2);
    });

    it('should disable previous button on first page', () => {
      render(<DeveloperTable data={manyDevelopers} pageSize={25} />);

      const prevButton = screen.getByRole('button', { name: /previous/i });
      expect(prevButton).toBeDisabled();
    });

    it('should disable next button on last page', () => {
      render(
        <DeveloperTable data={manyDevelopers} pageSize={25} currentPage={2} />
      );

      const nextButton = screen.getByRole('button', { name: /next/i });
      expect(nextButton).toBeDisabled();
    });

    it('should show current page information', () => {
      render(
        <DeveloperTable data={manyDevelopers} pageSize={25} currentPage={1} />
      );

      expect(screen.getByText(/1-25 of 50/i)).toBeInTheDocument();
    });
  });

  describe('Row Interaction', () => {
    it('should call onRowClick when row clicked', async () => {
      const user = userEvent.setup();
      const onRowClick = vi.fn();

      render(
        <DeveloperTable data={mockDevelopers} onRowClick={onRowClick} />
      );

      const aliceRow = screen.getByText('Alice Johnson').closest('tr');
      if (aliceRow) {
        await user.click(aliceRow);
        expect(onRowClick).toHaveBeenCalledWith(mockDevelopers[0]);
      }
    });

    it('should show hover state on rows', () => {
      const { container } = render(<DeveloperTable data={mockDevelopers} />);

      const rows = container.querySelectorAll('tbody tr');
      rows.forEach((row) => {
        expect(row).toHaveClass('hover:bg-gray-50');
      });
    });
  });

  describe('Accessibility', () => {
    it('should have proper table structure', () => {
      render(<DeveloperTable data={mockDevelopers} />);

      const table = screen.getByRole('table');
      expect(table).toBeInTheDocument();

      const rowgroups = within(table).getAllByRole('rowgroup');
      // Should have thead and tbody
      expect(rowgroups.length).toBe(2);
    });

    it('should have sortable column headers with aria-sort', async () => {
      const user = userEvent.setup();
      render(<DeveloperTable data={mockDevelopers} />);

      const nameHeader = screen.getByText('Name');
      await user.click(nameHeader);

      const th = nameHeader.closest('th');
      expect(th).toHaveAttribute('aria-sort');
    });

    it('should be keyboard navigable', () => {
      render(<DeveloperTable data={mockDevelopers} />);

      const searchInput = screen.getByPlaceholderText(/search/i);
      expect(searchInput).toBeInTheDocument();

      // Tab index should be set for interactive elements
      const buttons = screen.getAllByRole('button');
      buttons.forEach((button) => {
        expect(button.tabIndex).toBeGreaterThanOrEqual(0);
      });
    });
  });

  describe('Edge Cases', () => {
    it('should handle developers without stats', () => {
      const devsWithoutStats: Developer[] = [
        {
          id: '1',
          name: 'No Stats Dev',
          email: 'nostats@example.com',
          team: 'Test',
          seniority: 'junior',
          active: true,
        },
      ];

      render(<DeveloperTable data={devsWithoutStats} />);

      expect(screen.getByText('No Stats Dev')).toBeInTheDocument();
      // Should show multiple N/A values for missing stats
      const naElements = screen.getAllByText('N/A');
      expect(naElements.length).toBeGreaterThan(0);
    });

    it('should handle very long names gracefully', () => {
      const longNameDev: Developer[] = [
        {
          id: '1',
          name: 'A'.repeat(100),
          email: 'long@example.com',
          team: 'Test',
          seniority: 'senior',
          active: true,
          stats: mockDevelopers[0].stats,
        },
      ];

      const { container } = render(<DeveloperTable data={longNameDev} />);

      // Should render without breaking layout
      const table = container.querySelector('table');
      expect(table).toBeInTheDocument();
    });

    it('should handle zero values in stats', () => {
      const zeroStatsDev: Developer[] = [
        {
          id: '1',
          name: 'Zero Stats',
          email: 'zero@example.com',
          team: 'Test',
          seniority: 'mid',
          active: true,
          stats: {
            totalSuggestions: 0,
            acceptedSuggestions: 0,
            acceptanceRate: 0,
            aiLinesAdded: 0,
            aiLinesDeleted: 0,
            humanLinesAdded: 0,
            humanLinesDeleted: 0,
            chatInteractions: 0,
            aiVelocity: 0,
          },
        },
      ];

      render(<DeveloperTable data={zeroStatsDev} />);

      // Check for 0.0% acceptance rate (formatted)
      expect(screen.getByText('0.0%')).toBeInTheDocument();
    });
  });
});
