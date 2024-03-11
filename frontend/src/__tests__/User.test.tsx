import { render, screen } from '@testing-library/react';
import { UserTokenContext, ListSummaryContext, ListContext } from '../lib/contexts';
import UserAuth from '../routes/UserAuth';
import { List, ListSummary, UserToken } from '../lib/definitions';
import { MemoryRouter } from 'react-router-dom';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

const testToken: UserToken = {
  Token: 'token',
  TokenExpiresIn: 1000,
  TokenRefresh: 'refresh',
  TokenType: 'type',
  UserGivenName: 'John',
  UserId: '123',
  Username: 'JohnDoe123',
};

const testListSummary: ListSummary[] = [
  { id: '1', description: 'List 1 with 5 films', filmCount: 5, name: 'List 1', version: 1 },
  { id: '2', description: 'List 2 with 10 films', filmCount: 10, name: 'List 2', version: 1 },
  { id: '3', description: 'List 3 with 50 films', filmCount: 50, name: 'List 3', version: 1 },
];

const testList: List = {
  id: '2',
  description: 'List 2 with 10 films',
  filmCount: 10,
  name: 'List 2',
  version: 1,
  entries: [
    {
      adult: false,
      adultPosterUrl: '',
      entryId: '0',
      filmId: 'abc',
      ImageInfo: { Color: { R: 255, G: 0, B: 0 }, Hex: '#FF0000', Hue: 0, Path: '123' },
      name: 'Test Film',
      posterCustomisable: false,
      posterUrl: '123',
      releaseYear: 1997,
    },
  ],
};

const mockUserPage = (mockToken: UserToken | null, mockList: List | null) => {
  return (
    <UserTokenContext.Provider value={{ userToken: mockToken, setUserToken: jest.fn() }}>
      <ListSummaryContext.Provider value={{ listSummary: testListSummary, setListSummary: jest.fn() }}>
        <ListContext.Provider value={{ list: mockList, setList: jest.fn() }}>
          <MemoryRouter>
            <UserAuth />
          </MemoryRouter>
        </ListContext.Provider>
      </ListSummaryContext.Provider>
    </UserTokenContext.Provider>
  );
};

test('renders john doe user page, with 3 lists', () => {
  render(mockUserPage(testToken, null));
  const listLabel1 = screen.getByText(testListSummary[0].name);
  const listLabel2 = screen.getByText(testListSummary[1].name);
  const listLabel3 = screen.getByText(testListSummary[2].name);
  expect(listLabel1).toBeInTheDocument();
  expect(listLabel2).toBeInTheDocument();
  expect(listLabel3).toBeInTheDocument();
});

test('sets active list, expect ui change', () => {
  render(mockUserPage(testToken, testList));
  const hintText = screen.getByText(/Hint: Click an item to make it the start of the list./i);
  expect(hintText).toBeInTheDocument();
});
