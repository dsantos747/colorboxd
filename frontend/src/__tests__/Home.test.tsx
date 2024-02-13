import { render, screen } from '@testing-library/react';
import ComingSoon from '../routes/ComingSoon';

test('renders coming soon text', () => {
  render(<ComingSoon />);
  const comingSoon = screen.getByText(/ColorBoxd is under construction./i);
  expect(comingSoon).toBeInTheDocument();
});
