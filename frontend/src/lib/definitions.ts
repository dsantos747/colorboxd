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
  TokenRefresh: string;
  UserId: string;
  Username: string;
  UserGivenName: string;
}

export interface List {
  id: string;
  name: string;
  version: number;
  filmCount: number;
  description: string;
  entries: ListEntry[];
}

export interface ListEntry {
  rank: number; // What is this number if the list isn't ranked?
  entryId: string;
  posterPickerUrl: string; // This might not be accessible with our auth level
  film: Film;
}

export interface Film {
  id: string;
  poster: ImageURLs;
  adultPoster: ImageURLs;
  adult: boolean;
}

export interface ImageURLs {
  width: number;
  height: number;
  url: string;
}

export interface EntryWithImage {
  entryId: string;
  filmId: string;
  name: string;
  releaseYear: number;
  adult: boolean;
  posterCustomisable: boolean;
  posterUrl: string;
  adultPosterUrl: string;
  ImageInfo: ImageInfo;
}

interface ImageInfo {
  Path: string;
  Hex: string;
  Color: {
    R: number;
    G: number;
    B: number;
  };
  Hue: number;
}
