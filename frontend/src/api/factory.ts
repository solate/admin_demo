import http from './http'

export interface Factory {
  id: number
  name: string
  address: string
  owner: string
  createdAt: string
}

export interface FactoryListParams {
  page?: number
  pageSize?: number
  keyword?: string
}

export interface FactoryListResponse {
  list: Factory[]
  total: number
  page: number
  pageSize: number
}

export const factoryApi = {
  getList: (params: FactoryListParams): Promise<FactoryListResponse> => {
    return http.get('/factories', { params })
  },
  
  create: (data: Omit<Factory, 'id' | 'createdAt'>): Promise<Factory> => {
    return http.post('/factories', data)
  },
  
  update: (id: number, data: Omit<Factory, 'id' | 'createdAt'>): Promise<Factory> => {
    return http.put(`/factories/${id}`, data)
  },
  
  delete: (id: number): Promise<void> => {
    return http.delete(`/factories/${id}`)
  },
  
  batchDelete: (ids: number[]): Promise<void> => {
    return http.delete('/factories/batch', { data: { ids } })
  }
}
