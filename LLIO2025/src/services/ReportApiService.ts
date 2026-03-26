import { GET } from '../ts/server';

const getReportCsv = async ():Promise<Blob> => {
  try {
    const response:Response = await GET('/report/csv');
    return await response.blob();
  } catch (error) {
    console.error('Erreur lors de la écuperation du fichier csv : ', error);
    throw error;
  }
};

export const ReportApiService = {
  getReportCsv,
};
