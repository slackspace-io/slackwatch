export async function fetchPods(): Promise<string> {
  const response = await fetch('/api/pods');
  if (!response.ok) {
      throw new Error('Failed to fetch pods');
  }
  return response.text();
}
