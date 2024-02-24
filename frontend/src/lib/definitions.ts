export interface GetTokenResponse {
  token: AccessTokenResponse;
  user: User;
  lists: ListSummary[];
}

export interface AccessTokenResponse {
  access_token: string;
  token_type: string;
  refresh_token: string;
  expires_in: number;
  notBefore: number;
  issuer: string;
  encodedToken: string;
}

export interface User {
  id: string;
  displayName: string;
  givenName: string;
  username: string;
}

export interface ListSummary {
  id: string;
  name: string;
  version: number;
  filmCount: number;
  description: string;
}

export interface UserToken {
  Token: string;
  TokenType: string;
  TokenExpiresIn: number;
  UserId: string;
  Username: string;
  UserGivenName: string;
}
