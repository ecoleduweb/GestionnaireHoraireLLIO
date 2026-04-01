export function getHoursColor(spent: number, estimated: number): string {
  if (estimated && spent > estimated) return 'text-red-700';
  return 'text-gray-700';
}
