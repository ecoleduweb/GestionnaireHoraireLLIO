import { ProjectStatus, UserRole } from "../lib/types/enums";

export interface Activity {
    id?: number;
    name: string;
    description: string;
    startDate: Date;
    endDate: Date;
    userId: number;
    projectId: number;
    projectName: string;
    categoryId: number;
}

export interface UserInfo {
    firstName: string;
    lastName: string;
    role : UserRole;
    email?: string;
}

export interface User {
    id: number;
    firstName: string;
    lastName: string;
    email: string;
    role: UserRole;
}

export interface Employee {
    name: string;
    categories: Category[];
}

export interface Category {
    id: number;
    name: string;
    description?: string;
    billable?: boolean;
    timeSpent: number;
    timeEstimated: number;
}

export interface ActivityUpdateResponse {
    updated_activity: Activity;
}

/*************** Project ***************/
export interface ProjectBase {
    uniqueId: string;
    name?: string;
    managerId: number;
    billable: boolean;
    status?: ProjectStatus;
    estimatedHours?: number;
}

export interface UpdateProject extends ProjectBase {
    id: number;
}

export interface Project extends UpdateProject {
    color: string;
    totalTimeEstimated: number;
    totalTimeRemaining: number;
    totalTimeSpent: number;
    isArchived?: boolean;
    status: ProjectStatus;
    createdAt: Date;
    updatedAt: Date;
}

export interface DetailedProject extends Project {
    id: number;
    color: string;
    lead: string;
    coLeads: CoLead[];
    employees: Employee[];
    totalTimeEstimated: number;
    totalTimeRemaining: number;
    totalTimeSpent: number;
    isArchived?: boolean;
    status: ProjectStatus;
    createdAt: Date;
    updatedAt: Date;
}

export interface CoLead {
  id: number;
  name: string;
}