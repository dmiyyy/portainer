import {
  createContext,
  useContext,
  useMemo,
  useReducer,
  PropsWithChildren,
} from 'react';

interface RowContextState {
  isLoading: boolean;
  toggleIsLoading(): void;
}

const RowContext = createContext<RowContextState | null>(null);

export function RowProvider({ children }: PropsWithChildren<unknown>) {
  const [isLoading, toggleIsLoading] = useReducer((state) => !state, false);

  const state = useMemo(
    () => ({ isLoading, toggleIsLoading }),
    [isLoading, toggleIsLoading]
  );

  return <RowContext.Provider value={state}>{children}</RowContext.Provider>;
}

export function useRowContext() {
  const context = useContext(RowContext);
  if (!context) {
    throw new Error('should be nested under RowProvider');
  }

  return context;
}
