<script lang="ts">
  import { onMount } from "svelte";
  import "../../style/app.css"
  import type { UserInfo } from '../../Models/index.ts';
  import { UserRole } from '../../lib/types/enums';
  import { UserApiService } from '../../services/UserApiService';
  import { goto } from "$app/navigation";

  let users = $state<UserInfo[]>([]);
  let isLoading = $state(false);
  let selectedUser = $state(null);
  let selectedRole = $state(undefined);

  const handleConfirmClick = async () =>{
    if (!selectedUser) return;
    if (selectedRole === undefined) return;
    try {
      await UserApiService.updateUserRole(selectedUser.id, selectedRole);
      alert('Rôle mis à jour avec succès');
      users = await UserApiService.getUsers();
    } catch (error) {
      console.error('Error updating user role:', error);
      alert('Erreur lors de la mise à jour du rôle');
    }
  }

  const handleDeleteUser = async () => {
    if (!selectedUser) return;
    try {
      await UserApiService.deleteUser(selectedUser.id);
      users = users.filter(user => user.email !== selectedUser.email);
      alert('Utilisateur supprimé avec succès');
    } catch (error) {
      alert('Erreur lors de la suppression de l\'utilisateur');
    }

  }

  onMount(async () => {
    isLoading = true;
    try {
      users = await UserApiService.getUsers()
    } catch (error) {
      console.error('Error fetching users:', error);
    } finally {
      isLoading = false;
    }
  });
</script>

<div class="bg-gray-100 min-h-screen">
  <div class="p-4 flex flex-col items-center">
    <h1 class="text-2xl font-medium text-gray-800">Modifier le rôle d'un utilisateur</h1>
    
    <div class="mt-4 p-4 bg-white rounded shadow w-1/2">
      {#if isLoading}
        <p>Chargement...</p>
      {:else}
        <div class="flex flex-col gap-4">
          <div>
            <label for="userSelect" class="block text-sm font-medium text-gray-700">Sélectionner un utilisateur</label>
            <select id="userSelect" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2" bind:value={selectedUser}>
              <option value={null}>Choisir un utilisateur</option>
              {#each users as user}
                <option value={user}>{user.firstName} {user.lastName} | {user.email}</option>
              {/each}
            </select>
          </div>
            <div>
            <label for="roleSelect" class="block text-sm font-medium text-gray-700">Sélectionner un rôle</label>
            <select id="roleSelect" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2" bind:value={selectedRole}>
              <option value={UserRole.User}>Utilisateur</option>
              <option value={UserRole.ProjectManager}>Chargé de projet</option>
              <option value={UserRole.Admin}>Administrateur</option>
            </select>
            </div>
          
            <div class="flex gap-4">
            <button 
              onclick={handleConfirmClick} 
              class="bg-[#015e61] hover:bg-[#014446] text-white font-bold py-2 px-4 rounded transition-colors cursor-pointer"
              disabled={!selectedUser || !selectedRole}>
              Confirmer
            </button>
            <button
              onclick={() => {
                handleDeleteUser();
              }} 
              class="bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded transition-colors cursor-pointer"
              disabled={!selectedUser}>
              Supprimer
            </button>
            <button 
              onclick={() => goto('./calendar')} 
              class="bg-gray-500 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded transition-colors cursor-pointer">
              Retour
            </button>
            </div>
        </div>
      {/if}
    </div>
  </div>
</div>