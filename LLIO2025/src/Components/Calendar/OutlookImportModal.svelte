<script lang="ts">
    import { X, ChevronLeft, ChevronRight } from 'lucide-svelte';
    import {type Activity, type OutlookEvent, type Project} from "../../Models";
    import ActivityEntryForm from "./ActivityEntryForm.svelte";
    import {activityTemplate} from "../../forms/activity";

    type Props = {
        date: Date;
        events : OutlookEvent[];
        projects: Project[];
        onClose: () => void;
        onSuccess: () => void;
    };

    let {
        date,
        events,
        projects,
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
            alert(err.message);
        } finally {
        }
    };

    const handleNext = () => {
        if (selectedEventInt >= events.length - 1) {
            selectedEventInt = 0;
        } else selectedEventInt++;
    }
    const handlePrevious = () => {
        if (selectedEventInt <= 0) {
            selectedEventInt = events.length - 1;
        } else selectedEventInt--;
    }

    const roundToNearest15 = (dateStr) => {
        const date = new Date(dateStr);
        const ms = 15 * 60 * 1000;
        return new Date(Math.round(date.getTime() / ms) * ms);
    }

    let selectedEventInt = $state<number>(0);
    let selectedEvent = $derived(events[selectedEventInt]);
    let selectedEventProject = $derived(projects.find((p) => selectedEvent.subject.includes(p.uniqueId)))
    let activity = $state<Activity>({} as Activity);

    $effect(() => {
        const baseActivity = activityTemplate.generate();
        activity = {
            ...baseActivity,
            name: selectedEvent.subject,
            description: selectedEvent.body.content,
            startDate: roundToNearest15(selectedEvent.start),
            endDate: roundToNearest15(selectedEvent.end),
            projectId: selectedEventProject ? selectedEventProject.id : baseActivity.projectId,
            projectName: selectedEventProject ? selectedEventProject.name : baseActivity.projectName,
        }
    })
</script>

<div class="modal-overlay">
    <div class="modal">
        <div class="modal-header">
            <h2 class="modal-title">Importation des évènements Outlook du {date.toLocaleDateString("fr-CA")}<br/>Évènement {selectedEventInt + 1} de {events.length}</h2>
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
                <div class="nav-row">
                    <button
                            type="button"
                            class="nav-btn"
                            onclick={handlePrevious}
                            disabled={events.length <= 1}
                    >
                        <ChevronLeft size={16} />
                        Précédent
                    </button>
                    <span class="nav-label">{selectedEventInt + 1} / {events.length}</span>
                    <button
                            type="button"
                            class="nav-btn"
                            onclick={handleNext}
                            disabled={events.length <= 1}
                    >
                        Suivant
                        <ChevronRight size={16} />
                    </button>
                </div>

                <div class="form-group">
                    <ActivityEntryForm importIndex={selectedEventInt} activityToImport={activity} projects={projects} onClose={()=>{}} onDelete={()=>{}} onSubmit={()=>{}} onUpdate={()=>{}} activityToEdit={null} />

                    <div class="field-row">
                        <div class="field">
                            <label class="field-label">Nom</label>
                            <p class="field-value">{activity.name ?? '—'}</p>
                        </div>
                        <div class="field">
                            <label class="field-label">Projet</label>
                            <p class="field-value">{selectedEventProject ? `${selectedEventProject.uniqueId} — ${selectedEventProject.name}` : '—'}</p>
                        </div>
                    </div>

                    <div class="field-row">
                        <div class="field">
                            <label class="field-label">Début</label>
                            <p class="field-value">
                                {activity.startDate
                                    ? activity.startDate.toLocaleString("fr-CA", { dateStyle: "short", timeStyle: "short" })
                                    : '—'}
                            </p>
                        </div>
                        <div class="field">
                            <label class="field-label">Fin</label>
                            <p class="field-value">
                                {activity.endDate
                                    ? activity.endDate.toLocaleString("fr-CA", { dateStyle: "short", timeStyle: "short" })
                                    : '—'}
                            </p>
                        </div>
                        <div class="field">
                            <label class="field-label">Catégorie</label>
                            <p class="field-value">{activity.categoryId ?? '—'}</p>
                        </div>
                    </div>

                    <div class="field">
                        <label class="field-label">Description</label>
                        <p class="field-value description">{activity.description ?? '—'}</p>
                    </div>
                </div>

                <div class="modal-footer">
                    <button
                            type="button"
                            class="py-3 px-6 bg-gray-100 text-gray-700 rounded-lg font-medium hover:bg-gray-200 hover:-translate-y-0.5 hover:shadow-sm active:translate-y-0 transition border border-gray-200"
                            onclick={handleClose}
                    >
                        Annuler l'import
                    </button>
                    <button
                            type="submit"
                            class="py-3 px-6 bg-[#015e61] text-white rounded-lg font-medium hover:bg-[#014446] hover:-translate-y-0.5 hover:shadow-md active:translate-y-0 transition disabled:opacity-50"
                            onclick={handleSubmit}
                    >
                        Importer les évènements sélectionnés
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
        width: 900px;
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

    .nav-row {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 20px;
    }

    .nav-btn {
        display: flex;
        align-items: center;
        gap: 4px;
        padding: 6px 14px;
        font-size: 14px;
        font-weight: 500;
        color: #374151;
        background: #f3f4f6;
        border: 1px solid #d1d5db;
        border-radius: 6px;
        cursor: pointer;
        transition: background 0.15s;
    }

    .nav-btn:hover:not(:disabled) {
        background: #e5e7eb;
    }

    .nav-btn:disabled {
        opacity: 0.4;
        cursor: default;
    }

    .nav-label {
        font-size: 14px;
        color: #6b7280;
    }

    .form-group {
        display: flex;
        flex-direction: column;
        gap: 16px;
        margin-bottom: 16px;
        padding: 16px;
        background: #f9fafb;
        border: 1px solid #e5e7eb;
        border-radius: 8px;
    }

    .field-row {
        display: flex;
        gap: 24px;
    }

    .field {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 4px;
    }

    .field-label {
        font-size: 12px;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        color: #6b7280;
    }

    .field-value {
        font-size: 14px;
        color: #111827;
        margin: 0;
    }

    .field-value.description {
        white-space: pre-wrap;
        max-height: 80px;
        overflow-y: auto;
        color: #374151;
    }

    .modal-footer {
        display: flex;
        justify-content: flex-end;
        gap: 12px;
        margin-top: 24px;
    }
</style>