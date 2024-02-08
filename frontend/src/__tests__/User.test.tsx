import { render, screen } from '@testing-library/react';
import Dashboard from '../routes/Dashboard';

test('renders coming soon text', () => {
  render(<Dashboard />);
  const dashboardText = screen.getByText(/This will be the user dashboard/i);
  expect(dashboardText).toBeInTheDocument();
});
