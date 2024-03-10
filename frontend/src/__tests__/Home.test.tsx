import { render, screen } from '@testing-library/react';
import Home from '../routes/Home';

test('renders slogan text', () => {
  render(<Home />);
  const slogan = screen.getByAltText(/A list sorted with Colorboxd/i);
  expect(slogan).toBeInTheDocument();
});
