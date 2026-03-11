import type { TimeBankConfig } from '../Models/index';
import type { TimeBankResponse } from '../Models/index';
import type { User, UserInfo } from '../Models/index';
import { DELETE, GET, PATCH, POST } from '../ts/server';

const getUserInfo = async (): Promise<UserInfo> => {
  const response = await GET<UserInfo>('/user/me');
  return response;
};

const getAllUsers = async (): Promise<User[]> => {
  const response = await GET<User[]>("/users");
  return response;
};

const getAllManagersAdmin = async (): Promise<User[]> => {
  const response = await GET<User[]>("/users?role=1&role=2");
  return response;
};

const getUsers = async (): Promise<UserInfo[]> => {
  const response = await GET<UserInfo[]>('/users');
  return response;
};

const updateUserRole = async (userId: number, role: number): Promise<void> => {
  await PATCH(`/user/${userId}/role`, { role });
};

const logOut = async (): Promise<void> => {
  await POST('/logout', {});
};

const deleteUser = async (userId: number): Promise<void> => {
  await DELETE(`/user/${userId}`);
};

const saveTimeBankConfig = async (config: TimeBankConfig): Promise<TimeBankConfig> => {
  const response = await POST<TimeBankConfig, TimeBankConfig>(
    '/user/time-bank/config',
    config
  );
  return response;
};

const getTimeBankConfig = async (): Promise<TimeBankConfig | null> => {
  try {
    const response = await GET<TimeBankConfig>('/user/time-bank/config');
    return response;
  } catch (error: any) {
    if (error?.name === 'NotFound') {
      return null;
    }
    throw error;
  }
};

const getTimeInBank = async (): Promise<TimeBankResponse> => {
  const response = await GET<TimeBankResponse>('/user/time-bank');
  return response;
};

export const UserApiService = {
  getAllUsers,
  getAllManagersAdmin,
  getUserInfo,
  getUsers,
  updateUserRole,
  logOut,
  deleteUser,
  saveTimeBankConfig,
  getTimeBankConfig,
  getTimeInBank,
};