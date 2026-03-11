<script lang="ts">
import { X } from "lucide-svelte";
import { UserApiService } from "../../services/UserApiService";
import type { TimeBankConfig } from "../../Models/index";
import { validateTimeBankForm } from "../../Validation/TimeBank";

type Props = {
  onClose: () => void;
  onSave: () => void;
};

let { onClose, onSave }: Props = $props();

const config = $state<TimeBankConfig>({
  startDate: "",
  hoursPerWeek: 0,
  offset: 0
});

let isSubmitting = $state(false);

const handleSubmit = async () => {
  if (isSubmitting) return;

  try {
    isSubmitting = true;

    await UserApiService.saveTimeBankConfig(config);

    onSave();
    onClose();

  } catch (err) {
    console.error(err);
    alert("Erreur lors de la configuration de la banque d'heure");
  } finally {
    isSubmitting = false;
  }
};

const { form, errors } = validateTimeBankForm(handleSubmit, config);
</script>

<div class="modal-overlay">

<div class="modal">

<div class="modal-header">

<h2 class="modal-title">Configuration des heures en banque</h2>

<button
type="button"
class="text-black hover:text-gray-600"
onclick={onClose}
>
<X />
</button>

</div>

<div class="modal-content">

<form
class="flex flex-col h-full"
use:form
onsubmit={(e) => e.preventDefault()}
>

<div class="form-group">

<label for="startDate">Début de la période</label>

<input
id="startDate"
name="startDate"
type="date"
bind:value={config.startDate}
/>

{#if $errors.startDate}
<span class="error-text">
{$errors.startDate}
</span>
{/if}

</div>


<div class="form-group">

<label for="hoursWorked">Nombre d'heures par semaine</label>

<input
id="hoursWorked"
name="hoursPerWeek"
type="number"
min="0"
step="0.5"
bind:value={config.hoursPerWeek}
/>

{#if $errors.hoursPerWeek}
<span class="error-text">
{$errors.hoursPerWeek}
</span>
{/if}

</div>


<div class="form-group">

<label for="offset">Décalage (offset)</label>

<input
id="offset"
name="offset"
type="number"
min="0"
step="0.5"
bind:value={config.offset}
/>

{#if $errors.offset}
<span class="error-text">
{$errors.offset}
</span>
{/if}

</div>


<div class="modal-footer">

<button
type="button"
class="btn-cancel"
onclick={onClose}
>
Annuler
</button>

<button
type="submit"
class="btn-save"
disabled={isSubmitting}
>
{isSubmitting ? "En cours..." : "Enregistrer"}
</button>

</div>

</form>

</div>

</div>

</div>

<style>

.modal-overlay{
position:fixed;
top:0;
left:0;
right:0;
bottom:0;
background-color:rgba(0,0,0,0.5);
display:flex;
align-items:center;
justify-content:center;
z-index:1000;
}

.modal{
background-color:white;
border-radius:4px;
width:400px;
max-width:90%;
box-shadow:0 2px 10px rgba(0,0,0,0.1);
}

.modal-header{
padding:12px 24px;
display:flex;
justify-content:space-between;
align-items:center;
border-bottom:1px solid #eee;
}

.modal-title{
font-size:18px;
margin:0;
color:#333;
}

.modal-content{
padding:24px;
}

.form-group{
margin-bottom:16px;
}

label{
display:block;
margin-bottom:8px;
color:#666;
font-size:14px;
}

input{
width:100%;
padding:8px 12px;
border:1px solid #ccc;
border-radius:4px;
box-sizing:border-box;
}

.modal-footer{
display:flex;
justify-content:flex-end;
gap:12px;
margin-top:24px;
}

.btn-cancel,
.btn-save{
padding:10px 20px;
border:none;
border-radius:6px;
font-weight:500;
cursor:pointer;
font-size:14px;
transition:all 0.2s ease;
}

.btn-cancel{
background:#f3f3f3;
color:#333;
border:1px solid #ddd;
}

.btn-cancel:hover{
background:#e8e8e8;
}

.btn-save{
background:#015e61;
color:white;
}

.btn-save:hover{
background:#014446;
}

.error-text{
color:#e53e3e;
font-size:14px;
margin-top:4px;
}

</style>