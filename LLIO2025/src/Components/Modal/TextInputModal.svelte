<script lang="ts">
  import BaseModal from "./BaseModal.svelte";


  type Props = {
    modalTitle: string;
    modalText: string;
    confirmText?: string;
    cancelText?: string;
    defaultTextInValue: string;
    errorText?: string;
    onClose: () => void;
    onSuccess: (value: string) => void | Promise<void>;
  };

let {
    modalTitle,
    modalText,
    confirmText = 'Confirmer',
    cancelText = 'Annuler',
    defaultTextInValue = '',
    errorText = '',
    onClose,
    onSuccess
  }: Props = $props();

  let textValue = $state("");
  
  $effect(() => {
    textValue = defaultTextInValue;
  });

  const submit = async () => {
    await onSuccess(textValue);
    onClose();
  };
</script>

<BaseModal
  {modalTitle}
  {confirmText}
  {cancelText}
  {errorText}
  {onClose}
  onSuccess={submit}
>
  {#snippet children()}
    <div class="form-group">
      <h2 class="modal-text">
        <b>{modalText}</b>
      </h2>
    </div>

    <input bind:value={textValue} class="form-input" type="text" />
  {/snippet}
</BaseModal>

<style>
  .modal-text {
    padding-bottom: 15px
  }

  .form-input {
    width: 100%;
    padding: 0.75rem 1rem;
    border: 1px solid #d1d5db;
    border-radius: 0.5rem;
    transition: all 0.2s;
    background-color: white;
    color: #4b5563;
  }
</style>
