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
        onActivityImported: (activity : Activity) => void;
    };

    let {
        date,
        events,
        projects,
        onClose,
        onActivityImported
    }: Props = $props();

    const handleClose = () => {
        onClose();
    };

    const roundToNearest15 = (dateStr) => {
        const date = new Date(dateStr);
        const ms = 15 * 60 * 1000;
        return new Date(Math.round(date.getTime() / ms) * ms);
    }

    let selectedEventInt = $state<number>(0);
    let selectedEvent = $derived(events[selectedEventInt]);
    let selectedEventProject = $derived(projects.find((p) => selectedEvent.subject.includes(p.uniqueId)))
    let activity = $state<Activity>({} as Activity);
    let eventsDoneInt = $state<number[]>([]);
    let eventsLeftInt = $derived<number[]>((() : number[] => {
        let e = [];
        for (let i = 0; i < events.length; i++) {
            if(!eventsDoneInt.includes(i)) e.push(i);
        }
        return e;
    })())

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

    const handleNext = () => {
        let nextEventToDisplay = selectedEventInt;
        do {
            if (nextEventToDisplay >= events.length - 1) {
                nextEventToDisplay = 0;
            } else nextEventToDisplay++;
        } while (eventsLeftInt.length > 0 && !eventsLeftInt.includes(nextEventToDisplay))
        selectedEventInt = nextEventToDisplay;
    }

    const handlePrevious = () => {
        let nextEventToDisplay = selectedEventInt;
        do {
            if (nextEventToDisplay <= 0) {
                nextEventToDisplay = events.length - 1;
            } else nextEventToDisplay--;
        } while (eventsLeftInt.length > 0 && !eventsLeftInt.includes(nextEventToDisplay))
        selectedEventInt = nextEventToDisplay;
    }

    const handleDone = () => {
        eventsDoneInt.push(selectedEventInt);
        handleNext();
        if (eventsLeftInt.length === 0) onClose();
    }
</script>

<div class="modal-overlay">
    <div class="modal">
        <div class="modal-header">
            <h2 class="modal-title">Importation des évènements Outlook du {date.toLocaleDateString("fr-CA")}
                {#if eventsLeftInt.length > 0}
                    <br/>Évènement {selectedEventInt + 1} de {events.length}
                {/if}
            </h2>
            <button type="button" class="text-black hover:text-gray-600" onclick={handleClose}>
                <X />
            </button>
        </div>

        <div class="modal-content">
            {#if eventsLeftInt.length > 0}
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
                            <ActivityEntryForm activityToImport={activity} projects={projects} onClose={handleDone} onDelete={()=>{}} onSubmit={onActivityImported} onUpdate={()=>{}} activityToEdit={null} />
                        </div>
                </form>
            {:else}
                <p class="text-center">Tous les évènements ont été importés !<br/>Vous pouvez maintenant quitter cette fenêtre.</p>
            {/if}
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
        max-height: 90%;
        display: flex;
        flex-direction: column;
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
        overflow-y: auto;
        flex: 1;
        min-height: 0;
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
</style>