import { render, screen } from '@testing-library/react';
import Nav from '../ui/nav';
import { useAuth0 } from '@auth0/auth0-react';
import { mocked } from 'jest-mock';

jest.mock('@auth0/auth0-react');
const mockedUseAuth0 = mocked(useAuth0, true);

describe('Nav login/signout button tests', () => {
  test('checks login button is visible when signed out', () => {
    mockedUseAuth0.mockReturnValue({
      isAuthenticated: false, // Mock signed out
      getAccessTokenSilently: jest.fn(),
      getAccessTokenWithPopup: jest.fn(),
      getIdTokenClaims: jest.fn(),
      loginWithRedirect: jest.fn(),
      loginWithPopup: jest.fn(),
      logout: jest.fn(),
      handleRedirectCallback: jest.fn(),
      isLoading: false,
    });
    render(<Nav />);
    const logIn = screen.getByText(/Log In/i);
    expect(logIn).toBeInTheDocument();
  });

  test('checks signout button is visible when logged in', () => {
    mockedUseAuth0.mockReturnValue({
      isAuthenticated: true, // Mock logged in
      getAccessTokenSilently: jest.fn(),
      getAccessTokenWithPopup: jest.fn(),
      getIdTokenClaims: jest.fn(),
      loginWithRedirect: jest.fn(),
      loginWithPopup: jest.fn(),
      logout: jest.fn(),
      handleRedirectCallback: jest.fn(),
      isLoading: false,
    });
    render(<Nav />);
    const signOut = screen.getByText(/Sign Out/i);
    expect(signOut).toBeInTheDocument();
  });
});
