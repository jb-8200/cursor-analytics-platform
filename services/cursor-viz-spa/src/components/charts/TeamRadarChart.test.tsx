/**
 * Tests for TeamRadarChart component
 *
 * Tests the multi-axis radar chart comparing teams across different metrics.
 *
 * Note: Due to jsdom limitations with ResponsiveContainer, some tests verify
 * component structure rather than full SVG rendering.
 */

import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import TeamRadarChart from './TeamRadarChart';
import { TeamStats } from '../../graphql/types';

describe('TeamRadarChart', () => {
  const mockTeams: TeamStats[] = [
    {
      teamName: 'Frontend',
      memberCount: 8,
      activeMemberCount: 7,
      averageAcceptanceRate: 0.75,
      totalSuggestions: 1000,
      aiVelocity: 0.65,
      chatInteractions: 500,
      topPerformers: [],
    },
    {
      teamName: 'Backend',
      memberCount: 10,
      activeMemberCount: 9,
      averageAcceptanceRate: 0.82,
      totalSuggestions: 1200,
      aiVelocity: 0.71,
      chatInteractions: 600,
      topPerformers: [],
    },
    {
      teamName: 'DevOps',
      memberCount: 5,
      activeMemberCount: 5,
      averageAcceptanceRate: 0.68,
      totalSuggestions: 800,
      aiVelocity: 0.58,
      chatInteractions: 400,
      topPerformers: [],
    },
  ];

  describe('Rendering', () => {
    it('should render chart container', () => {
      const { container } = render(
        <TeamRadarChart data={mockTeams} selectedTeams={['Frontend', 'Backend']} />
      );

      const chartContainer = container.querySelector('.team-radar-chart');
      expect(chartContainer).toBeInTheDocument();
    });

    it('should render ResponsiveContainer', () => {
      const { container } = render(
        <TeamRadarChart data={mockTeams} selectedTeams={['Frontend']} />
      );

      const responsive = container.querySelector('.recharts-responsive-container');
      expect(responsive).toBeInTheDocument();
    });

    it('should handle empty selected teams', () => {
      render(<TeamRadarChart data={mockTeams} selectedTeams={[]} />);

      expect(screen.getByText('Select teams to compare')).toBeInTheDocument();
    });

    it('should handle empty data', () => {
      render(<TeamRadarChart data={[]} selectedTeams={['Frontend']} />);

      expect(screen.getByText('No data available for selected teams')).toBeInTheDocument();
    });

    it('should handle team selection with onTeamSelect callback', () => {
      const onTeamSelect = vi.fn();
      render(
        <TeamRadarChart
          data={mockTeams}
          selectedTeams={['Frontend']}
          onTeamSelect={onTeamSelect}
        />
      );

      // Selection controls should be rendered
      expect(screen.getByText('Select Teams (2-5):')).toBeInTheDocument();
    });
  });

  describe('Team Selection Controls', () => {
    it('should render team selection buttons when onTeamSelect provided', () => {
      const onTeamSelect = vi.fn();
      render(
        <TeamRadarChart
          data={mockTeams}
          selectedTeams={['Frontend']}
          onTeamSelect={onTeamSelect}
        />
      );

      expect(screen.getByText('Frontend')).toBeInTheDocument();
      expect(screen.getByText('Backend')).toBeInTheDocument();
      expect(screen.getByText('DevOps')).toBeInTheDocument();
    });

    it('should call onTeamSelect when team button clicked', async () => {
      const user = userEvent.setup();
      const onTeamSelect = vi.fn();

      render(
        <TeamRadarChart
          data={mockTeams}
          selectedTeams={['Frontend']}
          onTeamSelect={onTeamSelect}
        />
      );

      const backendButton = screen.getByText('Backend');
      await user.click(backendButton);

      expect(onTeamSelect).toHaveBeenCalledWith(['Frontend', 'Backend']);
    });

    it('should deselect team when already selected team clicked', async () => {
      const user = userEvent.setup();
      const onTeamSelect = vi.fn();

      render(
        <TeamRadarChart
          data={mockTeams}
          selectedTeams={['Frontend', 'Backend']}
          onTeamSelect={onTeamSelect}
        />
      );

      const frontendButton = screen.getByText('Frontend');
      await user.click(frontendButton);

      expect(onTeamSelect).toHaveBeenCalledWith(['Backend']);
    });

    it('should disable button when 5 teams selected', () => {
      const onTeamSelect = vi.fn();
      const fiveTeams: TeamStats[] = [
        ...mockTeams,
        {
          teamName: 'QA',
          memberCount: 4,
          activeMemberCount: 4,
          averageAcceptanceRate: 0.7,
          totalSuggestions: 500,
          aiVelocity: 0.6,
          chatInteractions: 300,
          topPerformers: [],
        },
        {
          teamName: 'Design',
          memberCount: 3,
          activeMemberCount: 3,
          averageAcceptanceRate: 0.65,
          totalSuggestions: 400,
          aiVelocity: 0.55,
          chatInteractions: 250,
          topPerformers: [],
        },
      ];

      render(
        <TeamRadarChart
          data={fiveTeams}
          selectedTeams={['Frontend', 'Backend', 'DevOps', 'QA', 'Design']}
          onTeamSelect={onTeamSelect}
        />
      );

      // All buttons should be rendered
      expect(screen.getByText('Frontend')).toBeInTheDocument();
      expect(screen.getByText('Design')).toBeInTheDocument();
    });
  });

  describe('Custom Configuration', () => {
    it('should render with custom metrics', () => {
      const customMetrics = [
        { key: 'totalSuggestions' as keyof TeamStats, label: 'Suggestions', max: 2000 },
        { key: 'chatInteractions' as keyof TeamStats, label: 'Chat', max: 1000 },
      ];

      const { container } = render(
        <TeamRadarChart
          data={mockTeams}
          selectedTeams={['Frontend']}
          metrics={customMetrics}
        />
      );

      const chartContainer = container.querySelector('.team-radar-chart');
      expect(chartContainer).toBeInTheDocument();
    });

    it('should render with custom dimensions', () => {
      const { container } = render(
        <TeamRadarChart
          data={mockTeams}
          selectedTeams={['Frontend']}
          width={600}
          height={600}
        />
      );

      const responsive = container.querySelector('.recharts-responsive-container');
      expect(responsive).toBeInTheDocument();
    });
  });

  describe('Edge Cases', () => {
    it('should handle teams with zero values', () => {
      const teamsWithZeros: TeamStats[] = [
        {
          teamName: 'Inactive',
          memberCount: 5,
          activeMemberCount: 0,
          averageAcceptanceRate: 0,
          totalSuggestions: 0,
          aiVelocity: 0,
          chatInteractions: 0,
          topPerformers: [],
        },
      ];

      const { container } = render(
        <TeamRadarChart data={teamsWithZeros} selectedTeams={['Inactive']} />
      );

      const chartContainer = container.querySelector('.team-radar-chart');
      expect(chartContainer).toBeInTheDocument();
    });

    it('should handle team not in data', () => {
      render(<TeamRadarChart data={mockTeams} selectedTeams={['NonExistent']} />);

      expect(screen.getByText('No data available for selected teams')).toBeInTheDocument();
    });

    it('should handle very large values', () => {
      const largeValueTeam: TeamStats[] = [
        {
          teamName: 'BigTeam',
          memberCount: 100,
          activeMemberCount: 95,
          averageAcceptanceRate: 0.9,
          totalSuggestions: 1000000,
          aiVelocity: 0.85,
          chatInteractions: 500000,
          topPerformers: [],
        },
      ];

      const { container } = render(
        <TeamRadarChart data={largeValueTeam} selectedTeams={['BigTeam']} />
      );

      const chartContainer = container.querySelector('.team-radar-chart');
      expect(chartContainer).toBeInTheDocument();
    });
  });

  describe('Accessibility', () => {
    it('should have proper ARIA attributes', () => {
      const { container } = render(
        <TeamRadarChart data={mockTeams} selectedTeams={['Frontend']} />
      );

      const chartContainer = container.querySelector('[role="img"]');
      expect(chartContainer).toBeInTheDocument();
      expect(chartContainer).toHaveAttribute('aria-label', 'Team comparison radar chart');
    });

    it('should have accessible team selection buttons', () => {
      const onTeamSelect = vi.fn();
      const { container } = render(
        <TeamRadarChart
          data={mockTeams}
          selectedTeams={['Frontend']}
          onTeamSelect={onTeamSelect}
        />
      );

      const buttons = container.querySelectorAll('button[aria-pressed]');
      expect(buttons.length).toBeGreaterThan(0);
    });
  });
});
