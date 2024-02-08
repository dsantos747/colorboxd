import { render, screen } from '@testing-library/react';
import Home from '../routes/Home';

test('renders coming soon text', () => {
  render(<Home />);
  const comingSoon = screen.getByText(/ColorBoxd is under construction./i);
  expect(comingSoon).toBeInTheDocument();
});
