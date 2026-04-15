<script lang="ts">
  import { Pencil } from "lucide-svelte";
  import type { Category } from "../Models";
  import TextInputModal from "./TextInputModal.svelte";

    let { category, sendRenameCategory, setNewCategories, getCategories } = $props();

    let selectedCategory: Category = $state();
    let enableRenameCategoryPrompt = $state(false);

    function handleRenameCategory(category: Category) 
    {
        selectedCategory = category;
        enableRenameCategoryPrompt = true;
    }

</script>

{#if category.name !== "Par défaut"}
        <button
            class="renameButton justify-end ml-1 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-full transition-colors"
            onclick={() => handleRenameCategory(category)}
            aria-label="Renommer la catégorie"
        >
        <Pencil size={16} />
    </button>
{/if}

{#if enableRenameCategoryPrompt}
    <TextInputModal
      modalTitle="Modification du nom de la catégorie"
      modalText={`Renommez la catégorie "` + selectedCategory.name + `"`}
      defaultTextInValue={selectedCategory.name}
      onSuccess={async (val: string) => {
        if (await sendRenameCategory(selectedCategory, val)) {
            setNewCategories(getCategories().map((cat) => cat.id === selectedCategory.id ? {...cat, name: val} : cat));
        }
        enableRenameCategoryPrompt = false;
      }}
      onClose={() => {
        enableRenameCategoryPrompt = false;
      }}
    />
{/if}

<style>
    .renameButton {
        display: none;
    }
</style>