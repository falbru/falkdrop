export interface AuthProvider {
  init(): Promise<boolean>;
  login(options?: { redirectUri?: string }): Promise<void>;
  logout(options?: { redirectUri?: string }): Promise<void>;
  isAuthenticated(): boolean;
  getToken(): string | undefined;
}
