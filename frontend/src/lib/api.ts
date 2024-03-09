export async function fetchPods(): Promise<Array<{name: string, timeScanned: string}>> {
  const baseUrl = process.env.VITE_API_BASE_URL; // Access the environment variable correctly
  const response = await fetch(`${baseUrl}/pods`);
  if (!response.ok) {
      throw new Error('Failed to fetch pods');
  }
  return response.json();
}
