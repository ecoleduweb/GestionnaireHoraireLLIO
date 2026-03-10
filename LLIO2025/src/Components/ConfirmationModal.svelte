<script lang="ts">
  import { X } from 'lucide-svelte';
  import { string } from 'yup';

  type Props = {
    modalTitle: string;
    modalText: string;
    confirmText?: string;
    cancelText?: string;
    errorText: string;
    onClose: () => void;
    onSuccess: () => void;
  };

  let { 
    modalTitle, 
    modalText,
    confirmText = 'Confirmer',
    cancelText = 'Annuler',
    errorText, 
    onClose, 
    onSuccess 
  }: Props = $props();

  const handleClose = () => {
    onClose();
  };

  const handleSubmit = async () => {
    try {
      await onSuccess();
      onClose();
    } catch (err) {
        if (errorText != "") alert(errorText);
    } finally {
    }
  };
</script>

  <div class="modal-overlay">
    <div class="modal">
      <div class="modal-header">
        <h2 class="modal-title">{modalTitle}</h2>
        <button type="button" class="text-black hover:text-gray-600" onclick={handleClose}>
          <X />
        </button>
      </div>

      <div class="modal-content">
          <form
            class="flex flex-col h-full"
            onsubmit={(e) => {
              e.preventDefault();
            }}
          >
            <div class="form-group">
                <h2 class="modal-text">{modalText}</h2>
            </div>

            <div class="modal-footer">
                <button
                  type="button"
                  class="py-3 px-6 bg-gray-100 text-gray-700 rounded-lg font-medium hover:bg-gray-200 hover:-translate-y-0.5 hover:shadow-sm active:translate-y-0 transition border border-gray-200"
                  onclick={handleClose}
                >
                  {cancelText}
                </button>
                <button
                  type="submit"
                  class="py-3 px-6 bg-[#015e61] text-white rounded-lg font-medium hover:bg-[#014446] hover:-translate-y-0.5 hover:shadow-md active:translate-y-0 transition disabled:opacity-50"
                  onclick={handleSubmit}
                >
                  {confirmText}
                </button>
            </div>
          </form>
      </div>
    </div>
  </div>


<style>
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .modal {
    background-color: white;
    border-radius: 4px;
    width: 400px;
    max-width: 90%;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  }

  .modal-header {
    padding: 12px 24px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid #eee;
  }

  .modal-title {
    font-size: 18px;
    margin: 0;
    color: #333;
  }

  .modal-content {
    padding: 24px;
  }

  .form-group {
    margin-bottom: 16px;
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 24px;
  }
</style>
