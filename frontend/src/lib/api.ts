import { BE_BASE_URL } from "../config";

export async function fetchPods(): Promise<Array<{name: string, timeScanned: string}>> {
  // @ts-ignore
  const response = await fetch(`${BE_BASE_URL}/pods`);
  if (!response.ok) {
      throw new Error('Failed to fetch pods');
  }
  return response.json();
}
