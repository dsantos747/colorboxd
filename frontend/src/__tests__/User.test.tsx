import { render, screen } from '@testing-library/react';
import User from '../routes/UserContent';

// Need to somehow mock the access token - how tf...
test('renders user dashboard page', () => {
  render(<User accessToken='' />);
  const dashboardText = screen.getByText(/This will be the user dashboard/i);
  expect(dashboardText).toBeInTheDocument();
});
