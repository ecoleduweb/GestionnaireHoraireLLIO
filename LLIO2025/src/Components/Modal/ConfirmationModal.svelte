<!-- ConfirmModal.svelte -->
<script lang="ts">
  import BaseModal from "./BaseModal.svelte";

  type Props = {
    modalTitle: string;
    modalText: string;
    confirmText?: string;
    cancelText?: string;
    errorText?: string;
    onClose: () => void;
    onSuccess: () => void | Promise<void>;
  };

  let {
    modalTitle,
    modalText,
    confirmText = 'Confirmer',
    cancelText = 'Annuler',
    errorText = '',
    onClose,
    onSuccess
  }: Props = $props();
</script>

<BaseModal
  {modalTitle}
  {confirmText}
  {cancelText}
  {errorText}
  {onClose}
  {onSuccess}
>
  {#snippet children()}
    <div class="form-group">
      <h2 class="modal-text">
        <b>{modalText}</b>
      </h2>
    </div>
  {/snippet}
</BaseModal>
