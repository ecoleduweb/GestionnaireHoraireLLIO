import type { User } from "lucide-svelte";

export interface UsersResponse {
  users: User[];
}

export interface TimeBankConfig {
  startDate: string;
  hoursPerWeek: number;
  offset: number;
}

export  interface TimeBankResponse {
  isConfigured: boolean;
  timeInBank: number;
}

 export interface HoursConfig {
    startDate: string;
    offset: number;
    hoursWorked: number;
  }

  
