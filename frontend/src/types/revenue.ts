export interface RevenueData {
    date?: string;
    count: number;
  }
  
  export const fetchRevenue = (): RevenueData[] => {
    return [{ date: '2025-04-01', count: 10 }]
  }