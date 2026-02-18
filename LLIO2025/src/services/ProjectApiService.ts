// ProjectApiService.ts
import { ProjectStatus } from "$lib/types/enums";
import type { ProjectBase, Project, DetailedProject } from "../Models/index";
import { DELETE, GET, POST, PUT } from "../ts/server";

interface ProjectDeleteResponse {
  deleted: boolean;
}

interface ProjectResponse {
  project: Project;
}

const createProject = async (project: ProjectBase): Promise<Project> => {
  if (project.status === undefined) {
    project.status = ProjectStatus.InProgress;
  }
  try {
    const response = await POST<ProjectBase, ProjectResponse>("/project", project);
    return response.project;
  } catch (error) {
    console.error("Erreur lors de la création du projet:", error);
    throw new Error("Échec de la création du projet. Veuillez réessayer.");
  }
};

const updateProject = async (project: ProjectBase): Promise<Project> => {
  if (project.status === undefined) {
    project.status = ProjectStatus.InProgress;
  }
  try {
    const response = await PUT<ProjectBase, ProjectResponse>("/project", project);
    return response.project;
  } catch (error) {
    console.error("Erreur lors de la création du projet:", error);
    throw new Error("Échec de la mise à jour du projet. Veuillez réessayer.");
  }
};

const deleteProject = async (projectId: number): Promise<void> => {
  try {
    const response = await DELETE(`/project/${projectId}`);
    return response;
  } catch (error) {
    console.error("Erreur lors de la suppression du projet:", error);
    throw new Error("Erreur lors de la suppression du projet. Veuillez réessayer.");
  }
};

const getProjects = async(): Promise<Project[]> => {
  try {
    const response = await GET<{projects: Project[]}>("/projects");
    return response.projects;
  } catch (error) {
    console.error("Erreur lors de la récupération des projets:", error);
    throw new Error("Échec de la récupération des projets. Veuillez réessayer.");
  }
}

const getProject = async(id: number): Promise<Project> => {
  try {
    const response = await GET<{project: Project}>(`/project/${id}`);
    return response.project;
  } catch (error) {
    console.error("Erreur lors de la récupération du projet:", error);
    throw new Error("Échec de la récupération du projet. Veuillez réessayer.");
  }
}

const getDetailedProjects = async(): Promise<DetailedProject[]> => {
  try {
    const response = await GET<{projects:   DetailedProject[]}>("/projects/detailed");
    return response.projects;
  } catch (error) {
    console.error("Erreur lors de la récupération des projets:", error);
    throw new Error("Échec de la récupération des projets. Veuillez réessayer.");
  }
}

const getCurrentUserProjects = async(): Promise<DetailedProject[]> => {
  try {
    const response = await GET<{projects: DetailedProject[]}>("/projects/me/detailed");
    return response.projects;
  } catch (error) {
    console.error("Erreur lors de la récupération des projets:", error);
    throw new Error("Échec de la récupération des projets. Veuillez réessayer.");
  }
}

const addCoManagerToProject = async(projectId: number, userId: number): Promise<void> => {
  try {
    await POST(`/project/${projectId}/coManager/${userId}`, {});
    return;
  } catch (error) {
    console.error("Erreur lors de l'ajout du co-chargé au projet:", error);
    throw new Error("Erreur à l'ajout du co-chargé : " + error.message)
  }
}

export const ProjectApiService = {
  createProject,
  updateProject,
  deleteProject,
  getProjects,
  getProject,
  getDetailedProjects,
  getCurrentUserProjects,
  addCoManagerToProject
};
