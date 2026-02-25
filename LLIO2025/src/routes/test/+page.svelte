<script lang="ts">
  import { ReportApiService } from '../../services/ReportApiService';

  const downloadReport = async () => {
    try {
      const blob = await ReportApiService.getReportCsv();
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'export.csv';
      document.body.appendChild(a);
      a.click();
      a.remove();
      URL.revokeObjectURL(url);
    } catch (err) {
      console.error('Erreur lors du téléchargement du rapport :', err);
      alert('Impossible de récupérer le rapport. Vérifiez la console.');
    }
  }
</script>

<style>
  .btn-download {
    padding: 0.6rem 1rem;
    background-color: #2b6cb0;
    color: white;
    border: none;
    border-radius: 6px;
    cursor: pointer;
  }
</style>

<button class="btn-download" on:click={downloadReport}>Télécharger le rapport CSV</button>