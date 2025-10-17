import http from './http'

export interface Product {
  id: number
  name: string
  sku: string
  category: string
  price: number
  unit: string
  stock: number
  remark: string
  createdAt: string
}

export interface ProductListParams {
  page?: number
  pageSize?: number
  keyword?: string
}

export interface ProductListResponse {
  list: Product[]
  total: number
  page: number
  pageSize: number
}

export interface StockOperation {
  productId: number
  quantity: number
  remark?: string
}

export const productApi = {
  getList: (params: ProductListParams): Promise<ProductListResponse> => {
    return http.get('/products', { params })
  },
  
  create: (data: Omit<Product, 'id' | 'createdAt' | 'stock'>): Promise<Product> => {
    return http.post('/products', data)
  },
  
  update: (id: number, data: Omit<Product, 'id' | 'createdAt' | 'stock'>): Promise<Product> => {
    return http.put(`/products/${id}`, data)
  },
  
  delete: (id: number): Promise<void> => {
    return http.delete(`/products/${id}`)
  },
  
  stockIn: (data: StockOperation): Promise<void> => {
    return http.post('/inventory/in', data)
  },
  
  stockOut: (data: StockOperation): Promise<void> => {
    return http.post('/inventory/out', data)
  },
  
  getStock: (productId: number): Promise<{ quantity: number }> => {
    return http.get(`/inventory/${productId}`)
  }
}
