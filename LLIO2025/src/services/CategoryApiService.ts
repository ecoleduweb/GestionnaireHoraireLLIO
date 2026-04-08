// src/services/CategoryApiService.ts
import type { Category } from "../Models/index";
import { DELETE, GET, POST, PUT } from "../ts/server";

interface CategoriesResponse {
  categories: Category[];
}

interface CategoryResponse {
  category: Category;
}

// Correspondance avec la structure de réponse réelle du backend
interface CategoryCreateResponse {
  reponse: string;
  category?: Category;
}

// Récupérer toutes les catégories
const getAllCategories = async (): Promise<Category[]> => {
  try {
    const response = await GET<CategoriesResponse>("/categories");
    return response.categories;
  } catch (error) {
    console.error("Erreur lors de la récupération des catégories:", error);
    throw error;
  }
};

// Récupérer les catégories d'un projet spécifique
const getCategoriesByProject = async (projectId: number): Promise<Category[]> => {
  try {
    const response = await GET<CategoriesResponse>(`/project/${projectId}/categories`);
    return response?.categories || []; // Retourne un tableau vide si null ou undefined
  } catch (error) {
    console.error(`Erreur lors de la récupération des catégories du projet (ID: ${projectId}):`, error);
    return [];
  }
};

// Récupérer une catégorie par son ID
const getCategoryById = async (id: number): Promise<Category> => {
  try {
    const response = await GET<CategoryResponse>(`/category/${id}`);
    return response.category;
  } catch (error) {
    console.error(`Erreur lors de la récupération de la catégorie (ID: ${id}):`, error);
    throw error;
  }
};

// Créer une nouvelle catégorie
const createCategory = async (
  categoryData: Partial<Category>, 
  projectId?: number
): Promise<Category> => {
  try {
    // Préparer le payload avec les champs obligatoires
    const categoryPayload = {
      name: categoryData.name || "",
      description: "Créé depuis l'application",
      billable: categoryData.billable || false,
      projectId: projectId
    };
    
    // Envoyer la requête
    const response = await POST<any, CategoryCreateResponse>("/category", categoryPayload);
    
    // Traiter la réponse selon différents cas possibles
    if (response?.category) return response.category;
    // if (response?.activity) return response.activity;
    
    throw new Error('Format de réponse API inattendu');
  } catch (error) {
    // Propager l'erreur pour traitement par l'appelant
    throw error;
  }
};

// Récupérer les catégories d'un projet spécifique
const deleteCategory = async (categoryId: number): Promise<void> => {
  try {
    await DELETE(`/category/${categoryId}`, false, false);
  } catch (error) {
    if (error)
      throw new Error(`Erreur lors de la suppression de la catégorie (ID: ${categoryId}):`, error)
  }
};

const changeCategoryName = async (newName: string, category: Category): Promise<boolean> => {
  try {
    let newCategory = {
      id: category.id,
      name: newName,
      description: category.description,
      billable: category.billable,
      timeSpent: category.timeSpent,
      timeEstimated: category.timeEstimated
    }

    await PUT<Category, Category>(`/category`, newCategory);

    return true;

  } catch (error) {
    if (error)
      throw new Error(`Erreur lors du renommage de la catégorie (ID: ${category.id}):`, error)
  }

  return false;
}

export const CategoryApiService = {
  getAllCategories,
  getCategoryById,
  deleteCategory,
  createCategory,
  getCategoriesByProject,
  changeCategoryName
};