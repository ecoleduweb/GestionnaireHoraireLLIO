<script lang="ts">
    import Select from 'svelte-select';
    import type { SelectItem } from "$lib/types/SelectItem";
    import type { FieldsSetter } from '@felte/core';

    type Props = {
        id?: string;
        name?: string;
        items?: SelectItem[];
        selectedValue?: any;
        placeholder?: string;
        required?: boolean;
        setFields?: FieldsSetter<Record<string, unknown>>;
        focused?: boolean;
        onSubmit?: () => void;
    }

    let {
        id,
        name,
        items = [],
        selectedValue = $bindable(null),
        placeholder,
        required,
        setFields,
        focused = $bindable(false),
        onSubmit,
    } : Props = $props();

    // L'item sélectionné, pour affichage de svelte-select
    // Le derived est comme un effect appliqué directement à la variable. Elle ne peut donc pas être modifiée directement.
    let selectedItem = $derived(items.find(i => i.value === selectedValue) || null);

    /* ** Comment fonctionne la sélection ? **
    *
    * `selectedItem` contrôle seulement l'affichage de svelte-select.
    *
    * selectedValue (des Props) est la source de vérité.
    *   - Quand on choisit ou supprime dans le formulaire, ceci déclenche la méthode `handleSelect` qui
    *       modifie selectedValue, modifiant à son tour selectedItem change en raison du $derived.
    *   - Quand la valeur est changée par la page/composante parente, selectItem change aussi en raison du $derived.
    *
    * setFields() sert pour la validation des formulaires par Felte, qui ne fonctionne pas en utilisant un champ
    * hidden ou le nom fourni au svelte-select
    * */
    function handleSelect(e: CustomEvent<SelectItem | null>) {
        const nextValue = e.detail?.value ?? null;
        selectedValue = nextValue;
        if (setFields && name) {
            setFields(name, nextValue);
        }
    }
    // La partie suivante gère l'envoi du formulaire par la touche entrée du clavier.
    //
    // listOpen est utilisé pour détecter si la liste des options est ouverte. Ceci est utile pour détecter manuellement
    // si un appui sur la touche Entrée sert à choisir une option ou pour envoyer le formulaire.
    //
    // inputElement est utilisé pour y insérer l'évènement de keydown pour détecter l'appui d'une touche, puisque le
    // vrai on:keydown du svelte-select ne fonctionne pas.
    //
    // onKeydownHandle gère la logique de quelle touche a été appuyée et pourquoi, pour ensuite déclencher l'évènement
    // onSubmit.
    //
    // L'effect attache le *event handler* au select.
    let listOpen = $state(false);
    let inputElement = $state<HTMLInputElement | null>(null);

    function onKeydownHandler(e: KeyboardEvent) {
        if (selectedItem != null && !listOpen && (e.key === 'Enter' || e.key === 'Return')) {
            onSubmit?.();
        }
    }

    $effect(() => {
        if (!inputElement) return;
        inputElement.addEventListener('keydown', onKeydownHandler);
        return () => inputElement?.removeEventListener('keydown', onKeydownHandler);
    });
</script>

<Select
        {id}
        {name}
        {items}
        value={selectedItem}
        on:change={handleSelect}
        on:clear={() => handleSelect(new CustomEvent('clear', { detail: null }))}
        {placeholder}
        {required}
        bind:focused={focused}
        bind:listOpen
        bind:input={inputElement}
        --item-height="auto">
    <div slot="item" class="item" let:item>
        {item.label}
    </div>
</Select>

<style>
    .item {
        min-height: 16px;
        padding: 10px 0;
        line-height: 16px;
        display: flex;
        line-break: auto;
        white-space: pre-wrap;
        align-items: center;
    }
</style>