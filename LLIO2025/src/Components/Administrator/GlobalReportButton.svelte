<script>
 import { ReportApiService } from '../../services/ReportApiService';
 import ExcelSvg from '../../assets/svg/excel.svg.svelte'

 let isLoading = false;

    const handleButtonClick = async () => {
        isLoading = true
        try{
            await ReportApiService.getReportExcel()
        } catch (error) {  
            alert('Erreur lors du téléchargement du fichier Excel');
        } finally{
            isLoading = false
        }
    };
</script>
<button 
    onclick={() => handleButtonClick()}
    id="excel-button"
    type="button" 
    class="w-full mt-4 py-2 px-4 text-sm font-medium transition-colors bg-[#e6f0f0] text-[#005e61] rounded-md hover:bg-[#d0e6e6] flex items-center justify-center cursor-pointer"
    disabled={isLoading}
    >
    <ExcelSvg />
    {#if !isLoading}
        Générer un export Excel
    {:else}
        Génération du fichier en cours...
    {/if} 
</button>