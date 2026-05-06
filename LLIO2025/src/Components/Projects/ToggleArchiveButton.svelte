<script lang="ts">
    import { Archive, ArchiveX } from 'lucide-svelte';
    import type { DetailedProject, Project } from '../../Models';
    import { ProjectApiService } from '../../services/ProjectApiService';
    import ConfirmationModal from '../Modal/ConfirmationModal.svelte';

    type Props = {
        project: DetailedProject;
        onArchive: (project: Project) => void;
    };

    let showArchiveModal = $state(false);

    let { project, onArchive }: Props = $props();

    const handleCloseModal = () =>{
        showArchiveModal = false;
    }

    const handleSuccessArchive = async () => {
        const data = await ProjectApiService.toggleArchiveProject(project.id, !project.isArchived)
        console.log(data);
        if (data) onArchive(data);
    }

    const handleOnArchiveClick = (event: MouseEvent) => {
        event.stopPropagation();
        showArchiveModal = true;
    }

</script>

<button
    class="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-full transition-colors"
    onclick={handleOnArchiveClick}
    aria-label={!project.isArchived ? "Archiver le projet" : "Désarchiver le projet"}
>
    {#if !project.isArchived}
        <Archive size={16} />
    {:else}
        <ArchiveX size={16} />
    {/if}
</button>

{#if showArchiveModal}
<ConfirmationModal
  modalTitle="Archiver un projet"
  modalText="Voulez-vous vraiment {!project.isArchived ? "archiver" : "désarchiver"} le projet {project.name} ?"
  errorText="Erreur lors de la archivation du projet."
  onSuccess={handleSuccessArchive}
  onClose={handleCloseModal}
/>
{/if}