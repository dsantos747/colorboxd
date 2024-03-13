export interface UserToken {
  Token: string;
  TokenType: string;
  TokenExpiresIn: number;
  TokenRefresh: string;
  UserId: string;
  Username: string;
  UserGivenName: string;
}

export interface ListSummary {
  id: string;
  name: string;
  version: number;
  filmCount: number;
  description: string;
}

export interface List extends ListSummary {
  entries: EntryWithImage[];
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
  sorts: SortRanks;
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

export const sorts = [
  { id: 'hue', name: 'Hue' },
  { id: 'val', name: 'Luminosity' },
  { id: 'brightHue', name: 'Bright Hue' },
  { id: 'brightDomHue', name: 'Dominant Bright Hue' },
  // { id: 'step', name: 'Alternating Step Sort' },
  // { id: 'hilbert', name: 'Hilbert Sort' },
  // { id: 'cie2000', name: 'CIELAB2000 Sort' },
] as const;

type SortTypes = (typeof sorts)[number];

type SortIds = SortTypes['id'];

type SortRanks = {
  [K in SortIds]: number;
};

export type SortModeType = {
  sortMode: SortTypes;
  visible: boolean;
  reverse: boolean;
};
