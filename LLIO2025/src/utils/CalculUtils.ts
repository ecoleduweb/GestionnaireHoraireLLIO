export const calculateEmployeeTime = (employee: any, type: 'spent' | 'estimated'): number => {
  return employee.categories.reduce(
    (sum: number, cat: any) => sum + (type === 'spent' ? cat.timeSpent : cat.timeEstimated),
    0
  );
};
export const calculateRemainingTime = (timeSpent: number, timeEstimated: number): number => {
  return timeEstimated - timeSpent;
};
