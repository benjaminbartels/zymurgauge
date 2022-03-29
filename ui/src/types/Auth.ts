export interface AuthRequest {
  username: string;
  password: string;
}

export interface AuthResponse {
  token: string;
}
