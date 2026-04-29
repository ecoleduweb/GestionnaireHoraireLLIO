<script lang="ts">
  import { formatHours } from '../../utils/date';
  import { Pencil, Trash2, Archive, ArchiveX } from 'lucide-svelte';
  import type { DetailedProject, Project, UserInfo } from '../../Models';
  import { UserRole } from '$lib/types/enums';
  import { getHoursColor } from '../../utils/displayUtils';

  type Props = {
    project: DetailedProject;
    currentUser: UserInfo | null;
    onEdit: (project: Project) => void;
    onArchive: (project: Project) => void;
    onDelete: (project: Project) => void;
  };

  let { project, currentUser, onArchive, onEdit, onDelete }: Props = $props();

  const handleEdit = (event: MouseEvent) =>{
    event.stopPropagation(); // Empêche la propagation du clic aux éléments parents
    onEdit(project);
  }

  const handleArchive = (event: MouseEvent) => {
    event.stopPropagation();
    onArchive(project);
  }

  const handleDelete = (event: MouseEvent) => {
    event.stopPropagation();
    onDelete(project);
  }
</script>

<div class="border-l-10 border-b" style="border-left-color: {project.color}">
  <div class="p-4">
    <div class="flex justify-between items-center">
      <div>
        <span class="text-black">{project.uniqueId}</span>
        <span class="text-gray-500 ml-2">|</span>
        <span class="text-black">{project.name}</span>
        <span class="text-gray-500 ml-2">|</span>
        <span class="text-black ml-2">{project.lead}</span>
      </div>
      <!-- svelte-ignore event_directive_deprecated -->
      {#if currentUser.role == UserRole.Admin || currentUser.role == UserRole.ProjectManager}
        <button
          class="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-full transition-colors"
          onclick={handleEdit}
          aria-label="Modifier le projet"
        >
          <Pencil size={16} />
        </button>
        <button
          class="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-full transition-colors"
          onclick={handleArchive}
          aria-label={!project.isArchived ? "Archiver le projet" : "Désarchiver le projet"}
        >
          {#if !project.isArchived}
            <Archive size={16} />
          {:else}
            <ArchiveX size={16} />
          {/if}
      
        </button>
      {/if}
      {#if currentUser.role == UserRole.Admin}
        <button
          class="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-full transition-colors"
          onclick={handleDelete}
          aria-label="Supprimer le projet"
        >
          <Trash2 size={16} />
        </button>
      {/if}
    </div>

    <div class="mt-2 flex items-center text-xs text-gray-500">
      <div class="mr-3">
        <span>Temps passé</span>
        <hr class="my-1" />
        <div class="font-bold text-black mr-4">{formatHours(project.totalTimeSpent)}</div>
      </div>
      <div class="mr-3">
        <span>Temps estimé</span>
        <hr class="my-1" />
        <div class="text-gray-400">
          {formatHours(project.totalTimeEstimated)}
        </div>
      </div>
      <div>
        <span>Temps restant</span>
        <hr class="my-1" />
        {#if project.totalTimeRemaining < 0}
          <div class="font-medium {getHoursColor(project.totalTimeSpent, project.totalTimeEstimated)}">
          {formatHours(project.totalTimeRemaining)}
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
