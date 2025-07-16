<script lang="ts">
  // Propriétés du composant
  type Props = {
    show: boolean;
    title: string;
    message: string;
    confirmText?: string;
    cancelText?: string;
    onConfirm: () => void;
    onCancel: () => void;
  };

  let {
    show,
    title,
    message,
    confirmText = 'Confirmer',
    cancelText = 'Annuler',
    onConfirm,
    onCancel,
  }: Props = $props();
</script>

{#if show}
  <div class="fixed inset-0 z-50 flex items-center justify-center">
    <!-- Overlay semi-transparent avec la nouvelle syntaxe de Tailwind v4 -->
    <div class="absolute inset-0 bg-gray-800/50 transition-opacity" onclick={onCancel}></div>

    <!-- Modal -->
    <div class="bg-white rounded-lg shadow-xl p-6 w-80 max-w-md z-10">
      <h3 class="text-lg font-medium text-gray-900 mb-4">{title}</h3>

      <p class="text-gray-600 mb-6">
        {@html message}
      </p>

      <div class="flex justify-end gap-3">
        <button
          type="button"
          class="py-2 px-4 bg-gray-100 text-gray-700 rounded-lg font-medium hover:bg-gray-200 transition"
          onclick={onCancel}
        >
          {cancelText}
        </button>
        <button
          type="button"
          class="py-2 px-4 bg-[#015e61] text-white rounded-lg font-medium hover:bg-[#014446] transition"
          onclick={onConfirm}
        >
          {confirmText}
        </button>
      </div>
    </div>
  </div>
{/if}
