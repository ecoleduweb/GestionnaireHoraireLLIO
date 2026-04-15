const VITE_API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

const getReportExcel = async () => {
  let url
  try {
    const response = await fetch(`${VITE_API_BASE_URL}/report/excel`, {
      method: 'GET',
      credentials: 'include',
    });
    if (!response.ok) throw new Error('Erreur lors du téléchargement');
    const blob = await response.blob();
    url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'activities.xlsx';
    document.body.appendChild(a);
    a.click();
    a.remove();
  } catch (error) {
    console.error(error);
    throw error;
  } finally{
    if (url) window.URL.revokeObjectURL(url);
  }
};

export const ReportApiService = { getReportExcel };
