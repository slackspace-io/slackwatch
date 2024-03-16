import { unstable_noStore as noStore } from 'next/cache';
import UpdateCard from '@/components/UpdateCard';


interface CombinedData {
  containerName: string;
  currentTag?: string;
  newTag?: string;
  foundAt?: string;
  image: string;
  includePattern: string;
  excludePattern: string;
  podName: string;
  timeScanned: string;
  updateAvailable: boolean;
  namespace: string;
  sentTime?: string;
  gitopsRepo: string;
  directory?: string;
}
async function getData(): Promise<CombinedData[]> {
  noStore();
  const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;
  if (!baseUrl) {
    console.warn('NEXT_PUBLIC_API_BASE_URL is not defined. Skipping fetch.');
    return [];
  }
  const res = await fetch(`${baseUrl}/api/data/combined`);
  if (!res.ok) {
    throw new Error('Failed to fetch data');
  }
  return res.json();
}

export default async function Page() {
  let data = await getData();

  data = data.sort((a, b) => Number(b.updateAvailable) - Number(a.updateAvailable));

  return (
    <main className="p-4">
      {data.map((update, index) => (
        <UpdateCard key={index} update={update} />
      ))}
    </main>
  );
}
