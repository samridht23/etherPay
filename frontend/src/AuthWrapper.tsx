import React, {
  useState,
  useEffect,
  createContext,
  useMemo,
  ReactNode,
} from 'react';
import { useNavigate, useLocation } from 'react-router-dom';

interface User {
  address: string;
  about: string;
  username: string;
  profileImageUrl: string;
  bannerImageUrl: string;
}

interface AuthContextValue {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  authStatus: () => void;
}

const defaultAuthContextValue: AuthContextValue = {
  user: null,
  isAuthenticated: false,
  isLoading: true,
  authStatus: () => { },
};

export const AuthContext = createContext<AuthContextValue>(
  defaultAuthContextValue
);

const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const navigate = useNavigate();
  const location = useLocation();

  const [authState, setAuthState] = useState<AuthContextValue>(
    defaultAuthContextValue
  );

  const authStatus = async () => {
    try {
      const response = await fetch(
        `${import.meta.env.VITE_API_URL}/auth-status`,
        {
          method: 'GET',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );
      if (response.ok) {
        const data = await response.json();
        setAuthState({
          user: {
            address: data.data.name,
            about: data.data.abobut,
            username: data.data.username,
            profileImageUrl: data.data.profile_image_url,
            bannerImageUrl: data.data.banner_image_url,
          },
          isAuthenticated: true,
          isLoading: false,
          authStatus
        });
      } else {
        setAuthState({
          user: null,
          isAuthenticated: false,
          isLoading: false,
          authStatus
        });
      }
    } catch (err) {
      setAuthState({
        user: null,
        isAuthenticated: false,
        isLoading: false,
        authStatus
      });
    }
  };

  useEffect(() => {
    authStatus();
  }, [navigate, location]);

  const authValue = useMemo(
    () => ({
      user: authState.user,
      isAuthenticated: authState.isAuthenticated,
      isLoading: authState.isLoading,
      authStatus, // Expose authStatus function
    }),
    [authState]
  );

  return (
    <AuthContext.Provider value={authValue}>{children}</AuthContext.Provider>
  );
};

export default AuthProvider;
