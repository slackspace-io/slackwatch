import { unstable_noStore as noStore } from 'next/cache';
import UpdateCard from '@/components/UpdateCard';


interface CombinedData {
    name: string,
    exclude_pattern: string,
    git_ops_repo: string,
    include_pattern: string,
    update_available: string,
    image: string,
    last_scanned: string,
    namespace: string,
    current_version: string,
    latest_version: string,
}
async function getData(): Promise<CombinedData[]> {
  noStore();
  const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;
  if (!baseUrl) {
    console.warn('NEXT_PUBLIC_API_BASE_URL is not defined. Skipping fetch.');
    return [];
  }
  const res = await fetch(`${baseUrl}/api/workloads/all`);
  if (!res.ok) {
    throw new Error('Failed to fetch data');
  }
  return res.json();
}

export default async function Page() {
  let data = await getData();

  data = data.sort((a, b) => Number(b.update_available) - Number(a.update_available));

  return (
    <main className="p-4">
      {data.map((update, index) => (
        <UpdateCard key={index} update={update} />
      ))}
    </main>
  );
}
