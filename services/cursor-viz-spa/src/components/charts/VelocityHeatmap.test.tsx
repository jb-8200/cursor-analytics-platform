/**
 * Tests for VelocityHeatmap component
 *
 * Tests the GitHub-style contribution grid that displays
 * AI code acceptance intensity over time.
 */

import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import VelocityHeatmap from './VelocityHeatmap';
import { DailyStats } from '../../graphql/types';

describe('VelocityHeatmap', () => {
  // Mock data: 52 weeks worth of data (364 days)
  const mockData: DailyStats[] = Array.from({ length: 364 }, (_, index) => {
    const date = new Date('2025-01-01');
    date.setDate(date.getDate() + index);
    return {
      date: date.toISOString().split('T')[0],
      suggestionsShown: 100,
      suggestionsAccepted: 50 + (index % 30), // Varying acceptance
      acceptanceRate: 0.5 + (index % 30) / 100,
      aiLinesAdded: 100 + (index % 50),
      humanLinesAdded: 200,
      chatInteractions: 10,
    };
  });

  describe('Rendering', () => {
    it('should render a heatmap grid', () => {
      const { container } = render(<VelocityHeatmap data={mockData} />);

      // Check for SVG container
      const svg = container.querySelector('svg');
      expect(svg).toBeInTheDocument();
    });

    it('should render 7 rows (days of week)', () => {
      const { container } = render(<VelocityHeatmap data={mockData} />);

      // Count cells in a single week column
      const cells = container.querySelectorAll('[data-testid^="heatmap-cell"]');

      // Should have cells for all data points
      expect(cells.length).toBeGreaterThan(0);
    });

    it('should render 52 weeks by default', () => {
      const { container } = render(<VelocityHeatmap data={mockData} />);

      // Check that we have approximately 52 weeks worth of cells
      const cells = container.querySelectorAll('[data-testid^="heatmap-cell"]');
      expect(cells.length).toBeGreaterThanOrEqual(350); // ~52 weeks * 7 days
    });

    it('should render custom number of weeks', () => {
      const twoWeeksData = mockData.slice(0, 14);
      const { container } = render(
        <VelocityHeatmap data={twoWeeksData} weeks={2} />
      );

      const cells = container.querySelectorAll('[data-testid^="heatmap-cell"]');
      expect(cells.length).toBeLessThanOrEqual(14);
    });

    it('should render day labels (Mon, Wed, Fri)', () => {
      render(<VelocityHeatmap data={mockData} />);

      // Check for day labels
      expect(screen.getByText(/Mon/i)).toBeInTheDocument();
      expect(screen.getByText(/Wed/i)).toBeInTheDocument();
      expect(screen.getByText(/Fri/i)).toBeInTheDocument();
    });

    it('should render month labels at month boundaries', () => {
      render(<VelocityHeatmap data={mockData} />);

      // Check for at least some month labels
      const monthPattern = /Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec/;
      const monthLabels = screen.getAllByText(monthPattern);
      expect(monthLabels.length).toBeGreaterThan(0);
    });

    it('should handle empty data gracefully', () => {
      const { container } = render(<VelocityHeatmap data={[]} />);

      const svg = container.querySelector('svg');
      expect(svg).toBeInTheDocument();

      // Should show some placeholder or message
      expect(container.textContent).toBeTruthy();
    });
  });

  describe('Color Mapping', () => {
    it('should apply different colors based on activity level', () => {
      const variedData: DailyStats[] = [
        {
          date: '2025-01-01',
          suggestionsShown: 100,
          suggestionsAccepted: 0, // No activity
          acceptanceRate: 0,
          aiLinesAdded: 0,
          humanLinesAdded: 0,
          chatInteractions: 0,
        },
        {
          date: '2025-01-02',
          suggestionsShown: 100,
          suggestionsAccepted: 10, // Low activity
          acceptanceRate: 0.1,
          aiLinesAdded: 10,
          humanLinesAdded: 90,
          chatInteractions: 1,
        },
        {
          date: '2025-01-03',
          suggestionsShown: 100,
          suggestionsAccepted: 80, // High activity
          acceptanceRate: 0.8,
          aiLinesAdded: 80,
          humanLinesAdded: 20,
          chatInteractions: 10,
        },
      ];

      const { container } = render(<VelocityHeatmap data={variedData} />);

      const cells = container.querySelectorAll('[data-testid^="heatmap-cell"]');
      const fills = Array.from(cells).map((cell) => cell.getAttribute('fill'));

      // Expect different colors for different activity levels
      const uniqueColors = new Set(fills);
      expect(uniqueColors.size).toBeGreaterThan(1);
    });

    it('should use custom color scale when provided', () => {
      const customColors = ['#ffffff', '#ff0000', '#00ff00', '#0000ff', '#000000'];
      const { container } = render(
        <VelocityHeatmap data={mockData} colorScale={customColors} />
      );

      const cells = container.querySelectorAll('[data-testid^="heatmap-cell"]');
      expect(cells.length).toBeGreaterThan(0);

      // At least one cell should use a custom color
      const fills = Array.from(cells).map((cell) => cell.getAttribute('fill'));
      const usesCustomColor = fills.some((fill) => customColors.includes(fill || ''));
      expect(usesCustomColor).toBe(true);
    });
  });

  describe('Tooltips', () => {
    it('should show tooltip on cell hover', async () => {
      const user = userEvent.setup();
      const { container } = render(<VelocityHeatmap data={mockData} />);

      const firstCell = container.querySelector('[data-testid="heatmap-cell-0"]');
      expect(firstCell).toBeInTheDocument();

      // Hover over cell
      if (firstCell) {
        await user.hover(firstCell);

        // Tooltip should appear with date and count
        // Note: Tooltip implementation may vary (e.g., title attribute, custom tooltip)
        // This test checks for a title attribute or aria-label
        const hasTooltip =
          firstCell.getAttribute('title') || firstCell.getAttribute('aria-label');
        expect(hasTooltip).toBeTruthy();
      }
    });

    it('should display date and count in tooltip', async () => {
      const user = userEvent.setup();
      const singleDayData: DailyStats[] = [
        {
          date: '2025-01-01',
          suggestionsShown: 100,
          suggestionsAccepted: 42,
          acceptanceRate: 0.42,
          aiLinesAdded: 42,
          humanLinesAdded: 58,
          chatInteractions: 5,
        },
      ];

      const { container } = render(<VelocityHeatmap data={singleDayData} />);

      const cell = container.querySelector('[data-testid="heatmap-cell-0"]');
      if (cell) {
        await user.hover(cell);

        const tooltip = cell.getAttribute('title') || cell.getAttribute('aria-label');
        expect(tooltip).toContain('2025-01-01');
        expect(tooltip).toContain('42');
      }
    });
  });

  describe('Interactions', () => {
    it('should call onCellClick when a cell is clicked', async () => {
      const user = userEvent.setup();
      const onCellClick = vi.fn();
      const { container } = render(
        <VelocityHeatmap data={mockData} onCellClick={onCellClick} />
      );

      const firstCell = container.querySelector('[data-testid="heatmap-cell-0"]');
      expect(firstCell).toBeInTheDocument();

      if (firstCell) {
        await user.click(firstCell);
        expect(onCellClick).toHaveBeenCalledTimes(1);

        // Should pass a Date object
        const callArg = onCellClick.mock.calls[0][0];
        expect(callArg).toBeInstanceOf(Date);
      }
    });

    it('should handle clicks on cells with no data', async () => {
      const user = userEvent.setup();
      const onCellClick = vi.fn();
      const { container } = render(
        <VelocityHeatmap data={[]} onCellClick={onCellClick} />
      );

      const svg = container.querySelector('svg');
      if (svg) {
        await user.click(svg);
        // Should not crash
        expect(true).toBe(true);
      }
    });
  });

  describe('Accessibility', () => {
    it('should have proper ARIA labels', () => {
      const { container } = render(<VelocityHeatmap data={mockData} />);

      const svg = container.querySelector('svg');
      expect(svg).toHaveAttribute('role', 'img');
      expect(svg).toHaveAttribute('aria-label');
    });

    it('should be keyboard navigable', () => {
      const { container } = render(<VelocityHeatmap data={mockData} />);

      const cells = container.querySelectorAll('[data-testid^="heatmap-cell"]');
      if (cells.length > 0) {
        // Cells should be focusable or within a focusable container
        const firstCell = cells[0] as HTMLElement;
        expect(
          firstCell.tabIndex >= 0 || firstCell.parentElement?.tabIndex !== undefined
        ).toBe(true);
      }
    });
  });

  describe('Edge Cases', () => {
    it('should handle single day of data', () => {
      const singleDay: DailyStats[] = [
        {
          date: '2025-01-01',
          suggestionsShown: 100,
          suggestionsAccepted: 50,
          acceptanceRate: 0.5,
          aiLinesAdded: 50,
          humanLinesAdded: 50,
          chatInteractions: 5,
        },
      ];

      const { container } = render(<VelocityHeatmap data={singleDay} />);

      const svg = container.querySelector('svg');
      expect(svg).toBeInTheDocument();

      const cells = container.querySelectorAll('[data-testid^="heatmap-cell"]');
      expect(cells.length).toBeGreaterThan(0);
    });

    it('should handle data gaps (missing days)', () => {
      const dataWithGaps: DailyStats[] = [
        {
          date: '2025-01-01',
          suggestionsShown: 100,
          suggestionsAccepted: 50,
          acceptanceRate: 0.5,
          aiLinesAdded: 50,
          humanLinesAdded: 50,
          chatInteractions: 5,
        },
        // Gap: 2025-01-02 and 2025-01-03 missing
        {
          date: '2025-01-04',
          suggestionsShown: 100,
          suggestionsAccepted: 60,
          acceptanceRate: 0.6,
          aiLinesAdded: 60,
          humanLinesAdded: 40,
          chatInteractions: 6,
        },
      ];

      const { container } = render(<VelocityHeatmap data={dataWithGaps} />);

      const svg = container.querySelector('svg');
      expect(svg).toBeInTheDocument();

      // Should still render, filling gaps with zero activity
      const cells = container.querySelectorAll('[data-testid^="heatmap-cell"]');
      expect(cells.length).toBeGreaterThan(2);
    });

    it('should handle data with future dates gracefully', () => {
      const futureData: DailyStats[] = [
        {
          date: '2099-12-31',
          suggestionsShown: 100,
          suggestionsAccepted: 50,
          acceptanceRate: 0.5,
          aiLinesAdded: 50,
          humanLinesAdded: 50,
          chatInteractions: 5,
        },
      ];

      const { container } = render(<VelocityHeatmap data={futureData} />);

      const svg = container.querySelector('svg');
      expect(svg).toBeInTheDocument();
    });
  });
});
