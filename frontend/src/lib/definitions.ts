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
  hex1: string;
  hex2: string;
}

interface ImageInfo {
  Path: string;
  Hex: string;
  RGB: {
    R: number;
    G: number;
    B: number;
  };
  HSL: {
    H: number;
    S: number;
    L: number;
  };
}

export const sorts = [
  { id: 'hue', name: 'Hue' },
  { id: 'lum', name: 'Luminosity' },
  { id: 'brightDomHue', name: 'Bright Dominant Hue' },
  { id: 'inverseStep_8', name: 'Inverse Step (8)' },
  { id: 'inverseStep_12', name: 'Inverse Step (12)' },
  { id: 'inverseStep2_8', name: 'Inverse Step v2 (8)' },
  { id: 'inverseStep2_12', name: 'Inverse Step v2 (12)' },
  { id: 'BRBW1', name: 'BRBW1' },
  { id: 'BRBW2', name: 'BRBW2' },
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
