export async function fetchPods(): Promise<Array<{name: string, timeScanned: string}>> {
  // @ts-ignore
  const baseUrl = import.meta.env.VITE_API_BASE_URL; // Correctly access the environment variable

  const response = await fetch(`${baseUrl}/pods`);
  if (!response.ok) {
      throw new Error('Failed to fetch pods');
  }
  return response.json();
}
