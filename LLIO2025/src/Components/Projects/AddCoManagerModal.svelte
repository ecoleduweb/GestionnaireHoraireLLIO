<script lang="ts">
  // Propriétés du composant
  import type { User } from '../../Models';

  type Props = {
    show: boolean;
    projectName: string;
    users: User[];
    valueSelected: string;
    onAdd: () => void;
    onCancel: () => void;
  };

  let {
    show,
    projectName,
    users,
    valueSelected = $bindable(""),
    onAdd,
    onCancel,
  }: Props = $props();
</script>

{#if show}
  <div class="fixed inset-0 z-50 flex items-center justify-center">
    <!-- Overlay semi-transparent avec la nouvelle syntaxe de Tailwind v4 -->
    <div class="absolute inset-0 bg-gray-800/50 transition-opacity" onclick={onCancel}></div>

    <!-- Modal -->
    <div class="bg-white rounded-lg shadow-xl p-8 w-full max-w-xl z-10">
      <h3 class="text-lg font-medium text-gray-900 mb-4">Ajouter un co-chargé</h3>

      <p class="text-gray-600 mb-6">
        Sélectionnez un utilisateur pour devenir co-chargé de projet de {projectName}
      </p>

      <div class="flex gap-3 mb-6">
        <label for="userId" class="block text-gray-700">Utilisateur</label>
        <select name="userId" id="userId" bind:value={valueSelected}>
          <option value="" disabled hidden>{ users.length === 0 ? "Aucun utilisateur disponible à l'ajout" : "Choisir un utilisateur" }</option>
          {#each users as user}
            <option value={user.id}>
              {user.firstName} {user.lastName}
            </option>
          {/each}
        </select>
      </div>

      <div class="flex justify-end gap-3">
        <button
          type="button"
          class="py-2 px-4 bg-gray-100 text-gray-700 rounded-lg font-medium hover:bg-gray-200 transition"
          onclick={onCancel}
        >
          Annuler
        </button>
        <button
          type="button"
          class="py-2 px-4 bg-[#015e61] text-white rounded-lg
          font-medium hover:bg-[#014446] transition
          disabled:bg-gray-300 disabled:text-gray-500
          disabled:cursor-not-allowed disabled:hover:bg-gray-300"
          onclick={onAdd}
          disabled={!valueSelected}
        >
          Ajouter
        </button>
      </div>
    </div>
  </div>
{/if}
