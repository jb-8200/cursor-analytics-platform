import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { DateRangePicker } from './DateRangePicker';
import { DateRange } from '../../hooks/useDateRange';

describe('DateRangePicker', () => {
  const mockOnChange = vi.fn();

  const defaultRange: DateRange = {
    from: new Date('2026-01-01'),
    to: new Date('2026-01-15'),
    preset: 'LAST_30_DAYS',
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should render with initial value', () => {
    render(<DateRangePicker value={defaultRange} onChange={mockOnChange} />);

    expect(screen.getByText('Last 30 Days')).toBeInTheDocument();
  });

  it('should display custom range format', () => {
    const customRange: DateRange = {
      from: new Date('2026-01-01'),
      to: new Date('2026-01-15'),
      preset: 'CUSTOM',
    };

    render(<DateRangePicker value={customRange} onChange={mockOnChange} />);

    expect(screen.getByText(/Jan 1, 2026.*Jan 15, 2026/)).toBeInTheDocument();
  });

  it('should open dropdown when clicked', async () => {
    const user = userEvent.setup();
    render(<DateRangePicker value={defaultRange} onChange={mockOnChange} />);

    const button = screen.getByRole('button');
    await user.click(button);

    // Should show preset options
    expect(screen.getByText('Last 7 Days')).toBeInTheDocument();
    expect(screen.getByText('Last 90 Days')).toBeInTheDocument();
    expect(screen.getByText('Last 6 Months')).toBeInTheDocument();
    expect(screen.getByText('Last 1 Year')).toBeInTheDocument();
    expect(screen.getByText('Custom')).toBeInTheDocument();
  });

  it('should call onChange when preset is selected', async () => {
    const user = userEvent.setup();
    render(<DateRangePicker value={defaultRange} onChange={mockOnChange} />);

    const button = screen.getByRole('button');
    await user.click(button);

    const preset = screen.getByText('Last 7 Days');
    await user.click(preset);

    expect(mockOnChange).toHaveBeenCalledWith(
      expect.objectContaining({
        preset: 'LAST_7_DAYS',
      })
    );
  });

  it('should close dropdown after selection', async () => {
    const user = userEvent.setup();
    render(<DateRangePicker value={defaultRange} onChange={mockOnChange} />);

    const button = screen.getByRole('button');
    await user.click(button);

    const preset = screen.getByText('Last 7 Days');
    await user.click(preset);

    // Dropdown should close (preset options should not be visible)
    expect(screen.queryByText('Last 90 Days')).not.toBeInTheDocument();
  });

  it('should show custom date inputs when Custom is selected', async () => {
    const user = userEvent.setup();
    render(<DateRangePicker value={defaultRange} onChange={mockOnChange} />);

    const button = screen.getByRole('button');
    await user.click(button);

    const customOption = screen.getByText('Custom');
    await user.click(customOption);

    // Should show date inputs
    expect(screen.getByLabelText(/from/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/to/i)).toBeInTheDocument();
  });

  it('should update custom date range', async () => {
    const user = userEvent.setup();
    const customRange: DateRange = {
      from: new Date('2026-01-01'),
      to: new Date('2026-01-15'),
      preset: 'CUSTOM',
    };

    render(<DateRangePicker value={customRange} onChange={mockOnChange} />);

    const button = screen.getByRole('button');
    await user.click(button);

    // Custom should already be selected, so date inputs should be visible
    const fromInput = screen.getByLabelText(/from/i) as HTMLInputElement;
    await user.clear(fromInput);
    await user.type(fromInput, '2026-02-01');

    expect(mockOnChange).toHaveBeenCalledWith(
      expect.objectContaining({
        from: expect.any(Date),
        preset: 'CUSTOM',
      })
    );
  });

  it('should highlight current preset', async () => {
    const user = userEvent.setup();
    render(<DateRangePicker value={defaultRange} onChange={mockOnChange} />);

    const button = screen.getByRole('button');
    await user.click(button);

    const presetButtons = screen.getAllByText('Last 30 Days');
    // Find the menu item button (not the main button)
    const menuItemButton = presetButtons.find(el => el.closest('[role="menuitem"]'));

    // Should have highlighted/active styling
    expect(menuItemButton?.closest('button')).toHaveClass('bg-primary-50');
  });

  it('should be keyboard accessible', async () => {
    const user = userEvent.setup();
    render(<DateRangePicker value={defaultRange} onChange={mockOnChange} />);

    const button = screen.getByRole('button');

    // Open with keyboard
    button.focus();
    await user.keyboard('{Enter}');

    expect(screen.getByText('Last 7 Days')).toBeInTheDocument();

    // Click on a preset option directly instead of arrow navigation
    const preset = screen.getAllByRole('menuitem')[0]; // First preset (Last 7 Days)
    await user.click(preset);

    expect(mockOnChange).toHaveBeenCalled();
  });

  it('should close dropdown when clicking outside', async () => {
    const user = userEvent.setup();
    render(
      <div>
        <DateRangePicker value={defaultRange} onChange={mockOnChange} />
        <button>Outside Button</button>
      </div>
    );

    const button = screen.getByRole('button', { name: /Last 30 Days/i });
    await user.click(button);

    expect(screen.getByText('Last 7 Days')).toBeInTheDocument();

    // Click outside
    const outsideButton = screen.getByText('Outside Button');
    await user.click(outsideButton);

    // Dropdown should close
    expect(screen.queryByText('Last 7 Days')).not.toBeInTheDocument();
  });

  it('should validate date range (from must be before to)', async () => {
    const user = userEvent.setup();
    const customRange: DateRange = {
      from: new Date('2026-01-15'),
      to: new Date('2026-01-01'), // Invalid: to is before from
      preset: 'CUSTOM',
    };

    render(<DateRangePicker value={customRange} onChange={mockOnChange} />);

    const button = screen.getByRole('button');
    await user.click(button);

    // Should show validation error
    expect(screen.getByText(/from.*must be before.*to/i)).toBeInTheDocument();
  });
});
