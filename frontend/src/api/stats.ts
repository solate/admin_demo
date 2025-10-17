import http from './http'

export interface StatsOverview {
  totalValue: number
  productCount: number
  stockInCount: number
  stockOutCount: number
  stockCount: number
}

export interface TrendData {
  date: string
  inCount: number
  outCount: number
}

export interface StatsTrends {
  data: TrendData[]
  range: 'daily' | 'weekly' | 'monthly'
}

export const statsApi = {
  getOverview: (): Promise<StatsOverview> => {
    return http.get('/stats/overview')
  },
  
  getTrends: (range: 'daily' | 'weekly' | 'monthly'): Promise<StatsTrends> => {
    return http.get('/stats/trends', { params: { range } })
  }
}
