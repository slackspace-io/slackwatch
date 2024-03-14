import { Card } from '@/components/ui/card';
import { unstable_noStore as noStore } from 'next/cache'


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
  sentTime?: string;
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

  // Sort data so that entries with updates are at the top
  data = data.sort((a, b) => Number(b.updateAvailable) - Number(a.updateAvailable));

  return (
    <main className="p-4">
      <div className="mb-4">
      </div>
      {data.map((update: CombinedData, index: number) => (
        <Card key={index} className={`mb-4 p-4 shadow-lg ${update.updateAvailable ? 'border-l-4 border-green-500' : ''}`}>
          <p className="text-lg font-bold">
            Container: {update.containerName}
          </p>
          {update.currentTag && <p>Current Tag: {update.currentTag}</p>}
          {update.updateAvailable && update.newTag && (
            <p className="text-green-500">New Tag Available: {update.newTag}</p>
          )}
          <p>Image: {update.image}</p>
          <p>Pod Name: {update.podName}</p>
          <p>Time Scanned: {update.timeScanned}</p>
          {update.sentTime && <p>Notification Sent At: {update.sentTime}</p>}
        </Card>
      ))}
    </main>
  );
}