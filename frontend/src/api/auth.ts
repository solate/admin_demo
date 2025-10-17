import http from './http'

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: {
    id: number
    username: string
  }
}

export const authApi = {
  login: (data: LoginRequest): Promise<LoginResponse> => {
    return http.post('/auth/login', data)
  },
  
  register: (data: LoginRequest): Promise<LoginResponse> => {
    return http.post('/auth/register', data)
  }
}
