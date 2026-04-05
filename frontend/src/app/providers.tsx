import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { createContext, useContext, useMemo, useState, type ReactNode } from "react";

import type { User } from "@/entities/user/model";
import {
  clearSession,
  loadSelectedHouseId,
  loadUser,
  saveSelectedHouseId,
  saveTokens,
  saveUser,
  type AuthTokens,
} from "@/shared/lib/storage";

interface SessionContextValue {
  user: User | null;
  selectedHouseId: number | null;
  setAuth: (user: User, tokens: AuthTokens) => void;
  logout: () => void;
  setSelectedHouseId: (houseId: number | null) => void;
}

const SessionContext = createContext<SessionContextValue | null>(null);

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

export function AppProviders({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(() => loadUser());
  const [selectedHouseId, setSelectedHouseIdState] = useState<number | null>(() => loadSelectedHouseId());

  const value = useMemo<SessionContextValue>(
    () => ({
      user,
      selectedHouseId,
      setAuth(nextUser, tokens) {
        setUser(nextUser);
        saveUser(nextUser);
        saveTokens(tokens);
      },
      logout() {
        setUser(null);
        setSelectedHouseIdState(null);
        clearSession();
        queryClient.clear();
      },
      setSelectedHouseId(houseId) {
        setSelectedHouseIdState(houseId);
        if (houseId === null) {
          return;
        }
        saveSelectedHouseId(houseId);
      },
    }),
    [selectedHouseId, user],
  );

  return (
    <QueryClientProvider client={queryClient}>
      <SessionContext.Provider value={value}>{children}</SessionContext.Provider>
    </QueryClientProvider>
  );
}

export function useSession() {
  const context = useContext(SessionContext);
  if (!context) {
    throw new Error("useSession must be used within AppProviders");
  }
  return context;
}
