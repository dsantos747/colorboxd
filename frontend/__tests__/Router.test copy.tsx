import { render, screen } from '@testing-library/react';
import App from '../routes/Home';

test('renders coming soon text link', () => {
  render(<App />);
  const comingSoon = screen.getByText(/ColorBoxd is under construction./i);
  expect(comingSoon).toBeInTheDocument();
});
