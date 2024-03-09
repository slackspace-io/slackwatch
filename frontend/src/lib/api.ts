export async function fetchPods(): Promise<Array<{name: string, timeScanned: string}>> {
  const response = await fetch('/api/pods');
  if (!response.ok) {
      throw new Error('Failed to fetch pods');
  }
  return response.json();
}
