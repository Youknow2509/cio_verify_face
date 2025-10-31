// src/components/charts/AttendanceChart.tsx

import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import styles from './AttendanceChart.module.scss';

interface ChartData {
  date: string;
  checkIns: number;
  checkOuts: number;
  lateArrivals: number;
}

interface AttendanceChartProps {
  data: ChartData[];
  height?: number;
}

export function AttendanceChart({ data, height = 300 }: AttendanceChartProps) {
  return (
    <div className={styles.chartContainer} style={{ height }}>
      <ResponsiveContainer width="100%" height="100%">
        <LineChart data={data} margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
          <CartesianGrid strokeDasharray="3 3" stroke="var(--color-border)" />
          <XAxis 
            dataKey="date" 
            tickFormatter={(value) => new Date(value).toLocaleDateString('vi-VN', { month: 'short', day: 'numeric' })}
            tick={{ fontSize: 12, fill: 'var(--color-text-secondary)' }}
          />
          <YAxis 
            tick={{ fontSize: 12, fill: 'var(--color-text-secondary)' }}
          />
          <Tooltip 
            labelFormatter={(value) => new Date(value).toLocaleDateString('vi-VN')}
            formatter={(value: number, name: string) => [
              value, 
              name === 'checkIns' ? 'Check-ins' : 
              name === 'checkOuts' ? 'Check-outs' : 
              'Đi trễ'
            ]}
            contentStyle={{
              backgroundColor: 'var(--color-surface)',
              border: '1px solid var(--color-border)',
              borderRadius: 'var(--radius-sm)',
              fontSize: '14px'
            }}
          />
          <Legend 
            wrapperStyle={{ fontSize: '14px' }}
          />
          <Line 
            type="monotone" 
            dataKey="checkIns" 
            stroke="var(--color-primary)" 
            strokeWidth={2}
            name="Check-ins"
            dot={{ fill: 'var(--color-primary)', strokeWidth: 2, r: 4 }}
            activeDot={{ r: 6 }}
          />
          <Line 
            type="monotone" 
            dataKey="checkOuts" 
            stroke="var(--color-success)" 
            strokeWidth={2}
            name="Check-outs"
            dot={{ fill: 'var(--color-success)', strokeWidth: 2, r: 4 }}
            activeDot={{ r: 6 }}
          />
          <Line 
            type="monotone" 
            dataKey="lateArrivals" 
            stroke="var(--color-warning)" 
            strokeWidth={2}
            name="Đi trễ"
            dot={{ fill: 'var(--color-warning)', strokeWidth: 2, r: 4 }}
            activeDot={{ r: 6 }}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}