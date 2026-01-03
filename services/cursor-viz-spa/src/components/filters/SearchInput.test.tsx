import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { SearchInput } from './SearchInput';

describe('SearchInput', () => {
  const mockOnChange = vi.fn();

  beforeEach(() => {
    vi.useFakeTimers();
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('should render with placeholder', () => {
    render(<SearchInput value="" onChange={mockOnChange} placeholder="Search developers..." />);

    expect(screen.getByPlaceholderText('Search developers...')).toBeInTheDocument();
  });

  it('should display initial value', () => {
    render(<SearchInput value="John" onChange={mockOnChange} />);

    const input = screen.getByRole('textbox') as HTMLInputElement;
    expect(input.value).toBe('John');
  });

  it('should debounce onChange calls', async () => {
    const user = userEvent.setup({ delay: null }); // Disable userEvent delay for timer control
    render(<SearchInput value="" onChange={mockOnChange} debounceMs={300} />);

    const input = screen.getByRole('textbox');

    // Type quickly
    await user.type(input, 'John');

    // Should not call onChange immediately
    expect(mockOnChange).not.toHaveBeenCalled();

    // Fast-forward time
    act(() => {
      vi.runAllTimers();
    });

    // Now it should be called once with final value
    expect(mockOnChange).toHaveBeenCalledTimes(1);
    expect(mockOnChange).toHaveBeenCalledWith('John');
  });

  it('should use default debounce of 300ms', async () => {
    const user = userEvent.setup({ delay: null });
    render(<SearchInput value="" onChange={mockOnChange} />);

    const input = screen.getByRole('textbox');
    await user.type(input, 'Test');

    // Not called before 300ms
    act(() => {
      vi.advanceTimersByTime(299);
    });
    expect(mockOnChange).not.toHaveBeenCalled();

    // Called after 300ms
    act(() => {
      vi.advanceTimersByTime(1);
    });
    expect(mockOnChange).toHaveBeenCalledTimes(1);
  });

  it('should cancel previous debounce on new input', async () => {
    const user = userEvent.setup({ delay: null });
    render(<SearchInput value="" onChange={mockOnChange} debounceMs={300} />);

    const input = screen.getByRole('textbox');

    await user.type(input, 'Jo');
    act(() => {
      vi.advanceTimersByTime(100);
    });

    await user.type(input, 'hn');
    act(() => {
      vi.advanceTimersByTime(100);
    });

    // Should not have called yet
    expect(mockOnChange).not.toHaveBeenCalled();

    // Complete the debounce
    act(() => {
      vi.runAllTimers();
    });

    // Should be called only once with final value
    expect(mockOnChange).toHaveBeenCalledTimes(1);
    expect(mockOnChange).toHaveBeenCalledWith('John');
  });

  it('should show clear button when value is not empty', () => {
    render(<SearchInput value="John" onChange={mockOnChange} />);

    expect(screen.getByLabelText(/clear search/i)).toBeInTheDocument();
  });

  it('should not show clear button when value is empty', () => {
    render(<SearchInput value="" onChange={mockOnChange} />);

    expect(screen.queryByLabelText(/clear search/i)).not.toBeInTheDocument();
  });

  it('should clear input when clear button is clicked', async () => {
    vi.useRealTimers(); // Use real timers for this test
    const user = userEvent.setup();
    render(<SearchInput value="John" onChange={mockOnChange} />);

    const clearButton = screen.getByLabelText(/clear search/i);
    await user.click(clearButton);

    // Should call onChange immediately (no debounce for clear)
    expect(mockOnChange).toHaveBeenCalledWith('');
    vi.useFakeTimers(); // Restore fake timers
  });

  it('should show search icon', () => {
    render(<SearchInput value="" onChange={mockOnChange} />);

    // Check for SVG element
    const searchIcon = screen.getByLabelText(/search/i);
    expect(searchIcon).toBeInTheDocument();
  });

  it('should be keyboard accessible', async () => {
    vi.useRealTimers(); // Use real timers for this test
    const user = userEvent.setup();
    render(<SearchInput value="" onChange={mockOnChange} />);

    const input = screen.getByRole('textbox');

    // Can focus
    input.focus();
    expect(input).toHaveFocus();

    // Can type
    await user.keyboard('Test');
    expect(input).toHaveValue('Test');
    vi.useFakeTimers(); // Restore fake timers
  });

  it('should accept custom className', () => {
    const { container } = render(
      <SearchInput value="" onChange={mockOnChange} className="custom-class" />
    );

    const wrapper = container.firstChild;
    expect(wrapper).toHaveClass('custom-class');
  });

  it('should handle rapid typing correctly', async () => {
    const user = userEvent.setup({ delay: null });
    render(<SearchInput value="" onChange={mockOnChange} debounceMs={300} />);

    const input = screen.getByRole('textbox');

    // Simulate rapid typing
    await user.type(input, 'J');
    act(() => {
      vi.advanceTimersByTime(50);
    });
    await user.type(input, 'o');
    act(() => {
      vi.advanceTimersByTime(50);
    });
    await user.type(input, 'h');
    act(() => {
      vi.advanceTimersByTime(50);
    });
    await user.type(input, 'n');

    // Fast-forward past debounce
    act(() => {
      vi.runAllTimers();
    });

    // Should only call once with final value
    expect(mockOnChange).toHaveBeenCalledTimes(1);
    expect(mockOnChange).toHaveBeenCalledWith('John');
  });

  it('should update when controlled value changes externally', () => {
    const { rerender } = render(<SearchInput value="John" onChange={mockOnChange} />);

    const input = screen.getByRole('textbox') as HTMLInputElement;
    expect(input.value).toBe('John');

    // Update prop
    rerender(<SearchInput value="Jane" onChange={mockOnChange} />);

    expect(input.value).toBe('Jane');
  });
});
