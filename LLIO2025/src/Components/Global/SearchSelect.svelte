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
    }

    let {
        id,
        name,
        items = [],
        selectedValue = $bindable(null),
        placeholder,
        required,
        setFields
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
/>
