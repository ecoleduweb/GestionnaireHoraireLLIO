import type { User, UserInfo } from '../Models/index';
import { DELETE, GET, PATCH, POST } from '../ts/server';

interface UsersResponse {
  users: User[];
}

const getUserInfo = async (): Promise<UserInfo> => {
  try {
    const response = await GET<UserInfo>('/user/me');
    return response;
  } catch (error) {
    console.error("Erreur lors de la récupération des informations de l'utilisateur:", error);
    throw error;
  }
};

const getAllUsers = async (): Promise<User[]> => {
  try {
    const response = await GET<UsersResponse>("/users");
    return response.users;
  } catch (error) {
    console.error("Erreur lors de la récupération des utilisateurs:", error);
    throw new Error("Échec de la récupération des utilisateurs. Veuillez réessayer.");
  }
};

const getAllManagersAdmin = async (): Promise<User[]> => {
  try {
    const response = await GET<User[]>("/users?role=1&role=2");
    return response;
  } catch (error) {
    console.error("Erreur lors de la récupération des managers:", error);
    throw new Error("Échec de la récupération des managers. Veuillez réessayer.");
  }
};
const getUsers = async (): Promise<UserInfo[]> => {
  try {
    const response = await GET<UserInfo[]>('/users');
    return response;
  }
  catch (error) {
    console.error("Erreur lors de la récupération des utilisateurs:", error);
    throw error;
  }
};

const updateUserRole = async (userId: number, role: number): Promise<void> => {
  try {
    const response = await PATCH<{ role: number }>(
      `/user/${userId}/role`,
      { role }
    );
    return response;
  } catch (error) {
    console.error("Erreur lors de la mise à jour du rôle de l'utilisateur:", error);
    throw error;
  }
}

const logOut = async (): Promise<void> => {
  try {
    await POST('/logout', {});
  } catch (error) {
    console.error("Erreur lors de la déconnexion:", error);
    throw error;
  }
}

const deleteUser = async (userId: number): Promise<void> => {
  try {
    const response = await DELETE(`/user/${userId}`);
    return response;
  } catch (error) {
    throw new Error(error.message);
  }
};

export const UserApiService = {
  getAllUsers,
  getAllManagersAdmin,
  getUserInfo,
  getUsers,
  updateUserRole,
  logOut,
  deleteUser,
};
