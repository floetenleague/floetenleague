import { create } from "zustand";
import {
  Configuration,
  DefaultApi,
  FLBingoFieldStatus,
  FLOverview,
  LoginState,
  SetFieldStatusRequest,
} from "./generated";

export const API = new DefaultApi(
  new Configuration({ basePath: window.location.origin })
);

export interface GlobalState {
  auth: LoginState;
  overview: FLOverview;
  refreshAuth: () => Promise<void>;
  refreshOverview: () => Promise<void>;
  init: () => Promise<void>;
  poeLogin: () => Promise<string>;
  logout: () => Promise<void>;
  joinBoard: (bingoId: number) => Promise<void>;
  login: (username: string, password: string) => Promise<void>;
  setFieldStatus: (req: SetFieldStatusRequest) => Promise<void>;
}

export const useStore = create<GlobalState>()((set, get) => ({
  auth: {
    id: -1,
    username: "guest",
    permission: "Unverified",
    loggedIn: false,
  },
  init: async () => {
    await Promise.all([get().refreshAuth(), get().refreshOverview()]);
  },
  overview: { bingos: [] },
  setFieldStatus: async (req: SetFieldStatusRequest) => {
    await API.setFieldStatus(req);
    await get().refreshOverview();
  },
  logout: async () => {
    await API.logout();
    await get().refreshAuth();
  },
  joinBoard: async (bingoID) => {
    await API.joinBingo({ bingoID });
    await get().refreshOverview();
  },
  refreshAuth: async () => {
    const auth = await API.getAuth();
    set({ auth });
  },
  refreshOverview: async () => {
    const overview = await API.getOverview();
    set({ overview });
  },
  poeLogin: async () => {
    return JSON.parse(await API.loginPoe());
  },
  login: async (username, password) => {
    await API.loginInternal({ loginRequest: { username, password } });
    await get().refreshAuth();
  },
}));
